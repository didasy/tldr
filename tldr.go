/*
Dependencies:
  go get github.com/dcadenas/pagerank

BUG : 
1. if there is no space before \n, throw index out of range error from createNodes function. Somehow a word doesn't register on the dict and it cause the error because if not found in dict it returns 0.

FIX :
1. Added some more parameters at createDictionary in the if decision in the strings.Map, also add a guard in createNodes so if there is unknown word it would not crash.

TODO :
1. Try using idf-modified-cosine - done
*/
package tldr

import (
	"sort"
	"math"
	"bytes"
	"strings"
	// "unicode"
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
	VERSION = "0.1.0"
	ALGORITHM = "centrality"
	WEIGHING = "hamming"
	DAMPING = 0.85
	TOLERANCE = 0.0001
	THRESHOLD = 0.001
)

// Using pagerank algorithm will return many version of summary, unlike static summary result from centrality algorithm
var (
	Algorithm string = "centrality"
	Weighing string = "hamming"
	Damping float64 = 0.85
	Tolerance float64 = 0.0001
	Threshold float64 = 0.001
)

func Set(d, t, th float64, alg, w string) {
	Damping = d
	Tolerance = t
	Threshold = th
	Algorithm = alg
	Weighing = w
}

func (bag *Bag) Summarize(text string, num int) string {
	bag.createOriginalSentences(text)
	bag.createSentences(text)
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
	// ranks := make(map[int]float64)
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
	var idx []int
	for _, v := range ranks {
		idx = append(idx, v.idx)
	}
	
	bag.ranks = idx
}

type Edge struct {
	src int // index of node
	dst int // index of node
	weight float64 // weight of the similarity between two sentences, use Jaccard Coefficient
}

// this is also experimental, using triple bytes with ferret algorithm
func (bag *Bag) createByteFerretEdges() {
	for i, src := range bag.originalSentences {
		for j, dst := range bag.originalSentences {
			if i != j {
				// src = strings.TrimSpace(src)
				// dst = strings.TrimSpace(dst)
				src = strings.ToLower(src)
				dst = strings.ToLower(dst)
				// compact the strings
				src = strings.Map(func (r rune) rune {
					if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
						return -1
					}
					return r
				}, src)
				dst = strings.Map(func (r rune) rune {
					if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
						return -1
					}
					return r
				}, dst)
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
				// this is Jaccard's similarity
				weight = weight / ( ( float64( len(srcB) ) + float64( len(dstB) ) ) - weight )
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}
	}
}

// this is experimental, using raw sentences string instead of sanitized and tokenized one, try sanitize this
func (bag *Bag) createJaroWinklerEdges() {
	for i, src := range bag.originalSentences {
		for j, dst := range bag.originalSentences {
			if i != j {
				// src = strings.TrimSpace(src)
				// dst = strings.TrimSpace(dst)
				src = strings.ToLower(src)
				dst = strings.ToLower(dst)
				// compact the strings
				src = strings.Map(func (r rune) rune {
					if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
						return -1
					}
					return r
				}, src)
				dst = strings.Map(func (r rune) rune {
					if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
						return -1
					}
					return r
				}, dst)
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
				srcTfVector := createTfVector(srcS, seqDict)
				dstTfVector := createTfVector(dstS, seqDict)
				idf := createIdf(srcS, dstS, seqDict)
				srcTfIdfVector := createTfIdfVector(srcTfVector, idf)
				dstTfIdfVector := createTfIdfVector(dstTfVector, idf)
				// https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/ for more explanation
				// find the dot-product-idf-modified of them
				dotProduct := createDotProduct(srcTfIdfVector, dstTfIdfVector)
				// find the sum of magnitude of each tfidf in each sentences
				srcM := createMagnitude(srcTfIdfVector)
				dstM := createMagnitude(dstTfIdfVector)
				// now calculate tf-idf-modified-cosine-similarity between the sentences
				// http://en.wikipedia.org/wiki/Cosine_similarity exactly like this, but switch tf with tfidf
				// http://upload.wikimedia.org/math/f/3/6/f369863aa2814d6e283f859986a1574d.png for the formula
				weight = dotProduct / (srcM * dstM)
				// put them into the bag
				edge := &Edge{i, j, weight}
				bag.edges = append(bag.edges, edge)
			}
		}		
	}
}

func createTfVector(sen, seqDict []string) []float64 {
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
	return vector
}

func createIdf(src, dst, seqDict []string) []float64 {
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
	return idf
}

func createTfIdfVector(vector, idf []float64) []float64 {
	tfidf := make([]float64, len(vector))
	for i, tf := range vector {
		tfidf[i] = tf * idf[i]
	}
	return tfidf
}

func createDotProduct(srcVector, dstVector []float64) float64 {
	var dotProduct float64
	for i, v := range srcVector {
		dotProduct += (v * dstVector[i])
	}
	return dotProduct
}

func createMagnitude(vector []float64) float64 {
	var magnitude float64
	for _, tfidf := range vector {
		magnitude += (tfidf * tfidf)
	}
	magnitude = math.Sqrt(magnitude)
	return magnitude
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

func (bag *Bag) createSentences(text string) {
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
			word = sanitizeWord(word)
			sentence = append(sentence, word)
			sentences = append(sentences, sentence)
			sentence = []string{}
		} else {
			word = sanitizeWord(word)
			sentence = append(sentence, word)
		}
	}
	if len(sentence) > 0 {
		sentences = append(sentences, sentence)
	}
	// remove doubled sentence
	sentences = uniqSentences(sentences)
	// sanitize sentences before putting it into the bag
	bag.sentences = sanitizeSentences(sentences)
}

func (bag *Bag) createOriginalSentences(text string) {
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

func sanitizeWord(word string) string {
	word = strings.ToLower(word)
	var prev rune
	word = strings.Map(func (r rune) rune {
		// don't remove '-' if it exists after alphanumerics
		if r == '-' && ((prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || prev == 'ä' || prev == 'ö' || prev == 'ü' || prev == 'ß' || prev == 'é') {
			return r
		}
		if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
			return -1
		}
		prev = r
		return r
	}, word)
	return word
}

func sanitizeSentences(sentences [][]string) [][]string {
	var sanitizedSentence [][]string
	for _, sentence := range sentences {
		var newSentence []string
		for _, word := range sentence {
			var prev rune
			word = strings.ToLower(word)
			word = strings.Map(func (r rune) rune {
				// experimental, we don't remove '-' if it exists after alphanumerics
				if r == '-' && ((prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || prev == 'ä' || prev == 'ö' || prev == 'ü' || prev == 'ß' || prev == 'é') {
					return r
				}
				if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' {
					return -1
				}
				prev = r
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
	// remove all non alphanumerics but spaces
	var prev rune
	text = strings.Map(func (r rune) rune {
		// probably would be cleaner if use !unicode.IsDigit, !unicode.IsLetter, and !unicode.IsSpace
		// but could also be slower
		// This one is experimental, we don't remove '-' if it exists after alphanumerics
		if r == '-' && ((prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || prev == 'ä' || prev == 'ö' || prev == 'ü' || prev == 'ß' || prev == 'é') {
			return r
		}
		if (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != 'ä' && r != 'ö' && r != 'ü' && r != 'ß' && r != 'é' && r != ' ' && r != '\n' && r != '\t' && r != '\v' && r != '\f' && r!= '\r' {
			return -1
		}
		prev = r
		// if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsSpace(r) {
		// 	return -1
		// }
		return r
	}, text)
	// TRYING TO FIX BUG : remove all double spaces left
	text = strings.Replace(text, "  ", " ", -1)
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