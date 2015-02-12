# tldr
### When you are too lazy to read the entire text
------------------------------------------------------
[![GoDoc](https://godoc.org/github.com/JesusIslam/tldr?status.svg)](https://godoc.org/github.com/JesusIslam/tldr)

### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have three weighing and two ranking algorithm included, they are tfidf-modified-cosine distance, Jaccard coeficient, Hamming distance, and PageRank, centrality, respectively. The default settings use Hamming distance and centrality.

There is now new weighing algorithm: `byteferret` that will generate same summary as Jaccard and Hamming, but is slower than both of them.

The best combination that produced the best summaries are Hamming and centrality in my opinion (and tests, which is why they are the default.)

### Is this fast?
Test it yourself using `go text -bench . -cpu 4 -benchtime=5s`. 
My system has i3-3217U @1.8Ghz and this is the result (Windows, with some programs going on) :
```
BenchmarkSummarizeCentralityJaccard-4       1000           6649430 ns/op
BenchmarkSummarizePageRankJaccard-4         1000           6483326 ns/op
BenchmarkSummarizeCentralityHamming-4       1000           6469321 ns/op
BenchmarkSummarizePageRankHamming-4         1000           6464313 ns/op
BenchmarkSummarizeCentralityTfidf-4          200          49638106 ns/op
BenchmarkSummarizePageRankTfIdf-4            200          50038370 ns/op
BenchmarkSummarizeCentralityFerret-4         300          20933952 ns/op
BenchmarkSummarizePageRankFerret-4           300          20833901 ns/op
``` 

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
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	bag := tldr.New()
	result := bag.Summarize(text, intoSentences)
	fmt.Println(result)
}
```

Or just run sample.go in the test directory.

### Dependencies?
tldr depends on [pagerank](https://github.com/dcadenas/pagerank) package, and you can install it with `go get github.com/dcadenas/pagerank`.

### License?
Check the LICENSE file. tldr: MIT.

## Have fun!