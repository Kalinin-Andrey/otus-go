package hw03

import (
	"sort"
	"strings"
	"unicode"
)

// Returns sorted top 10 words with max frequency for a given text
func frequencyTextAnalysis(text string) (res []string) {
	frequenciesByWords := make(map[string]int, 20)
	wordsByFrequencies := make(map[int][]string, 20)
	frequencies := make([]int, 0, 20)

	text = cleanAndLowText(text)
	words := strings.Fields(text)

	for _, word := range words {
		if _, ok := frequenciesByWords[word]; !ok {
			frequenciesByWords[word] = 1
		} else {
			frequenciesByWords[word]++
		}
	}

	for word, count := range frequenciesByWords {
		wordsByFrequencies[count] = append(wordsByFrequencies[count], word)
	}

	for count := range wordsByFrequencies {
		frequencies = append(frequencies, count)
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i] > frequencies[j]
	})

	for _, count := range frequencies {
		res = append(res, wordsByFrequencies[count]...)

		if len(res) >= 10 {
			break
		}
	}

	if len(res) > 10 {
		res = res[:10]
	}
	sort.Strings(res)

	return res
}

// Cleans and turns given text to lower case
func cleanAndLowText(sourceText string) string {
	resText := strings.Builder{}

	for _, r := range sourceText {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			resText.WriteRune(r)
		}
	}

	return strings.ToLower(resText.String())
}