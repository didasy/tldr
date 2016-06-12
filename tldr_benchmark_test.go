package tldr_test

import (
	"github.com/JesusIslam/tldr"
	"io/ioutil"
	"testing"
)

const (
	SAMPLE_FILE_PATH = "./sample.txt"
	NUM_SENTENCES    = 3
)

var (
	benchSummarizer *tldr.Bag
	benchText       string
	resultText      string
)

func init() {
	raw, err := ioutil.ReadFile(SAMPLE_FILE_PATH)
	if err != nil {
		panic(err)
	}
	benchText = string(raw)
}

func BenchmarkSummarizeCentralityHamming(b *testing.B) {
	var rtxt string

	for n := 0; n < b.N; n++ {
		benchSummarizer = tldr.New()
		benchSummarizer.Algorithm = "centrality"
		benchSummarizer.Weighing = "hamming"
		rtxt, _ = benchSummarizer.Summarize(benchText, NUM_SENTENCES)
	}

	resultText = rtxt
}

func BenchmarkSummarizeCentralityJaccard(b *testing.B) {
	var rtxt string

	for n := 0; n < b.N; n++ {
		benchSummarizer = tldr.New()
		benchSummarizer.Algorithm = "centrality"
		benchSummarizer.Weighing = "jaccard"
		rtxt, _ = benchSummarizer.Summarize(benchText, NUM_SENTENCES)
	}

	resultText = rtxt
}

func BenchmarkSummarizePagerankHamming(b *testing.B) {
	var rtxt string

	for n := 0; n < b.N; n++ {
		benchSummarizer = tldr.New()
		benchSummarizer.Algorithm = "pagerank"
		benchSummarizer.Weighing = "hamming"
		rtxt, _ = benchSummarizer.Summarize(benchText, NUM_SENTENCES)
	}

	resultText = rtxt
}

func BenchmarkSummarizePagerankJaccard(b *testing.B) {
	var rtxt string

	for n := 0; n < b.N; n++ {
		benchSummarizer = tldr.New()
		benchSummarizer.Algorithm = "pagerank"
		benchSummarizer.Weighing = "jaccard"
		rtxt, _ = benchSummarizer.Summarize(benchText, NUM_SENTENCES)
	}

	resultText = rtxt
}
