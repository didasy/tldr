package tldr

import (
	"testing"
)

// TestSentences is more an illustration of how to pass sentences in, rather than a true test
func TestSentences(t *testing.T) {
	bag := New()
	bag.OriginalSentences = []string{
		"Mary had a little lamb,",
		"it's fleece was white as snow,",
		"and everywhere that Mary went,",
		"that lamb was sure to go.",
	}
	result, err := bag.Summarize("", 1)
	if err != nil {
		t.Error(err)
		return
	}
	if result != "it's fleece was white as snow," {
		t.Error("result not as expected")
	}
}
