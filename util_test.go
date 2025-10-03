package tldr

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utility functions", func() {
	Describe("SymmetricDifference", func() {
		Context("With different slices", func() {
			It("Should return indices where values differ", func() {
				src := []int{1, 2, 3, 4, 5}
				dst := []int{1, 0, 3, 0, 5}
				result := SymmetricDifference(src, dst)
				Expect(result).To(Equal([]int{1, 3}))
			})
		})

		Context("With identical slices", func() {
			It("Should return empty slice", func() {
				src := []int{1, 2, 3}
				dst := []int{1, 2, 3}
				result := SymmetricDifference(src, dst)
				Expect(result).To(BeEmpty())
			})
		})

		Context("With completely different slices", func() {
			It("Should return all indices", func() {
				src := []int{1, 2, 3}
				dst := []int{0, 0, 0}
				result := SymmetricDifference(src, dst)
				Expect(result).To(Equal([]int{0, 1, 2}))
			})
		})

		Context("With empty slices", func() {
			It("Should return empty slice", func() {
				src := []int{}
				dst := []int{}
				result := SymmetricDifference(src, dst)
				Expect(result).To(BeEmpty())
			})
		})

		Context("With single element slices", func() {
			It("Should handle single element comparison", func() {
				src := []int{1}
				dst := []int{0}
				result := SymmetricDifference(src, dst)
				Expect(result).To(Equal([]int{0}))
			})
		})
	})

	Describe("Intersection", func() {
		Context("With some matching elements", func() {
			It("Should return indices where values match", func() {
				src := []int{1, 2, 3, 4, 5}
				dst := []int{1, 0, 3, 0, 5}
				result := Intersection(src, dst)
				Expect(result).To(Equal([]int{0, 2, 4}))
			})
		})

		Context("With no matching elements", func() {
			It("Should return empty slice", func() {
				src := []int{1, 2, 3}
				dst := []int{0, 0, 0}
				result := Intersection(src, dst)
				Expect(result).To(BeEmpty())
			})
		})

		Context("With all matching elements", func() {
			It("Should return all indices", func() {
				src := []int{1, 2, 3}
				dst := []int{1, 2, 3}
				result := Intersection(src, dst)
				Expect(result).To(Equal([]int{0, 1, 2}))
			})
		})

		Context("With empty slices", func() {
			It("Should return empty slice", func() {
				src := []int{}
				dst := []int{}
				result := Intersection(src, dst)
				Expect(result).To(BeEmpty())
			})
		})

		Context("With single element slices", func() {
			It("Should handle single element comparison", func() {
				src := []int{1}
				dst := []int{1}
				result := Intersection(src, dst)
				Expect(result).To(Equal([]int{0}))
			})
		})

		Context("With larger slices to test capacity optimization", func() {
			It("Should handle larger slices efficiently", func() {
				src := make([]int, 100)
				dst := make([]int, 100)
				for i := 0; i < 100; i++ {
					src[i] = i
					dst[i] = i % 2 // Every other position matches when i is even
				}
				result := Intersection(src, dst)
				// Positions where src[i] == dst[i] (only when i=0 and i=1)
				expected := []int{0, 1}
				Expect(result).To(Equal(expected))
			})
		})
	})
})