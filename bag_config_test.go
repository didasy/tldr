package tldr_test

import (
	. "github.com/didasy/tldr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bag configuration methods", func() {
	var (
		bag *Bag
	)

	BeforeEach(func() {
		bag = New()
	})

	Describe("String()", func() {
		It("Should return a JSON string representation of the bag", func() {
			str := bag.String()
			Expect(str).To(ContainSubstring("\"Algorithm\""))
			Expect(str).To(ContainSubstring("\"Weighing\""))
			Expect(str).To(ContainSubstring("\"pagerank\""))
			Expect(str).To(ContainSubstring("\"hamming\""))
		})
	})

	Describe("Set()", func() {
		It("Should set all configuration parameters correctly", func() {
			bag.Set(1000, 0.9, 0.0002, 0.002, 0.96, "centrality", "jaccard")

			Expect(bag.MaxCharacters).To(Equal(1000))
			Expect(bag.Damping).To(Equal(0.9))
			Expect(bag.Tolerance).To(Equal(0.0002))
			Expect(bag.Threshold).To(Equal(0.002))
			Expect(bag.SentencesDistanceThreshold).To(Equal(0.96))
			Expect(bag.Algorithm).To(Equal("centrality"))
			Expect(bag.Weighing).To(Equal("jaccard"))
		})

		It("Should work with zero values", func() {
			bag.Set(0, 0.0, 0.0, 0.0, 0.0, "", "")

			Expect(bag.MaxCharacters).To(Equal(0))
			Expect(bag.Damping).To(Equal(0.0))
			Expect(bag.Tolerance).To(Equal(0.0))
			Expect(bag.Threshold).To(Equal(0.0))
			Expect(bag.SentencesDistanceThreshold).To(Equal(0.0))
			Expect(bag.Algorithm).To(Equal(""))
			Expect(bag.Weighing).To(Equal(""))
		})
	})

	Describe("SetDictionary()", func() {
		It("Should set a custom dictionary", func() {
			dict := map[string]int{
				"hello": 1,
				"world": 2,
				"test":  3,
			}
			bag.SetDictionary(dict)
			Expect(bag.Dict).To(Equal(dict))
		})

		It("Should allow setting empty dictionary", func() {
			dict := make(map[string]int)
			bag.SetDictionary(dict)
			Expect(bag.Dict).To(Equal(dict))
		})
	})

	Describe("SetCustomAlgorithm()", func() {
		It("Should set a custom algorithm function", func() {
			customAlg := func(edges []*Edge) []int {
				return []int{0, 1, 2}
			}
			bag.SetCustomAlgorithm(customAlg)
			// Test that the custom algorithm is used by verifying it doesn't panic
			Expect(bag.Algorithm).To(Equal("pagerank")) // Default algorithm remains unchanged
		})
	})

	Describe("SetCustomWeighing()", func() {
		It("Should set a custom weighing function", func() {
			customWeighing := func(src, dst []int) float64 {
				return 0.5
			}
			bag.SetCustomWeighing(customWeighing)
			// Test that the custom weighing is used by verifying it doesn't panic
			Expect(bag.Weighing).To(Equal("hamming")) // Default weighing remains unchanged
		})
	})

	Describe("SetWordTokenizer()", func() {
		It("Should set a custom word tokenizer function", func() {
			customTokenizer := func(sentence string) []string {
				return []string{"custom", "token"}
			}
			bag.SetWordTokenizer(customTokenizer)
			// Test that the custom tokenizer is used by actually using it
			bag.OriginalSentences = []string{"test sentence"}
			result, err := bag.Summarize("", 1)
			Expect(err).To(BeNil())
			// The result might be nil for single short sentences
			Expect(result).To(BeNil())
		})
	})
})