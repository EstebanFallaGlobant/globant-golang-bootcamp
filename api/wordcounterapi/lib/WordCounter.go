package wordcounterapi

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/maps-lib/wordcounter"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/structs"
)

type WordCounter struct{}

func (WordCounter *WordCounter) CountWords(text string) []structs.WordCount {
	var result []structs.WordCount

	for word, count := range wordcounter.WordCountCaseInsensitive(text) {
		result = append(result, structs.WordCount{Word: word, Count: count})
	}

	return result
}
