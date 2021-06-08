// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firestore_service

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

const idTokenName = "id"

type FirestoreKVService struct {
	client  *firestore.Client
	context context.Context
	sdk.UnimplementedKeyValuePlugin
}

func addMapKeys(values map[string]interface{}, keys map[string]interface{}) map[string]interface{} {
	// Clone values
	newMap := make(map[string]interface{})
	for key, value := range values {
		newMap[key] = value
	}

	// Add keys with "_" prefix
	for key, value := range keys {
		newMap["_"+key] = value
	}

	return newMap
}

func stripMapKeys(source map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, value := range source {
		if !strings.HasPrefix(key, "_") {
			newMap[key] = value
		}
	}
	return newMap
}

func (s *FirestoreKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return nil, err
	}

	keyValue := kv.GetKeyValue(key)

	value, err := s.client.Collection(collection).Doc(keyValue).Get(s.context)

	if err != nil {
		return nil, fmt.Errorf("error retrieving value: %v", err)
	}

	results := stripMapKeys(value.Data())

	return results, nil
}

func (s *FirestoreKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	keyValue := kv.GetKeyValue(key)

	// Add keys to value item to support querying
	value = addMapKeys(value, key)

	_, err = s.client.Collection(collection).Doc(keyValue).Set(s.context, value)

	if err != nil {
		return fmt.Errorf("error updating value: %v", err)
	}

	return nil
}

func (s *FirestoreKVService) Delete(collection string, key map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}
	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}

	keyValue := kv.GetKeyValue(key)

	_, err = s.client.Collection(collection).Doc(keyValue).Delete(s.context)

	if err != nil {
		return fmt.Errorf("error deleting value: %v", err)
	}

	return nil
}

func (s *FirestoreKVService) Query(collection string, expressions []sdk.QueryExpression, limit int, pagingToken map[string]interface{}) (*sdk.QueryResult, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}
	err = kv.ValidateExpressions(collection, expressions)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	var resultPagingToken map[string]interface{}

	if len(expressions) > 0 {
		stack, err := kv.Stack()
		if err != nil {
			return nil, err
		}
		indexes, _ := stack.CollectionIndexes(collection)

		query := s.client.Collection(collection).Offset(0)

		appliedOrderBy := false

		for _, exp := range expressions {
			// If operand is key/index then prefix with "_"
			expOperand := exp.Operand
			if utils.IndexOf(indexes, expOperand) != -1 {
				expOperand = "_" + expOperand
			}
			if exp.Operator == "startsWith" {
				endRangeValue := kv.GetEndRangeValue(exp.Value)
				query = query.Where(expOperand, ">=", exp.Value).Where(expOperand, "<", endRangeValue)

			} else {
				query = query.Where(expOperand, exp.Operator, exp.Value)
			}

			if exp.Operator != "==" && limit > 0 && !appliedOrderBy {
				query = query.OrderBy(expOperand, firestore.Asc)
				appliedOrderBy = true
			}
		}

		if limit > 0 {
			query = query.Limit(limit)

			if len(pagingToken) > 0 {
				if id, found := pagingToken[idTokenName]; found {
					query = query.StartAfter(id)
				}
			}
		}

		itr := query.Documents(s.context)

		for {
			docSnp, err := itr.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error querying value: %v", err)
			}
			value := stripMapKeys(docSnp.Data())
			results = append(results, value)

			// If query limit configured determine continue token
			if limit > 0 && len(results) == limit {
				resultPagingToken = make(map[string]interface{})
				resultPagingToken[idTokenName] = docSnp.Ref.ID
			}
		}

	} else {
		query := s.client.Collection(collection).Offset(0)

		if limit > 0 {
			if len(pagingToken) > 0 {
				if id, found := pagingToken[idTokenName]; found {
					query = query.StartAfter(id)
				}
			}
		}

		query = query.OrderBy(firestore.DocumentID, firestore.Asc)

		iter := query.Documents(s.context)

		for {
			docSnp, err := iter.Next()
			// Break done or fetch limit reached
			if err == iterator.Done || (limit > 0 && len(results) == limit) {
				break
			}
			if err != nil {
				return nil, err
			}
			value := stripMapKeys(docSnp.Data())
			results = append(results, value)

			// If query limit configured determine continue token
			if limit > 0 && len(results) == limit {
				resultPagingToken = make(map[string]interface{})
				resultPagingToken[idTokenName] = docSnp.Ref.ID
			}
		}
	}

	return &sdk.QueryResult{
		Data:        results,
		PagingToken: resultPagingToken,
	}, nil
}

func New() (sdk.KeyValueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := firestore.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("firestore client error: %v", clientError)
	}

	return &FirestoreKVService{
		client:  client,
		context: ctx,
	}, nil
}

func NewWithClient(client *firestore.Client, ctx context.Context) (sdk.KeyValueService, error) {
	return &FirestoreKVService{
		client:  client,
		context: ctx,
	}, nil
}
