package main

import (
	"testing"
	"io/ioutil"
	"github.com/JesusIslam/tldr-lr"
)

const (
	num = 3
)

var result string

func TestSummarizeCentralityJaccard(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "centrality", tldr.WEIGHING)
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func TestSummarizePageRankJaccard(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", tldr.WEIGHING)
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func TestSummarizeCentralityHamming(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "hamming")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func TestSummarizePageRankHamming(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "hamming")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func TestSummarizeCentralityTfidf(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "tfidf")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func TestSummarizePageRankTfidf(t *testing.T) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "pagerank", "tfidf")
	bag := tldr.New()
	result = bag.Summarize(text, 3)
}

func BenchmarkSummarizeCentralityJaccard(b *testing.B) {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "centrality", tldr.WEIGHING)
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankJaccard(b *testing.B) {
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
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "centrality", "hamming")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankHamming(b *testing.B) {
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
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	var r string
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "centrality", "tfidf")
	for n := 0; n < b.N; n++ {
		bag := tldr.New()
		r = bag.Summarize(text, 3)
	}
	result = r
}

func BenchmarkSummarizePageRankTfIdf(b *testing.B) {
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