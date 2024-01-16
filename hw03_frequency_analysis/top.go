package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reSpaceOrNewline = regexp.MustCompile(`\s+|\n`)

type wordCount struct {
	word  string
	count int
}

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}
	words := reSpaceOrNewline.Split(input, -1)

	wordFreq := make(map[string]int)
	for _, word := range words {
		word = strings.Trim(word, ".,!?:;'\"")
		word = strings.ToLower(word)
		if word == "-" || word == "—" {
			continue
		}
		wordFreq[word]++
	}

	wordCounts := make([]wordCount, 0, len(wordFreq))
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
