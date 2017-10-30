package tldr

import (
	"math"
	"regexp"
	"strings"
)

var sanitize, sentenceTokenizer *regexp.Regexp

func init() {
	sanitize = regexp.MustCompile(`([^\p{L}\d]{2,}|[^\p{L}\d_'-])`)
	sentenceTokenizer = regexp.MustCompile(`([\.\?\!])(?:\s|$)`)
}

func TokenizeSentences(text string) []string {
	tokens := []string{}

	text = strings.TrimSpace(text)

	// [][]int
	idxMap := sentenceTokenizer.FindAllStringIndex(text, -1)

	// cut by guide
	from := 0
	for _, c := range idxMap {
		str := text[from : c[0]+1]
		str = strings.TrimSpace(str)
		tokens = append(tokens, str)
		from = c[1]
	}

	return tokens
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

func UniqSentences(sentences [][]string, sentenceDistanceThreshold float64) {
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
			if Distance(sen, msens[j]) >= sentenceDistanceThreshold {
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
	word = sanitize.ReplaceAllString(word, "")

	return word
}
