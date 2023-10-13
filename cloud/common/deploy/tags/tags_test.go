// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tags_test

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tags", func() {
	Context("Tags", func() {
		When("Calling Tags", func() {
			resourcesTags := tags.Tags("testing", "my-resource", resources.Stack)

			It("tags should contain a name tag", func() {
				Expect(resourcesTags["x-nitric-testing-name"]).To(Equal("my-resource"))
			})

			It("tags should contain a type tag", func() {
				Expect(resourcesTags["x-nitric-testing-type"]).To(Equal("stack"))
			})
		})
	})
})
