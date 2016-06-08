package tldr_test

import (
	. "github.com/JesusIslam/tldr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"strings"
)

var (
	err          error
	raw          []byte
	text, result string
	summarizer   *Bag
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
	summarizer = New()
}

var _ = Describe("tldr", func() {
	Describe("Test summarizing", func() {
		Context("Summarize sample.txt to 2 sentences", func() {
			It("Should return a string match with result.txt without error", func() {
				sum, err := summarizer.Summarize(text, 2)
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal(strings.TrimSpace(string(result))))
			})
		})
		Context("Summarize sample.txt to 1 sentence but by giving it invalid parameter", func() {
			It("Should return a string with one sentence without error", func() {
				sum, err := summarizer.Summarize(text, 10000)
				Expect(err).To(BeNil())
				Expect(sum).To(BeAssignableToTypeOf(""))
				Expect(sum).NotTo(BeEmpty())
				Expect(sum).To(Equal("In honor of the Museum of Narrative Art and its star-studded cast of architects, here's a roundup of articles from Architizer that feature Star Wars-related architecture: Jeff Bennett's Wars on Kinkade are hilarious paintings that ravage the peaceful landscapes of Thomas Kinkade with the brutal destruction of Star Wars."))
			})
		})
	})

})
