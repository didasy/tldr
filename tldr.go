/*
Dependencies:
  go get github.com/dcadenas/pagerank

BUG : 
1. if there is no space before \n, throw index out of range error from createNodes function. Somehow a word doesn't register on the dict and it cause the error because if not found in dict it returns 0.

FIX :
1. Added some more parameters at createDictionary in the if decision in the strings.Map.

TODO :
1. Try Hamming distance instead of Jaccard Coeficient for calculating node weights - Done
2. Try using idf-modified-cosine
*/
package tldr

import (
	"sort"
	"strings"
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

func New() *Bag {
	return &Bag{}
}

// the default values of each settings
const (
	VERSION = "0.0.2"
	ALGORITHM = "centrality"
	WEIGHING = "jaccard"
	DAMPING = 0.85
	TOLERANCE = 0.0001
	THRESHOLD = 0.001
)

// Using pagerank algorithm will return many version of summary, unlike static summary result from centrality algorithm
var (
	Algorithm string = "centrality"
	Weighing string = "jaccard"
	Damping float64 = 0.85
	Tolerance float64 = 0.0001
	Threshold float64 = 0.001
)

func Set(d float64, t float64, th float64, alg string, w string) {
	Damping = d
	Tolerance = t
	Threshold = th
	Algorithm = alg
	Weighing = w
}

func (bag *Bag) Summarize(text string, num int) string {
	bag.createDictionary(text)
	bag.createSentences(text)
	bag.createNodes()
	bag.createEdges()
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
	for _, v := range idx {
		res += bag.originalSentences[v] + " "
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
	var rankBySrc []int
	for _, v := range newEdges {
		rankBySrc = append(rankBySrc, v.src)
	}
	// uniq it without disturbing the order
	// var uniq []int
	// uniq = append(uniq, rankBySrc[0])
	// for _, v := range rankBySrc {
	// 	same := false
	// 	for j := 0; j < len(uniq); j++ {
	// 		if uniq[j] == v {
	// 			same = true
	// 		}
	// 	}
	// 	if !same {
	// 		uniq = append(uniq, v)
	// 	}
	// }
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
	ranks := make(map[int]float64)
	graph.Rank(Damping, Tolerance, func (sentenceIndex int, rank float64) {
		ranks[sentenceIndex] = rank
	})
	// sort ranks into an array of sentence index, by rank descending
	var idx []int
	for i, v := range ranks {
		highest := i
		for j, x := range ranks {
			if i != j && x > v {
				highest = j
			}
		}
		idx = append(idx, highest)
		delete(ranks, highest)
		if len(ranks) == 2 {
			for l, z := range ranks {
				for m, r := range ranks {
					if r >= z {
						idx = append(idx, m)
						idx = append(idx, l)
						delete(ranks, m)
					}
				}
			}
		}
	}
	bag.ranks = idx
}

type Edge struct {
	src int // index of node
	dst int // index of node
	weight float64 // weight of the similarity between two sentences, use Jaccard Coefficient
}

func (bag *Bag) createEdges() {
	for i, src := range bag.nodes {
		for j, dst := range bag.nodes {
			// don't compare same node
			if i != j {
				var weight float64
				if Weighing == "jaccard" {
					commonElements := intersection(src.vector, dst.vector)
					weight = float64(len(commonElements)) / ((float64(vectorLength) * 2) - float64(len(commonElements)))
				} else if Weighing == "hamming" {
					differentElements := symetricDifference(src.vector, dst.vector)
					weight = float64(len(differentElements))
				} else {
					commonElements := intersection(src.vector, dst.vector)
					weight = float64(len(commonElements)) / ((float64(vectorLength) * 2) - float64(len(commonElements)))
				}
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}


func symetricDifference(src []int, dst []int) []int {
	var diff []int
	for i, v := range src {
		if v != dst[i] {
			diff = append(diff, i)
		} 
	}
	return diff
}

func intersection(src []int, dst []int) []int {
	intersect := make(map[int]bool)
	for i, v := range src {
		if v > 0 && dst[i] > 0 {
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
	vector = [1, 1, 0, 2]
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
			// check word dict position
			// minus 1, because array started from 0 and lowest dict is 1
			pos := bag.dict[word] - 1
			// increment the position
			vector[pos]++
		}
		// vector is now created, put it into the node
		node := &Node{i, vector}
		// node is now completed, put into the bag
		bag.nodes = append(bag.nodes, node)
	}
}

func (bag *Bag) createSentences(text string) {
	// trim all spaces
	text = strings.TrimSpace(text)
	words := strings.Fields(text)
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
	// remove doubled sentence
	sentences = uniqSentences(sentences)
	// now flatten them
	var bagOfSentence []string
	for _, s := range sentences {
		str := strings.Join(s, " ")
		bagOfSentence = append(bagOfSentence, str)
	}
	bag.originalSentences = bagOfSentence
	// sanitize sentences before putting it into the bag
	bag.sentences = sanitizeSentences(sentences)
}

func uniqSentences(sentences [][]string) [][]string {
	var z []string
	// create a sentence as one string and append it to z
	for _, v := range sentences {
		j := strings.Join(v ," ")
		z = append(z, j)
	}
	// var uniq []string
	// uniq = append(uniq, z[0])
	// for _, v := range z {
	// 	same := false
	// 	for j := 0; j < len(uniq); j++ {
	// 		if uniq[j] == v {
	// 			same = true
	// 		}
	// 	}
	// 	if !same {
	// 		uniq = append(uniq, v)
	// 	}
	// }
	m := make(map[string]bool)
	var uniq []string
	for _, v := range z {
		if m[v] {
			continue
		}
		uniq = append(uniq, v)
		m[v] = true
	}
	var unique [][]string
	for _, v := range uniq {
		unique = append(unique, strings.Fields(v))
	}
	return unique
}

func sanitizeSentences(sentences [][]string) [][]string {
	var sanitizedSentence [][]string
	for _, sentence := range sentences {
		var newSentence []string
		for _, word := range sentence {
			word = strings.ToLower(word)
			word = strings.Map(func (r rune) rune {
				if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != ' ' {
					return -1
				}
				return r
			}, word)
			newSentence = append(newSentence, word)
		}
		sanitizedSentence = append(sanitizedSentence, newSentence)
	}
	return sanitizedSentence
}

func (bag *Bag) createDictionary(text string) {
	// trim all spaces
	text = strings.TrimSpace(text)
	// lowercase the text
	text = strings.ToLower(text)
	// remove all non alphanumerics
	text = strings.Map(func (r rune) rune {
		// probably would be cleaner if use !unicode.IsDigit, !unicode.IsLetter, and !unicode.IsSpace
		if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != ' ' && r != '\n' && r != '\t' && r != '\v' && r != '\f' && r!= '\r' {
			return -1
		}
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