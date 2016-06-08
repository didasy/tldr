/*
Dependencies:
  go get github.com/dcadenas/pagerank
*/
package tldr

import (
	"errors"
	"github.com/dcadenas/pagerank"
	"math"
	"sort"
	"strings"
	"unicode"
)

type Bag struct {
	BagOfWordsPerSentence [][]string
	OriginalSentences     []string
	Dict                  map[string]int
	Nodes                 []*Node
	Edges                 []*Edge
	Ranks                 []int
}

// Create new summarizer
func New() *Bag {
	return &Bag{}
}

// The default values of each settings
const (
	VERSION           = "0.3.1"
	DEFAULT_ALGORITHM = "centrality"
	DEFAULT_WEIGHING  = "hamming"
	DEFAULT_DAMPING   = 0.85
	DEFAULT_TOLERANCE = 0.0001
	DEFAULT_THRESHOLD = 0.001
)

var (
	Algorithm string  = "centrality"
	Weighing  string  = "hamming"
	Damping   float64 = 0.85
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
func (bag *Bag) Summarize(text string, num int) (string, error) {
	createSentencesChan := make(chan bool)
	createBagOfWordsPerSentenceChan := make(chan bool)
	defer close(createSentencesChan)
	defer close(createBagOfWordsPerSentenceChan)

	go bag.CreateBagOfWordsPerSentence(text, createBagOfWordsPerSentenceChan)
	go bag.CreateSentences(text, createSentencesChan)
	<-createSentencesChan
	<-createBagOfWordsPerSentenceChan

	bag.CreateDictionary(text)
	bag.CreateNodes()
	bag.CreateEdges()

	switch Algorithm {
	case "centrality":
		bag.Centrality()
		break
	case "pagerank":
		bag.PageRank()
		break
	default:
		bag.Centrality()
	}

	// if no ranks, return error
	if len(bag.Ranks) == 0 {
		return "", errors.New("Ranks is empty")
	}

	// guard so it won't crash but return only the highest rank sentence
	// if num is invalid
	if num > (len(bag.Ranks)-1) || num < 1 {
		num = 1
	}
	// get only top num of ranks
	idx := bag.Ranks[:num]
	// sort it ascending by how the sentences appeared on the original text
	sort.Ints(idx)
	var res string
	for i := 0; i < len(idx); i++ {
		res += (bag.OriginalSentences[idx[i]] + " ")
		res += "\n"
	}

	// trim it from spaces
	res = strings.TrimSpace(res)

	return res, nil
}

func (bag *Bag) Centrality() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
		if edge.weight > Threshold {
			newEdges = append(newEdges, edge)
		}
	}
	// sort them by weight descending
	HeapSortEdge(newEdges)
	ReverseEdge(newEdges)
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
	bag.Ranks = uniq
}

type Rank struct {
	idx   int
	score float64
}

func (bag *Bag) PageRank() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
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
	graph.Rank(Damping, Tolerance, func(sentenceIndex int, rank float64) {
		ranks = append(ranks, &Rank{sentenceIndex, rank})
	})
	// sort ranks into an array of sentence index, by rank descending
	HeapSortRank(ranks)
	ReverseRank(ranks)
	idx := make([]int, len(ranks))
	for i, v := range ranks {
		idx[i] = v.idx
	}

	bag.Ranks = idx
}

type Edge struct {
	src    int     // index of node
	dst    int     // index of node
	weight float64 // weight of the similarity between two sentences
}

