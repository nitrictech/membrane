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
package kv_test

import (
	"errors"

	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases

var _ = Describe("KeyValue Plugin", func() {

	When("ValidateCollection", func() {
		When("Blank collection", func() {
			It("should return error", func() {
				err := kv.ValidateCollection("")
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank collection")))
			})
		})
		When("Valid collection", func() {
			It("should return nil", func() {
				err := kv.ValidateCollection("collection")
				Expect(err).To(BeNil())
			})
		})
	})

	When("ValidateKeyMap", func() {
		When("Nil key", func() {
			It("should return error", func() {
				err := kv.ValidateKeyMap(nil)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-nil key")))
			})
		})
		When("Empty key", func() {
			It("should return error", func() {
				keyMap := map[string]interface{}{}
				err := kv.ValidateKeyMap(keyMap)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-empty key")))
			})
		})
		When("Blank key value", func() {
			It("should return error", func() {
				keyMap := map[string]interface{}{
					"key": "",
				}
				err := kv.ValidateKeyMap(keyMap)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank key value")))
			})
		})
		When("Too many key values", func() {
			It("should return error", func() {
				keyMap := map[string]interface{}{
					"key1": "1",
					"key2": "2",
					"key3": "3",
				}
				err := kv.ValidateKeyMap(keyMap)
				Expect(err).To(BeEquivalentTo(errors.New("provide key with 1 or 2 items")))
			})
		})
	})

	When("GetKeyValue", func() {
		When("Single key", func() {
			It("should return single key value", func() {
				keyMap := map[string]interface{}{
					"key": "user@server.com",
				}
				key := kv.GetKeyValue(keyMap)
				Expect(key).To(BeEquivalentTo("user@server.com"))
			})
		})
		When("Multi key", func() {
			It("should return appended key values", func() {
				keyMap := map[string]interface{}{
					"pk": "Customer#123",
					"sk": "Order#456",
				}
				key := kv.GetKeyValue(keyMap)
				Expect(key).To(BeEquivalentTo("Customer#123_Order#456"))
			})
		})
	})

	When("GetKeyValues", func() {
		When("Single key", func() {
			It("should return single key value", func() {
				keyMap := map[string]interface{}{
					"key": "user@server.com",
				}
				keys := kv.GetKeyValues(keyMap)
				Expect(len(keys)).To(BeEquivalentTo(1))
				Expect(keys[0]).To(BeEquivalentTo("user@server.com"))
			})
		})
		When("Multi key", func() {
			It("should return key values", func() {
				keyMap := map[string]interface{}{
					"pk": "Customer#123",
					"sk": "Order#456",
				}
				keys := kv.GetKeyValues(keyMap)
				Expect(len(keys)).To(BeEquivalentTo(2))
				Expect(keys[0]).To(BeEquivalentTo("Customer#123"))
				Expect(keys[1]).To(BeEquivalentTo("Order#456"))
			})
		})
	})

	When("GetValueEndCode", func() {
		It("should get next value", func() {
			endCode := kv.GetEndRangeValue("Customer#")
			Expect(endCode).NotTo(BeNil())
			Expect(endCode).To(BeEquivalentTo("Customer$"))
		})
	})

	When("ValidateExpression", func() {
		When("expression is valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "==", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("expressions empty", func() {
			It("should be valid", func() {
				err := kv.ValidateExpressions([]sdk.QueryExpression{})
				Expect(err).To(BeNil())
			})
		})
		When("operand is nil", func() {
			It("should return error", func() {
				err := kv.ValidateExpressions(nil)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operand is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "", Operator: "==", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operator is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("value is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "==", Value: ""},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "=", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startWith", Value: "Order#"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("primary key operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: ">", Value: "Customer#1000"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
