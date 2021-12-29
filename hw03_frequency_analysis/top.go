package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var regClean = regexp.MustCompile(`[\p{L}-]+`)

func Top10(str string) []string {
	var topWords []string
	if len(str) == 0 {
		return topWords
	}

	counter := make(map[string]int8)
	words := strings.Fields(str)
	for _, word := range words {
		word = normalizeWord(word)
		if len(word) == 0 {
			continue
		}

		count, ok := counter[word]
		if !ok {
			topWords = append(topWords, word)
			count = 0
		}
		counter[word] = count + 1
	}

	sort.Slice(topWords, func(i, j int) bool {
		if counter[topWords[i]] == counter[topWords[j]] {
			return strings.Compare(topWords[i], topWords[j]) < 0
		}
		return counter[topWords[i]] > counter[topWords[j]]
	})

	if len(topWords) > 10 {
		topWords = topWords[:10]
	}

	return topWords
}

func normalizeWord(word string) string {
	word = strings.Trim(word, ` `)
	if !taskWithAsteriskIsCompleted || len(word) == 0 {
		return word
	}

	word = strings.ToLower(word)
	word = regClean.FindString(word)
	word = strings.Trim(word, `-`)

	return word
}
