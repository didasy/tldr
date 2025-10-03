# tldr
### When you are too lazy to read the entire text
------------------------------------------------------
[![Build Status](https://travis-ci.org/JesusIslam/tldr.svg?branch=master)](https://travis-ci.org/JesusIslam/tldr)
[![Coverage Status](https://coveralls.io/repos/github/JesusIslam/tldr/badge.svg?branch=master)](https://coveralls.io/github/JesusIslam/tldr?branch=master)
[![GoDoc](https://godoc.org/github.com/didasy/tldr?status.svg)](https://godoc.org/github.com/didasy/tldr)
[![Go Report Card](https://goreportcard.com/badge/github.com/didasy/tldr)](https://goreportcard.com/report/github.com/didasy/tldr)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FJesusIslam%2Ftldr.svg?type=small)](https://app.fossa.io/projects/git%2Bgithub.com%2FJesusIslam%2Ftldr?ref=badge_small)


### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have two weighing and two ranking algorithm included, they are Jaccard coeficient and Hamming distance, then PageRank and centrality, respectively. The default settings use Hamming distance and pagerank.

### Is This Fast?
```
$ go test -bench . -benchmem -benchtime 5s -cpu 4
Running Suite: Tldr Suite
=========================
Random Seed: 1759506557
Will run 8 of 8 specs

••••••••
Ran 8 of 8 Specs in 0.012 seconds
SUCCESS! -- 8 Passed | 0 Failed | 0 Pending | 0 Skipped
goos: linux
goarch: amd64
pkg: github.com/didasy/tldr
cpu: AMD Ryzen 5 5600G with Radeon Graphics
BenchmarkSummarizeCentralityHamming-4               5877            896338 ns/op          177320 B/op       1898 allocs/op
BenchmarkSummarizeCentralityJaccard-4               6562            885374 ns/op          177221 B/op       1898 allocs/op
BenchmarkSummarizePagerankHamming-4                 5832            962000 ns/op          200830 B/op       2086 allocs/op
BenchmarkSummarizePagerankJaccard-4                 5949            962579 ns/op          200865 B/op       2086 allocs/op
PASS
ok      github.com/didasy/tldr  22.840s
```
So, not bad huh?

### Installation
`go get github.com/didasy/tldr`

### Example

```
package main

import (
	"fmt"
	"io/ioutil"
	"github.com/didasy/tldr"
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
