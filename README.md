# tldr
### When you are too lazy to read the entire text
------------------------------------------------------
[![Build Status](https://travis-ci.org/JesusIslam/tldr.svg?branch=master)](https://travis-ci.org/JesusIslam/tldr)
[![Coverage Status](https://coveralls.io/repos/github/JesusIslam/tldr/badge.svg?branch=master)](https://coveralls.io/github/JesusIslam/tldr?branch=master)
[![GoDoc](https://godoc.org/github.com/JesusIslam/tldr?status.svg)](https://godoc.org/github.com/JesusIslam/tldr)

### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have two weighing and two ranking algorithm included, they are Jaccard coeficient and Hamming distance, then PageRank and centrality, respectively. The default settings use Hamming distance and pagerank.

### Is This Fast?
Test it yourself, my system is i3-3217@1.8GHz with single channel 4GB RAM using Ubuntu 15.10 with kernel 4.5.0
```
$ go test -bench . -benchmem -benchtime 5s -cpu 4
BenchmarkSummarizeCentralityHamming-4	    2000	   6429340 ns/op	  401204 B/op	    3551 allocs/op
BenchmarkSummarizeCentralityJaccard-4	     200	  30036357 ns/op	 3449461 B/op	   12543 allocs/op
BenchmarkSummarizePagerankHamming-4  	    1000	   7015008 ns/op	  420665 B/op	    3731 allocs/op
BenchmarkSummarizePagerankJaccard-4  	     200	  31066764 ns/op	 3469629 B/op	   12737 allocs/op
```
So, not bad huh?

### Installation
`go get github.com/JesusIslam/tldr`

### Example

```
package main

import (
	"fmt"
	"io/ioutil"
	"github.com/JesusIslam/tldr"
)

func main() {
	intoSentences := 3
	textB, _ := ioutil.ReadFile("./sample.txt")
	text := string(textB)
	bag := tldr.New()
	result, _ := bag.Summarize(text, intoSentences)
	fmt.Println(result)
}
```
### Testing
To test, just run `go test`, but you need to have [gomega](http://github.com/onsi/gomega) and [ginkgo](http://github.com/onsi/ginkgo) installed.

### Dependencies?
tldr depends on [pagerank](https://github.com/alixaxel/pagerank) package, and you can install it with `go get github.com/alixaxel/pagerank`.

### License?
Check the LICENSE file. tldr: MIT.

## Have fun!
