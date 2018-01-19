/*
Dependencies:
  go get github.com/alixaxel/pagerank

WARNING: This package is not thread safe, so you cannot use *Bag from many goroutines.
*/
package tldr

import (
	"errors"
	"sort"
	"strings"
	"unicode"

	"github.com/alixaxel/pagerank"
)

type Bag struct {
	BagOfWordsPerSentence [][]string
	OriginalSentences     []string
	Dict                  map[string]int
	Nodes                 []*Node
	Edges                 []*Edge
	Ranks                 []int

	MaxCharacters              int
	Algorithm                  string // "centrality" or "pagerank" or "custom"
	Weighing                   string // "hamming" or "jaccard" or "custom"
	Damping                    float64
	Tolerance                  float64
	Threshold                  float64
	SentencesDistanceThreshold float64

	customAlgorithm func(e []*Edge) []int
	customWeighing  func(src, dst []int) float64
	wordTokenizer   func(sentence string) []string

	vectorLength int
}

// The default values of each settings
const (
	VERSION                              = "0.5.0"
	DEFAULT_ALGORITHM                    = "pagerank"
	DEFAULT_WEIGHING                     = "hamming"
	DEFAULT_DAMPING                      = 0.85
	DEFAULT_TOLERANCE                    = 0.0001
	DEFAULT_THRESHOLD                    = 0.001
	DEFAULT_MAX_CHARACTERS               = 0
	DEFAULT_SENTENCES_DISTANCE_THRESHOLD = 0.95
)

func defaultWordTokenizer(sentence string) []string {
	words := strings.Fields(sentence)
	for i, word := range words {
		words[i] = SanitizeWord(word)
	}

	return words
}

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
		wordTokenizer:              defaultWordTokenizer,
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

// Useful if you already have your own dictionary (example: from your database)
// Dictionary is a map[string]int where the key is the word and int is the position in vector, starting from 1
func (bag *Bag) SetDictionary(dict map[string]int) {
	bag.Dict = dict
}

func (bag *Bag) SetCustomAlgorithm(f func(e []*Edge) []int) {
	bag.customAlgorithm = f
}

func (bag *Bag) SetCustomWeighing(f func(src, dst []int) float64) {
	bag.customWeighing = f
}

func (bag *Bag) SetWordTokenizer(f func(string) []string) {
	bag.wordTokenizer = f
}

// Summarize the text to num sentences
func (bag *Bag) Summarize(text string, num int) (string, error) {
	text = strings.TrimSpace(text)
	if len(text) < 1 && len(bag.OriginalSentences) == 0 {
		return "", nil
	}

	bag.createSentences(text) // only actually creates sentences if no OrignalSentences

	// If user already provide dictionary, pass creating dictionary
	if len(bag.Dict) < 1 {
		if text == "" {
			text = strings.TrimSpace(strings.Join(bag.OriginalSentences, " "))
		}
		bag.createDictionary(text)
	}

	bag.createNodes()
	bag.createEdges()

	switch bag.Algorithm {
	case "centrality":
		bag.centrality()
		break
	case "pagerank":
		bag.pageRank()
		break
	case "custom":
		bag.Ranks = bag.customAlgorithm(bag.Edges)
		break
	default:
		bag.pageRank()
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
		max := len(r)
		if max > bag.MaxCharacters {
			max = bag.MaxCharacters
		}
		r = r[:max]
		// then turn back to string
		res = string(r)
	}

	return res, nil
}

type Rank struct {
	idx   int
	score float64
}

func (bag *Bag) centrality() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
		if edge.weight > bag.Threshold {
			newEdges = append(newEdges, edge)
		}
	}
	// sort them by weight descending
	sort.Sort(ByWeight(newEdges))
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

func (bag *Bag) pageRank() {
	// first remove edges under Threshold weight
	var newEdges []*Edge
	for _, edge := range bag.Edges {
		if edge.weight > bag.Threshold {
			newEdges = append(newEdges, edge)
		}
	}
	// then page rank them
	graph := pagerank.NewGraph()
	defer graph.Reset()
	for _, edge := range newEdges {
		graph.Link(uint32(edge.src), uint32(edge.dst), edge.weight)
	}
	var ranks []*Rank
	graph.Rank(bag.Damping, bag.Tolerance, func(sentenceIndex uint32, rank float64) {
		ranks = append(ranks, &Rank{int(sentenceIndex), rank})
	})
	// sort ranks into an array of sentence index, by score descending
	sort.Sort(ByScore(ranks))
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

func (bag *Bag) createEdges() {
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
				case "custom":
					weight = bag.customWeighing(src.vector, dst.vector)
					break
				default:
					differentElements := SymmetricDifference(src.vector, dst.vector)
					weight = float64(len(differentElements))
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

func (bag *Bag) createNodes() {
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

func (bag *Bag) createSentences(text string) {
	if len(bag.OriginalSentences) == 0 {
		// trim all spaces
		// done by calling func: text = strings.TrimSpace(text)
		// tokenize text as sentences
		// sentence is a group of words separated by whitespaces or punctuation other than !?.
		bag.OriginalSentences = TokenizeSentences(text)
	}

	// from original sentences, explode each sentences into bag of words
	bag.BagOfWordsPerSentence = [][]string{}
	for _, sentence := range bag.OriginalSentences {
		words := bag.wordTokenizer(sentence)
		bag.BagOfWordsPerSentence = append(bag.BagOfWordsPerSentence, words)
	}

	// then uniq it
	UniqSentences(bag.BagOfWordsPerSentence, bag.SentencesDistanceThreshold)
}

func (bag *Bag) createDictionary(text string) {
	// trim all spaces
	// this already done by calling func:	text = strings.TrimSpace(text)
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
