package hw03

import (
	"sort"
	"strings"
	"unicode"
)

const (
	StartCap            = 20
	DefaultResultLength = 10
)

type WordsCount struct {
	word  string
	count int
}

// Returns sorted top 10 words with max frequency for a given text
func frequencyTextAnalysis(text string) (res []string) {
	frequenciesByWords := make(map[string]int, StartCap)
	wc := make([]WordsCount, 0, StartCap)

	text = cleanAndLowText(text)
	words := strings.Fields(text)

	for _, word := range words {

		if frequenciesByWords[word] == 0 {
			wc = append(wc, WordsCount{
				word: word,
			})
		}
		frequenciesByWords[word]++
	}

	for _, wordsCount := range wc {
		wordsCount.count = frequenciesByWords[wordsCount.word]
	}

	sort.SliceStable(wc, func(i, j int) bool {
		return wc[i].count < wc[j].count
	})

	resultLen := DefaultResultLength
	if resultLen > len(wc) {
		resultLen = len(wc)
	}

	for _, wordsCount := range wc[:resultLen] {
		res = append(res, wordsCount.word)
	}

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
