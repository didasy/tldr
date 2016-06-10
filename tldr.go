/*
Dependencies:
  go get github.com/dcadenas/pagerank

WARNING: This package is not thread safe, so you cannot use *Bag from many goroutines.
*/
package tldr

import (
	"errors"
	"github.com/dcadenas/pagerank"
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

	MaxCharacters              int
	Algorithm                  string // "centrality" or "pagerank"
	Weighing                   string // "hamming" or "jaccard"
	Damping                    float64
	Tolerance                  float64
	Threshold                  float64
	SentencesDistanceThreshold float64

	vectorLength int
}

// The default values of each settings
const (
	VERSION                              = "0.4.1"
	DEFAULT_ALGORITHM                    = "centrality"
	DEFAULT_WEIGHING                     = "hamming"
	DEFAULT_DAMPING                      = 0.85
	DEFAULT_TOLERANCE                    = 0.0001
	DEFAULT_THRESHOLD                    = 0.001
	DEFAULT_MAX_CHARACTERS               = 0
	DEFAULT_SENTENCES_DISTANCE_THRESHOLD = 0.95
)

// Create new summarizer
func New() *Bag {
	return &Bag{
		MaxCharacters:              DEFAULT_MAX_CHARACTERS,
		Algorithm:                  DEFAULT_ALGORITHM,
		Weighing:                   DEFAULT_WEIGHING,
		Damping:                    DEFAULT_DAMPING,
		Tolerance:                  DEFAULT_TOLERANCE,
		Threshold:                  DEFAULT_THRESHOLD,
		SentencesDistanceThreshold: DEFAULT_SENTENCES_DISTANCE_THRESHOLD,
	}
}

// Set max characters, damping, tolerance, threshold, sentences distance threshold, algorithm, and weighing
func (bag *Bag) Set(m int, d, t, th, sth float64, alg, w string) {
	bag.MaxCharacters = m
	bag.Damping = d
	bag.Tolerance = t
	bag.Threshold = th
	bag.Algorithm = alg
	bag.Weighing = w
	bag.SentencesDistanceThreshold = sth
}

// Summarize the text to num sentences
func (bag *Bag) Summarize(text string, num int) (string, error) {
	bag.CreateSentences(text)
	bag.CreateDictionary(text)
	bag.CreateNodes()
	bag.CreateEdges()

	switch bag.Algorithm {
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
	HeapSortInt(idx)
	var res string
	for i, _ := range idx {
		res += (bag.OriginalSentences[idx[i]] + " ")
		res += "\n\n"
	}

	// trim it from spaces
	res = strings.TrimSpace(res)

	// Truncate if it has more than n characters
	// Note this is not bytes length
	if bag.MaxCharacters > 0 {
		// turn into runes
		r := []rune(res)
		// cut
		r = r[:bag.MaxCharacters]
		// then turn back to string
		res = string(r)
	}

	return res, nil
}

type Rank struct {
	idx   int
	score float64
}

func (bag *Bag) Centrality() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
		if edge.weight > bag.Threshold {
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

func (bag *Bag) PageRank() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
		if edge.weight > bag.Threshold {
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
	graph.Rank(bag.Damping, bag.Tolerance, func(sentenceIndex int, rank float64) {
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

func (bag *Bag) CreateEdges() {
	for i, src := range bag.Nodes {
		for j, dst := range bag.Nodes {
			// don't compare same node
			if i != j {
				var weight float64
				switch bag.Weighing {
				case "jaccard":
					commonElements := Intersection(src.vector, dst.vector)
					weight = 1.0 - float64(len(commonElements))/((float64(bag.vectorLength)*2)-float64(len(commonElements)))
					break
				case "hamming":
					differentElements := SymmetricDifference(src.vector, dst.vector)
					weight = float64(len(differentElements))
					break
				default:
					// defaulted to jaccard
					commonElements := Intersection(src.vector, dst.vector)
					weight = 1.0 - float64(len(commonElements))/((float64(bag.vectorLength)*2)-float64(len(commonElements)))
				}
				edge := &Edge{i, j, weight}
				bag.Edges = append(bag.Edges, edge)
			}
		}
	}
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

func (bag *Bag) CreateNodes() {
	bag.vectorLength = len(bag.Dict)
	for i, sentence := range bag.BagOfWordsPerSentence {
		// vector length is len(dict)
		vector := make([]int, bag.vectorLength)
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

func (bag *Bag) CreateSentences(text string) {
	// trim all spaces
	text = strings.TrimSpace(text)
	// tokenize text as sentences
	// sentence is a group of words separated by whitespaces or punctuation other than !?.
	bag.OriginalSentences = TokenizeSentences(text)

	// from original sentences, explode each sentences into bag of words
	bag.BagOfWordsPerSentence = [][]string{}
	for _, sentence := range bag.OriginalSentences {
		words := strings.Fields(sentence)
		// then sanitize each word
		for i, word := range words {
			words[i] = SanitizeWord(word)
		}
		bag.BagOfWordsPerSentence = append(bag.BagOfWordsPerSentence, words)
	}

	// then uniq it
	UniqSentences(bag.BagOfWordsPerSentence, bag.SentencesDistanceThreshold)
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
