/*
Dependencies:
  go get github.com/dcadenas/pagerank
*/
package tldr

import (
	"sort"
	"math"
	"bytes"
	"strings"
	"unicode"
	"github.com/dcadenas/pagerank"
)

type Bag struct {
	sentences [][]string
	originalSentences []string
	dict map[string]int
	nodes []*Node
	edges []*Edge
	ranks []int
}

// Create new summarizer
func New() *Bag {
	return &Bag{}
}

// the default values of each settings
const (
	VERSION = "0.1.3"
	ALGORITHM = "centrality"
	WEIGHING = "hamming"
	DAMPING = 0.85
	TOLERANCE = 0.0001
	THRESHOLD = 0.001
)

var (
	Algorithm string = "centrality"
	Weighing string = "hamming"
	Damping float64 = 0.85
	Tolerance float64 = 0.0001
	Threshold float64 = 0.001
)

// Set damping, tolerance, threshold, algorithm, and weighing
func Set(d, t, th float64, alg, w string) {
	Damping = d
	Tolerance = t
	Threshold = th
	Algorithm = alg
	Weighing = w
}

// Summarize the text to num sentences
func (bag *Bag) Summarize(text string, num int) string {
	createOriginalSentencesChan := make(chan bool)
	createSentencesChan := make(chan bool)
	defer close(createOriginalSentencesChan)
	defer close(createSentencesChan)
	go bag.createOriginalSentences(text, createOriginalSentencesChan)
	go bag.createSentences(text, createSentencesChan)
	<- createOriginalSentencesChan
	<- createSentencesChan
	if Weighing == "tfidf" {
		bag.createTfIdfModifiedCosineSimilarityEdges()
	} else if Weighing == "jarowinkler" {
		bag.createJaroWinklerEdges()
	} else if Weighing == "ferret" {
		bag.createByteFerretEdges()
	} else {
		bag.createDictionary(text)
		bag.createNodes()
		bag.createEdges()
	}
	if Algorithm == "centrality" {
		bag.centrality()
	} else if Algorithm == "pagerank" {
		bag.pageRank()
	} else {
		bag.centrality()
	}
	// get only num top of idx
	idx := bag.ranks[:num]
	// sort it ascending
	sort.Ints(idx)
	var res string
	for i := 0; i < len(idx); i++ {
		res += bag.originalSentences[idx[i]] + " "
	}
	return res
}

func (bag *Bag) centrality() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.edges {
		if edge.weight > Threshold {
			newEdges = append(newEdges, edge)
		}
	}
	// sort them by weight descending, using insertion sort
	for i, v := range newEdges {
		j := i - 1
		for j >= 0 && newEdges[j].weight < v.weight {
			newEdges[j+1] = newEdges[j]
			j -= 1
		}
		newEdges[j+1] = v
	}
	rankBySrc := make([]int, len(newEdges))
	for i, v := range newEdges {
		rankBySrc[i] = v.src
	}
	// uniq it without disturbing the order
	m := make(map[int]bool)
	var uniq []int
	for _, v := range rankBySrc {
		if m[v] {
			continue
		}
		uniq = append(uniq, v)
		m[v] = true
	}
	bag.ranks = uniq
}

type Rank struct {
	idx int
	score float64
}

func (bag *Bag) pageRank() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.edges {
		if edge.weight > Threshold {
			newEdges = append(newEdges, edge)
		}
	}
	// then page rank them
	graph := pagerank.New()
	defer graph.Clear()
	for _, edge := range newEdges {
		graph.Link(edge.src, edge.dst)
	}
	var ranks []*Rank
	graph.Rank(Damping, Tolerance, func (sentenceIndex int, rank float64) {
		ranks = append(ranks, &Rank{sentenceIndex, rank})
	})
	// sort ranks into an array of sentence index, by rank descending
	for i, v := range ranks {
		j := i - 1
		for j >= 0 && ranks[j].score < v.score {
			ranks[j+1] = ranks[j]
			j -= 1
		}
		ranks[j+1] = v
	}
	idx := make([]int, len(ranks))
	for i, v := range ranks {
		idx[i] = v.idx
	}

	bag.ranks = idx
}

