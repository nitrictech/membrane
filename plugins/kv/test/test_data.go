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
package test

// Simple 'users' collection test data

var UserKey = map[string]interface{}{
	"key": "jsmith@server.com",
}
var UserItem = map[string]interface{}{
	"firstName": "John",
	"lastName":  "Smith",
	"email":     "jsmith@server.com",
}
var UserItem2 = map[string]interface{}{
	"firstName": "Johnson",
	"lastName":  "Smithers",
	"email":     "j.smithers@yahoo.com",
}

// Single Table Design 'application' collection test data

var CustomerKey = map[string]interface{}{
	"pk": "Customer#1000",
	"sk": "Customer#1000",
}
var CustomerItem = map[string]interface{}{
	"testName":  "CustomerItem",
	"firstName": "John",
	"lastName":  "Smith",
	"email":     "jsmith@server.com",
}
var OrderKey1 = map[string]interface{}{
	"pk": "Customer#1000",
	"sk": "Order#501",
}
var OrderItem1 = map[string]interface{}{
	"testName": "OrderItem1",
	"sku":      "ABC-501",
	"number":   "1",
	"price":    "13.95",
}
var OrderKey2 = map[string]interface{}{
	"pk": "Customer#1000",
	"sk": "Order#502",
}
var OrderItem2 = map[string]interface{}{
	"testName": "OrderItem2",
	"sku":      "ABC-502",
	"number":   "1",
	"price":    "13.95",
}
var OrderKey3 = map[string]interface{}{
	"pk": "Customer#1000",
	"sk": "Order#503",
}
var OrderItem3 = map[string]interface{}{
	"testName": "OrderItem3",
	"sku":      "ABC-503",
	"number":   "1",
	"price":    "13.95",
}
var ProductKey = map[string]interface{}{
	"pk": "Product#ABC-502",
	"sk": "Product#ABC-502",
}
var ProductItem = map[string]interface{}{
	"testName": "ProductItem",
	"sku":      "ABC-503",
	"mode":     "dark",
	"weight":   "13.95",
}
