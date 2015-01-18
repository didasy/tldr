package main

import (
	"fmt"
	"io/ioutil"
	"github.com/JesusIslam/tldr"
)

func main() {
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, tldr.WEIGHING)
	bag := tldr.New()
	result := bag.Summarize(text, 3)
	fmt.Println(result)
}