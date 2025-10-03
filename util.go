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
	// Pre-allocate slice with reasonable capacity to avoid multiple allocations
	result := make([]int, 0, len(src)/4) // Estimate 25% intersection

	for i, v := range src {
		if v == dst[i] {
			result = append(result, i)
		}
	}
	return result
}

func UniqSentences(sentences [][]string, sentenceDistanceThreshold float64) {
	// Pre-allocate msens with exact capacity
	msens := make([]string, 0, len(sentences))
	for _, sen := range sentences {
		msens = append(msens, strings.Join(sen, " "))
	}

	// Do jarowinkler then CSIS to deduplicate sentences
	reject := make(map[int]bool, len(msens))

	// First JaroWinkler - optimized to avoid redundant comparisons
	for i := 0; i < len(msens)-1; i++ {
		if reject[i] {
			continue // Skip if already rejected
		}
		sen := msens[i]
		for j := i + 1; j < len(msens); j++ {
			if !reject[j] && Distance(sen, msens[j]) >= sentenceDistanceThreshold {
				reject[j] = true
			}
		}
	}

	// Then CSIS - optimized to avoid redundant comparisons
	for i := 0; i < len(msens)-1; i++ {
		if reject[i] {
			continue
		}
		psen := msens[i]
		for j := i + 1; j < len(msens); j++ {
			if i != j && !reject[j] {
				nsen := msens[j]
				// if i subset of j, put i in reject
				if strings.Contains(nsen, psen) {
					reject[i] = true
					break // i is rejected, no need to check more j's
				}
				// if j subset of i, put j in reject
				if strings.Contains(psen, nsen) {
					reject[j] = true
				}
			}
		}
	}

	// Rebuild sentences slice in place
	keepCount := len(msens) - len(reject)
	result := make([][]string, 0, keepCount)
	for i, sen := range msens {
		if !reject[i] {
			result = append(result, strings.Fields(sen))
		}
	}

	// Copy back to original slice to maintain same reference
	if len(result) < len(sentences) {
		// Clear and rebuild the original slice
		sentences = sentences[:0]
		for _, sen := range result {
			sentences = append(sentences, sen)
		}
	}
}

func SanitizeWord(word string) string {
	word = strings.ToLower(word)
	word = sanitize.ReplaceAllString(word, "")

	return word
}
