package main

import (
	"fmt"
	"io/ioutil"
	"github.com/JesusIslam/tldr-lr"
)

func main() {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, "centrality", "tfidf")
	bag := tldr.New()
	result := bag.Summarize(text, 3)
	fmt.Println(result)
}