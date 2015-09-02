package tldr_test

import (
	. "github.com/JesusIslam/tldr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"strings"
)

var _ = Describe("Tldr", func() {

	summarizer := New()

	text, err := ioutil.ReadFile("./sample.txt")
	if err != nil {
		Fail("Failed to read sample txt: " + err.Error())
	}

	Describe("Test summarizing", func() {
		Context("Summarize to 3 sentences", func() {
			It("Should returns three sentences string", func() {
				var str string
				sum := summarizer.Summarize(string(text), 3)
				Expect(sum).To(BeAssignableToTypeOf(str))
				Expect(sum).ToNot(BeEmpty())
				Expect(strings.Split(strings.TrimSpace(sum), ".")).To(HaveLen(3))
			})
		})
	})

})
