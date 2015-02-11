package main

import (
	"testing"
	"io/ioutil"
	"github.com/JesusIslam/tldr"
)

const (
	num = 3
)

var result string

func TestSummarizeCentralityJaccard(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, tldr.WEIGHING)
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("Centrality Jaccard test failed")
	}
}

func TestSummarizePageRankJaccard(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", tldr.WEIGHING)
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("PageRank Jaccard test failed")
	}
}

func TestSummarizeCentralityHamming(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "hamming")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("Centrality Hamming test failed")
	}
}

func TestSummarizePageRankHamming(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "hamming")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("PageRank Hamming test failed")
	}
}

func TestSummarizeCentralityTfidf(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "tfidf")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("Centrality tfidf test failed")
	}
}

func TestSummarizePageRankTfidf(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "tfidf")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("PageRank tfidf test failed")
	}
}

func TestSummarizeCentralityByteFerret(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "ferret")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("Centrality byteferret test failed")
	}
}

func TestSummarizePageRankByteFerret(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "ferret")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("PageRank byteferret test failed")
	}
}

func TestSummarizeCentralityJaroWinkler(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "jarowinkler")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("Centrality Jaro Winkler test failed")
	}
}

func TestSummarizePageRankJaroWinkler(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "jarowinkler")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
	if len(result) <= 0 {
		t.Fatal("PageRank Jaro Winkler test failed")
	}
}

func BenchmarkSummarizeCentralityJaccard(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, tldr.WEIGHING)
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankJaccard(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", tldr.WEIGHING)
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizeCentralityHamming(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "hamming")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankHamming(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "hamming")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizeCentralityTfidf(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "tfidf")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankTfIdf(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "tfidf")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizeCentralityFerret(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "ferret")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankFerret(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "ferret")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizeCentralityJaroWinkler(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "jarowinkler")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankJaroWinkler(b *testing.B) {
	b.ResetTimer()
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "jarowinkler")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}