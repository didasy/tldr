# tldr
### When you are too lazy to read the entire text
------------------------------------------------------

### What?
tldr is a golang package to summarize a text automatically using [lexrank](http://www.cs.cmu.edu/afs/cs/project/jair/pub/volume22/erkan04a-html/erkan04a.html) algorithm.

### How?
There are two main steps in lexrank, weighing, and ranking. tldr have two weighing and ranking algorithm included, they are Jaccard coeficient, Hamming distance, PageRank, and centrality, respectively. The default settings use Jaccard coeficient and centrality, meanwhile if you are planning to use PageRank, don't be surprised if you got various version of summary of the text (it is PageRank after all.)

### Is this fast?
Test it yourself using `go text -bench .`. The results in my computer (i3-3217U 1.8GHz) are about 2-3ms per operation.

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
	intoSentences := 3s
	textB, _ := ioutil.ReadFile("../sample.txt")
	text := string(textB)
	tldr.Set(tldr.DAMPING, tldr.TOLERANCE, tldr.THRESHOLD, tldr.ALGORITHM, "hamming")
	bag := tldr.New()
	result := bag.Summarize(text, intoSentences)
	fmt.Println(result)
}
```

Or just run sample.go in the test directory.

### Dependencies?
tldr depends on (pagerank)[https://github.com/dcadenas/pagerank] package, and you can install it with `go get github.com/dcadenas/pagerank`.

### License?
Check the LICENSE file. tldr: MIT.

## Have fun!