# tldr
### When you are too lazy to read the entire text
------------------------------------------------------
[![GoDoc](https://godoc.org/github.com/JesusIslam/tldr?status.svg)](https://godoc.org/github.com/JesusIslam/tldr)

### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have three weighing and two ranking algorithm included, they are Jaccard coeficient, Hamming distance, and PageRank, centrality, respectively. The default settings use Hamming distance and centrality.

### Is This Fast?
Test it yourself, my system is i3-3217@1.8GHz with single channel 4GB RAM using Ubuntu 15.10 with kernel 4.5.0
```
$ go test -bench . -benchmem -benchtime 5s -cpu 4
BenchmarkSummarizeCentralityHamming-4	    1000	   6583990 ns/op	  401311 B/op	    3549 allocs/op
BenchmarkSummarizeCentralityJaccard-4	     300	  28813583 ns/op	 3449792 B/op	   12546 allocs/op
BenchmarkSummarizePagerankHamming-4  	    1000	   7122519 ns/op	  420581 B/op	    3729 allocs/op
BenchmarkSummarizePagerankJaccard-4  	     300	  30523688 ns/op	 3468906 B/op	   12725 allocs/op
ok  	github.com/JesusIslam/tldr	38.566s
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
	result := bag.Summarize(text, intoSentences)
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