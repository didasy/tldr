package tldr_test

import (
	. "github.com/JesusIslam/tldr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"strings"
)

var (
	err                                                                error
	raw                                                                []byte
	text, result, shortResult, resultCentrality, shortResultCentrality string
	summarizer                                                         *Bag
)

func init() {
	raw, err = ioutil.ReadFile("./sample.txt")
	if err != nil {
		panic(err)
	}
	text = string(raw)
	raw, err = ioutil.ReadFile("./result.txt")
	if err != nil {
		panic(err)
	}
	result = string(raw)
	raw, err = ioutil.ReadFile("./short.result.txt")
	if err != nil {
		panic(err)
	}
	shortResult = string(raw)
	raw, err = ioutil.ReadFile("./result_centrality.txt")
	if err != nil {
		panic(err)
	}
	resultCentrality = string(raw)
	raw, err = ioutil.ReadFile("./short.result_centrality.txt")
	if err != nil {
		panic(err)
	}
	shortResultCentrality = string(raw)
}

var _ = Describe("tldr", func() {
	Describe("Test summarizing using default hamming weighing and pagerank algorithm", func() {
		Context("Summarize sample.txt to 3 sentences", func() {
			It("Should return a string match with result.txt without error", func() {
				summarizer = New()
				sums, err := summarizer.Summarize(text, 3)
				summarizer.Algorithm = ""
				summarizer.Weighing = ""
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(result)))
			})
		})
		Context("Summarize sample.txt to 1 sentence but by giving it invalid parameter", func() {
			It("Should return a string with one sentence without error", func() {
				summarizer = New()
				summarizer.Algorithm = ""
				summarizer.Weighing = ""
				sums, err := summarizer.Summarize(text, 10000)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(string(shortResult))))
			})
		})
	})

	Describe("Test summarizing using jaccard weighing and pagerank algorithm", func() {
		Context("Summarize sample.txt to 3 sentences", func() {
			It("Should return a string match with result.txt without error", func() {
				summarizer = New()
				summarizer.Weighing = "jaccard"
				summarizer.Algorithm = ""
				sums, err := summarizer.Summarize(text, 3)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(result)))
			})
		})
		Context("Summarize sample.txt to 1 sentence but by giving it invalid parameter", func() {
			It("Should return a string with one sentence without error", func() {
				summarizer = New()
				summarizer.Weighing = "jaccard"
				summarizer.Algorithm = ""
				sums, err := summarizer.Summarize(text, 10000)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(string(shortResult))))
			})
		})
	})

	Describe("Test summarizing using invalid weighing name and invalid algorithm", func() {
		Context("Summarize sample.txt to 3 sentences", func() {
			It("Should return a string match with result.txt without error", func() {
				summarizer = New()
				summarizer.Weighing = "invalid"
				summarizer.Algorithm = "invalid"
				sums, err := summarizer.Summarize(text, 3)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(result)))
			})
		})
		Context("Summarize sample.txt to 1 sentence but by giving it invalid parameter", func() {
			It("Should return a string with one sentence without error", func() {
				summarizer = New()
				summarizer.Weighing = "invalid"
				summarizer.Algorithm = "invalid"
				sums, err := summarizer.Summarize(text, 10000)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(string(shortResult))))
			})
		})
	})

	Describe("Test summarizing using centrality algorithm", func() {
		Context("Summarize sample.txt to 3 sentences", func() {
			It("Should return a string match with result_centrality.txt without error", func() {
				summarizer = New()
				summarizer.Algorithm = "centrality"
				summarizer.Weighing = "pagerank"
				sums, err := summarizer.Summarize(text, 3)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(resultCentrality)))
			})
		})
		Context("Summarize sample.txt to 1 sentence but by giving it invalid parameter", func() {
			It("Should return a string with one sentence without error", func() {
				summarizer = New()
				summarizer.Algorithm = "centrality"
				summarizer.Weighing = "pagerank"
				sums, err := summarizer.Summarize(text, 10000)
				sum := strings.Join(sums, "\n\n")
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(string(shortResultCentrality))))
			})
		})
	})
})
