# tldr
### When you are too lazy to read the entire text
------------------------------------------------------

### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have three weighing and two ranking algorithm included, they are tfidf-modified-cosine distance, Jaccard coeficient, Hamming distance, and PageRank, centrality, respectively. The default settings use Hamming distance and centrality, meanwhile if you are planning to use PageRank, don't be surprised if you got various version of summary of the text (it is PageRank after all.)

If you want the same exact result as Flipboard auto-summarizer, use tfidf and centrality.

### Is this fast?
Test it yourself using `go text -bench . -cpu 4 -benchtime=5s`. 
My system has i3-3217U @1.8Ghz and this is the result (Windows, with some programs going on) :
```
BenchmarkSummarizeCentralityJaccard-4       2000           4086226 ns/op
BenchmarkSummarizePageRankJaccard-4         2000           4049201 ns/op
BenchmarkSummarizeCentralityHamming-4       2000           4059708 ns/op
BenchmarkSummarizePageRankHamming-4         2000           4050209 ns/op
BenchmarkSummarizeCentralityTfidf-4          100          50493720 ns/op
BenchmarkSummarizePageRankTfIdf-4            100          50473672 ns/op
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