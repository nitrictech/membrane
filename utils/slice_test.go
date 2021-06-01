package utils_test

import (
	"github.com/nitric-dev/membrane/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slice", func() {
	defer GinkgoRecover()

	Context("IndexOf", func() {
		slice := []string{"A", "B", "C"}
		When("Found First", func() {
			It("Should return 0", func() {
				index := utils.IndexOf(slice, "A")
				Expect(index).To(BeIdenticalTo(0))
			})
		})
		When("Found Second", func() {
			It("Should return 1", func() {
				index := utils.IndexOf(slice, "B")
				Expect(index).To(BeIdenticalTo(1))
			})
		})
		When("Found Third", func() {
			It("Should return 2", func() {
				index := utils.IndexOf(slice, "C")
				Expect(index).To(BeIdenticalTo(2))
			})
		})
		When("Not Found", func() {
			It("Should return -1", func() {
				index := utils.IndexOf(slice, "D")
				Expect(index).To(BeIdenticalTo(-1))
			})
		})
		When("Nil Slice", func() {
			It("Should return -1", func() {
				index := utils.IndexOf(nil, "")
				Expect(index).To(BeIdenticalTo(-1))
			})
		})
	})
	Context("Remove", func() {
		When("nil slice", func() {
			It("Should return nil", func() {
				results := utils.Remove(nil, 1)
				Expect(results).To(BeNil())
			})
		})
		When("empty slice", func() {
			It("Should return empty slice", func() {
				results := utils.Remove([]string{}, 1)
				Expect(results).To(HaveLen(0))
			})
		})
		When("index < 0", func() {
			It("Should return original slice", func() {
				slice := []string{"A", "B", "C"}
				results := utils.Remove(slice, -1)
				Expect(results).To(HaveLen(3))
			})
		})
		When("index > length", func() {
			It("Should return original slice", func() {
				slice := []string{"A", "B", "C"}
				results := utils.Remove(slice, 3)
				Expect(results).To(HaveLen(3))
			})
		})
		When("index = 0", func() {
			It("Should return modified slice", func() {
				slice := []string{"A", "B", "C"}
				results := utils.Remove(slice, 0)
				Expect(results).To(HaveLen(2))
				Expect(results).To(BeEquivalentTo([]string{"B", "C"}))
			})
		})
		When("index = 2", func() {
			It("Should return modified slice", func() {
				slice := []string{"A", "B", "C"}
				results := utils.Remove(slice, 2)
				Expect(results).To(HaveLen(2))
				Expect(results).To(BeEquivalentTo([]string{"A", "B"}))
			})
		})
	})
})