type Edge struct {
	src int // index of node
	dst int // index of node
	weight float64 // weight of the similarity between two sentences
}

// this is experimental, using triple bytes with ferret algorithm
func (bag *Bag) createByteFerretEdges() {
	for i, srcT := range bag.sentences {
		for j, dstT := range bag.sentences {
			if i != j {
				src := strings.Join(srcT, " ")
				dst := strings.Join(dstT, " ")
				weight := 0.0
				srcB := []byte(src)
				dstB := []byte(dst)
				var lesser []byte
				if len(srcB) < len(dstB) {
					lesser = srcB
				} else {
					lesser = dstB
				}
				lesserLen := len(lesser)
				shingleLen := 3
				exists := false
				for i := 0; i < lesserLen-shingleLen+1; i++ {
					exists = bytes.Contains(srcB, dstB[i:(i+shingleLen)])
					if exists {
						weight++
					}
				}
				// Jaccard distance
				weight = 1.0 - weight / ( ( float64( len(srcB) ) + float64( len(dstB) ) ) - weight )
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}

// this is experimental, using sanitized sentence strings but not tokenized
func (bag *Bag) createJaroWinklerEdges() {
	for i, srcT := range bag.sentences {
		for j, dstT := range bag.sentences {
			if i != j {
				src := strings.Join(srcT, " ")
				dst := strings.Join(dstT, " ")
				weight := createJaroWinklerDistance(src, dst)
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}

/*
Whole code adapted to Go from:
https://github.com/NaturalNode/natural/blob/master/lib/natural/distance/jaro-winkler_distance.js
*/
func distance(str1 string, str2 string) float64 {
	if len(str1) == 0 && len(str2) == 0 {
		return 0.0
	}
	if str1 == str2 {
		return 1.0
	}
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)

	// s1 is lesser, s2 is higher
	var s1, s2 string
	if len(str1) <= len(str2) {
		s1 = str1
		s2 = str2
	} else {
		s1 = str2
		s2 = str1
	}

	matchWindow := int(math.Floor(math.Max(float64(len(s1)), float64(len(s2))) / 2.0) - 1.0)
	matches1 := make([]bool, len(s1))
	matches2 := make([]bool, len(s2))
	var m float64
	var t float64

	for i, v := range s1 {
		matched := false
		if v == rune(s2[i]) {
			matches1[i] = true
			matches2[i] = true
			matched = true
			m++
		} else {
			var k int
			if i <= matchWindow {
				k = 0
			} else {
				k = i - matchWindow
			}
			for {
				// Guard so we would not call uninitialized index
				var x int
				dif := len(s2) - len(s1)
				if dif < 2 {
					x = 0
				} else {
					x = (dif - 2)
				}
				if k == (len(s2) - x) {
					break
				}
				//
				if v == rune(s2[k]) {
					if !matches1[i] && !matches2[k] {
						m++
					}
					matches1[i] = true
					matches2[k] = true
					matched = true
				}
				k++
				if (k <= (i + matchWindow)) && k < len(s2) && matched {
					break
				}
			}
		}
	}

	if m == 0 {
		return 0.0
	}

	k := 0
	for _, v := range s1 {
		// guard from index out of range error
		if k >= len(matches1) - 1 {
			break
		}
		//
		if matches1[k] {
			for k < len(matches2) && !matches2[k] {
				k++
			}
			if k < len(matches2) && v != rune(s2[k]) {
				t++
			}
			k++
		}
	}

	t = t / 2.0
	x1 := m / float64(len(s1))
	x2 := m / float64(len(s2))
	return (x1 + x2 + ((m - t) / m) ) / 3
}

func createJaroWinklerDistance(s1 string, s2 string) float64 {
	if s1 == s2 {
		return 1
	}
	d := distance(s1, s2)
	p := 0.1
	l := 0
	for s1[l] == s2[l] && l < 4 {
		l++
	}
	return d + float64(l) * p * (1 - d)
}

func (bag *Bag) createTfIdfModifiedCosineSimilarityEdges() {
	// find tf of each word in each sentence in the bag.sentences
	for i, src := range bag.originalSentences {
		for j, dst := range bag.originalSentences {
			if i != j {
				var weight float64
				srcS := bag.sentences[i]
				dstS := bag.sentences[j]
				// create a dict
				bag.createDictionary(src + " " + dst)
				// transform to seq dict
				seqDict := createSeqDict(bag.dict)
				// find tf, idf and tfidf of their words
				srcDoneChan := make(chan []float64)
				dstDoneChan := make(chan []float64)
				idfDoneChan := make(chan []float64)
				defer close(srcDoneChan)
				defer close(dstDoneChan)
				defer close(idfDoneChan)
				go createTfVector(srcS, seqDict, srcDoneChan)
				go createTfVector(dstS, seqDict, dstDoneChan)
				go createIdf(srcS, dstS, seqDict, idfDoneChan)
				srcTfVector := <- srcDoneChan
				dstTfVector := <- dstDoneChan
				idf := <- idfDoneChan
				go createTfIdfVector(srcTfVector, idf, srcDoneChan)
				go createTfIdfVector(dstTfVector, idf, dstDoneChan)
				srcTfIdfVector := <- srcDoneChan
				dstTfIdfVector := <- dstDoneChan
				// https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/ for more explanation
				// find the dot-product-idf-modified of them
				dotProductDoneChan := make(chan float64)
				defer close(dotProductDoneChan)
				go createDotProduct(srcTfIdfVector, dstTfIdfVector, dotProductDoneChan)
				// find the sum of magnitude of each tfidf in each sentences
				srcMagnitudeDoneChan := make(chan float64)
				dstMagnitudeDoneChan := make(chan float64)
				defer close(srcMagnitudeDoneChan)
				defer close(dstMagnitudeDoneChan)
				go createMagnitude(srcTfIdfVector, srcMagnitudeDoneChan)
				go createMagnitude(dstTfIdfVector, dstMagnitudeDoneChan)
				// now calculate tf-idf-modified-cosine-similarity between the sentences
				// http://en.wikipedia.org/wiki/Cosine_similarity exactly like this, but switch tf with tfidf
				// http://upload.wikimedia.org/math/f/3/6/f369863aa2814d6e283f859986a1574d.png for the formula
				dotProduct := <- dotProductDoneChan
				srcM := <- srcMagnitudeDoneChan
				dstM := <- dstMagnitudeDoneChan
				weight = dotProduct / (srcM * dstM)
				// put them into the bag
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}

func createTfVector(sen, seqDict []string, done chan<- []float64) {
	vector := make([]float64, len(seqDict))
	for i, term := range seqDict {
		for _, word := range sen {
			if term == word {
				vector[i]++
			}
		}
	}
	// find the tf
	senTermsCount := float64(len(sen))
	for i, count := range vector {
		vector[i] = count / senTermsCount
	}
	done <- vector
}

func createIdf(src, dst, seqDict []string, done chan<- []float64) {
	// count occurence of each word in each sentence
	idf := make([]float64, len(seqDict))
	for i, term := range seqDict {
		for _, word := range src {
			if term == word {
				idf[i] = 1.0
			}
		}
		for _, word := range dst {
			if term == word {
				if idf[i] == 1.0 {
					idf[i] = 2.0
				} else {
					idf[i] = 1.0
				}
			}
		}
	}
	// calculate idf of each words
	for i, count := range idf {
		idf[i] = math.Log(2.0 / count)
	}
	done <- idf
}

func createTfIdfVector(vector, idf []float64, done chan<- []float64) {
	tfidf := make([]float64, len(vector))
	for i, tf := range vector {
		tfidf[i] = tf * idf[i]
	}
	done <- tfidf
}

func createDotProduct(srcVector, dstVector []float64, done chan<- float64) {
	var dotProduct float64
	for i, v := range srcVector {
		dotProduct += (v * dstVector[i])
	}
	done <- dotProduct
}

func createMagnitude(vector []float64, done chan<- float64) {
	var magnitude float64
	for _, tfidf := range vector {
		magnitude += (tfidf * tfidf)
	}
	magnitude = math.Sqrt(magnitude)
	done <- magnitude
}

func createSeqDict(dict map[string]int) []string {
	var seq []string
	for term, _ := range dict {
		seq = append(seq, term)
	}
	return seq
}

func (bag *Bag) createEdges() {
	for i, src := range bag.nodes {
		for j, dst := range bag.nodes {
			// don't compare same node
			if i != j {
				var weight float64
				if Weighing == "jaccard" {
					commonElements := intersection(src.vector, dst.vector)
					// Old version, Jaccard's coeficient, not distance
					// weight = float64(len(commonElements)) / ((float64(vectorLength) * 2) - float64(len(commonElements)))
					weight = 1.0 - float64(len(commonElements)) / ((float64(vectorLength) * 2) - float64(len(commonElements)))
				} else if Weighing == "hamming" {
					differentElements := symetricDifference(src.vector, dst.vector)
					weight = float64(len(differentElements))
				} else {
					commonElements := intersection(src.vector, dst.vector)
					weight = 1.0 - float64(len(commonElements)) / ((float64(vectorLength) * 2) - float64(len(commonElements)))
				}
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}


func symetricDifference(src, dst []int) []int {
	var diff []int
	for i, v := range src {
		if v != dst[i] {
			diff = append(diff, i)
		}
	}
	return diff
}

func intersection(src, dst []int) []int {
	intersect := make(map[int]bool)
	for i, v := range src {
		// Old version, only counting vector value with more than 0. So we only count occurence of a word in both sentence as similarity.
		// if v > 0 && dst[i] > 0 {
		// this one also counting whether if a word doesn't occur on both sentences
		if v == dst[i] {
			intersect[i] = true
		}
	}
	var result []int
	for k, _ := range intersect {
		result = append(result, k)
	}
	return result
}

type Node struct {
	sentenceIndex int // index of sentence from the bag
	vector []int // map of word count in respect with dict, should we use map instead of slice?
	// for example :
	/*
	dict = {
		i : 1
		am : 2
		the : 3
		shit : 4
	}
	str = "I am not shit, you effin shit"
	vector = [1, 1, 0, 2] => [1, 1, 0, 1] because should be binary
	*/
}

var vectorLength int

func (bag *Bag) createNodes() {
	vectorLength = len(bag.dict)
	for i, sentence := range bag.sentences {
		// vector length is len(dict)
		vector := make([]int, vectorLength)
		// word for word now
		for _, word := range sentence {
			// check word dict position, if doesn't exist, skip
			if bag.dict[word] == 0 {
				continue
			}
			// minus 1, because array started from 0 and lowest dict is 1
			pos := bag.dict[word] - 1
			// set 1 to the position
			vector[pos] = 1
		}
		// vector is now created, put it into the node
		node := &Node{i, vector}
		// node is now completed, put into the bag
		bag.nodes = append(bag.nodes, node)
	}
}

func (bag *Bag) createSentences(text string, done chan<- bool) {
	// trim all spaces
	text = strings.TrimSpace(text)
	words := strings.Fields(text)
	var sentence []string
	var sentences [][]string
	for _, word := range words {
		// FIX
		word = strings.ToLower(word)

		// if there isn't . ? or !, append to sentence. If found, also append but reset the sentence
		if strings.ContainsRune(word, '.') || strings.ContainsRune(word, '!') || strings.ContainsRune(word, '?') {
			// FIX
			word = strings.Map(func (r rune) rune {
				if r == '.' || r == '!' || r == '?' {
					return -1
				}
				return r
			}, word)
			// FIX
			sanitizeWord(&word)
			sentence = append(sentence, word)
			sentences = append(sentences, sentence)
			sentence = []string{}
		} else {
			sanitizeWord(&word)
			sentence = append(sentence, word)
		}
	}
	if len(sentence) > 0 {
		sentences = append(sentences, sentence)
	}
	// remove doubled sentence
	uniqSentences(sentences)
	bag.sentences = sentences
	done <- true
}

func (bag *Bag) createOriginalSentences(text string, done chan<- bool) {
	// trim all spaces
	text = strings.TrimSpace(text)
	// tokenize text
	words := strings.Fields(text)
	// build sentence
	var sentence []string
	var sentences [][]string
	for _, word := range words {
		// if there isn't . ? or !, append to sentence. If found, also append but reset the sentence
		if strings.ContainsRune(word, '.') || strings.ContainsRune(word, '!') || strings.ContainsRune(word, '?') {
			sentence = append(sentence, word)
			sentences = append(sentences, sentence)
			sentence = []string{}
		} else {
			sentence = append(sentence, word)
		}
	}
	if len(sentence) > 0 {
		sentences = append(sentences, sentence)
	}
	// now flatten them
	var bagOfSentence []string
	for _, s := range sentences {
		str := strings.Join(s, " ")
		bagOfSentence = append(bagOfSentence, str)
	}
	bag.originalSentences = bagOfSentence
	done <- true
}

func uniqSentences(sentences [][]string) {
	var msens []string
	for _, sen := range sentences {
		merged := strings.Join(sen, " ")
		msens = append(msens, merged)
	}
	// Do jarowinkler then CSIS to deduplicate sentences
	reject := make(map[int]bool, len(msens))
	// First JaroWinkler
	next := 1
	for _, sen := range msens {
		for j := next; j < len(msens); j++ {
			if distance(sen, msens[j]) >= 0.95 {
				reject[j] = true
			}
		}
		next++
	}
	// Then CSIS
	for i, psen := range msens {
		for j, nsen := range msens {
			if i != j {
				// if i subset of j, put i in reject
				if strings.Contains(nsen, psen) {
					reject[i] = true
					continue
				}
				// if j subset of i, put j in reject
				if strings.Contains(psen, nsen) {
					reject[j] = true
					continue
				}
			}
		}
	}

	sentences = [][]string{}
	for i, sen := range msens {
		if reject[i] {
			continue
		}
		sentences = append(sentences, strings.Fields(sen))
	}
}

func sanitizeWord(word *string) {
	*word = strings.ToLower(*word)
	var prev rune
	*word = strings.Map(func (r rune) rune {
		// don't remove '-' if it exists after alphanumerics
		if r == '-' && (unicode.IsDigit(prev) || unicode.IsLetter(prev)) {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return -1
		}
		prev = r
		return r
	}, *word)
}

func (bag *Bag) createDictionary(text string) {
	// trim all spaces
	text = strings.TrimSpace(text)
	// lowercase the text
	text = strings.ToLower(text)
	// remove all non alphanumerics but spaces
	var prev rune
	text = strings.Map(func (r rune) rune {
		if r == '-' && (unicode.IsDigit(prev) || unicode.IsLetter(prev)) {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return -1
		}
		prev = r
		return r
	}, text)
	// turn it into bag of words
	words := strings.Fields(text)
	// turn it into dictionary
	dict := make(map[string]int)
	i := 1
	for _, word := range words {
		if dict[word] == 0 {
			dict[word] = i
			i++
		}
	}
	bag.dict = dict
}