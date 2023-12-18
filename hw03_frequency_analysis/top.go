package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	words := strings.Fields(input)

	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	type wordCount struct {
		word  string
		count int
	}

	//nolint:prealloc
	var wordCounts []wordCount
	for word, count := range wordFreq {
		wordCounts = append(wordCounts, wordCount{word, count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	top := 10
	if len(wordCounts) < top {
		top = len(wordCounts)
	}

	result := make([]string, top)
	for i := 0; i < top; i++ {
		result[i] = wordCounts[i].word
	}

	return result
}