/*
Whole code adapted to Go from:
https://github.com/NaturalNode/natural/blob/master/lib/natural/distance/jaro-winkler_distance.js
*/
func Distance(str1 string, str2 string) float64 {
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

	matchWindow := int(math.Floor(math.Max(float64(len(s1)), float64(len(s2)))/2.0) - 1.0)
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
		if k >= len(matches1)-1 {
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
	return (x1 + x2 + ((m - t) / m)) / 3
}

func (bag *Bag) CreateEdges() {
	for i, src := range bag.Nodes {
		for j, dst := range bag.Nodes {
			// don't compare same node
			if i != j {
				var weight float64
				switch Weighing {
				case "jaccard":
					commonElements := Intersection(src.vector, dst.vector)
					weight = 1.0 - float64(len(commonElements))/((float64(vectorLength)*2)-float64(len(commonElements)))
					break
				case "hamming":
					differentElements := SymmetricDifference(src.vector, dst.vector)
					weight = float64(len(differentElements))
					break
				default:
					// defaulted to jaccard
					commonElements := Intersection(src.vector, dst.vector)
					weight = 1.0 - float64(len(commonElements))/((float64(vectorLength)*2)-float64(len(commonElements)))
				}
				edge := &Edge{i, j, weight}
				bag.Edges = append(bag.Edges, edge)
			}
		}
	}
}

func SymmetricDifference(src, dst []int) []int {
	var diff []int
	for i, v := range src {
		if v != dst[i] {
			diff = append(diff, i)
		}
	}
	return diff
}

func Intersection(src, dst []int) []int {
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
	sentenceIndex int   // index of sentence from the bag
	vector        []int // map of word count in respect with dict, should we use map instead of slice?
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

func (bag *Bag) CreateNodes() {
	vectorLength = len(bag.Dict)
	for i, sentence := range bag.BagOfWordsPerSentence {
		// vector length is len(dict)
		vector := make([]int, vectorLength)
		// word for word now
		for _, word := range sentence {
			// check word dict position, if doesn't exist, skip
			if bag.Dict[word] == 0 {
				continue
			}
			// minus 1, because array started from 0 and lowest dict is 1
			pos := bag.Dict[word] - 1
			// set 1 to the position
			vector[pos] = 1
		}
		// vector is now created, put it into the node
		node := &Node{i, vector}
		// node is now completed, put into the bag
		bag.Nodes = append(bag.Nodes, node)
	}
}

func (bag *Bag) CreateBagOfWordsPerSentence(text string, done chan<- bool) {
	// trim all spaces
	text = strings.TrimSpace(text)
	words := strings.Fields(text)
	var sentence []string
	var sentences [][]string
	for _, word := range words {
		word = strings.ToLower(word)

		// if there isn't . ? or !, append to sentence. If found, also append but reset the sentence
		if strings.ContainsRune(word, '.') || strings.ContainsRune(word, '!') || strings.ContainsRune(word, '?') {
			word = strings.Map(func(r rune) rune {
				if r == '.' || r == '!' || r == '?' {
					return -1
				}
				return r
			}, word)
			word = SanitizeWord(word)
			sentence = append(sentence, word)
			sentences = append(sentences, sentence)
			sentence = []string{}
		} else {
			word = SanitizeWord(word)
			sentence = append(sentence, word)
		}
	}
	if len(sentence) > 0 {
		sentences = append(sentences, sentence)
	}
	// remove doubled sentence
	UniqSentences(sentences)
	bag.BagOfWordsPerSentence = sentences
	done <- true
}

func (bag *Bag) CreateSentences(text string, done chan<- bool) {
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
	bag.OriginalSentences = bagOfSentence
	done <- true
}

func UniqSentences(sentences [][]string) {
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
			if Distance(sen, msens[j]) >= 0.95 {
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

func SanitizeWord(word string) string {
	word = strings.ToLower(word)
	var prev rune
	word = strings.Map(func(r rune) rune {
		// don't remove '-' if it exists after alphanumerics
		if r == '-' && (unicode.IsDigit(prev) || unicode.IsLetter(prev)) {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return -1
		}
		prev = r
		return r
	}, word)

	return word
}

func (bag *Bag) CreateDictionary(text string) {
	// trim all spaces
	text = strings.TrimSpace(text)
	// lowercase the text
	text = strings.ToLower(text)
	// remove all non alphanumerics but spaces
	var prev rune
	text = strings.Map(func(r rune) rune {
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
	bag.Dict = dict
}
