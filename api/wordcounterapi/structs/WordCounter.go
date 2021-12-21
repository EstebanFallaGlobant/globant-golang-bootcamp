package structs

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/maps-lib/wordcounter"
)

type WordCounter struct{}

func (WordCounter *WordCounter) CountWords(text string) []WordCount {
	var result []WordCount

	for word, count := range wordcounter.WordCountCaseInsensitive(text) {
		result = append(result, WordCount{Word: word, Count: count})
	}

	return result
}
