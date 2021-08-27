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

package secret

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unimplemented Secret Plugin Tests", func() {
	uisp := &UnimplementedSecretPlugin{}

	Context("Put", func() {
		When("Calling Put on UnimplementedSecretPlugin", func() {
			_, err := uisp.Put(nil, nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Access", func() {
		When("Calling Access on UnimplementedSecretPlugin", func() {
			_, err := uisp.Access(nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Secret.String", func() {
		It("should print ReceiveOptions", func() {
			secret := &Secret{Name: "secret"}
			Expect(secret.String()).To(BeEquivalentTo("{Name: secret}"))
		})
	})

	Context("SecretVersion.String", func() {
		It("should print ReceiveOptions", func() {
			secret := &SecretVersion{
				Secret:  &Secret{Name: "secret"},
				Version: "version",
			}
			Expect(secret.String()).To(BeEquivalentTo("{Secret: {Name: secret}, Version: version}"))
		})
	})
})
