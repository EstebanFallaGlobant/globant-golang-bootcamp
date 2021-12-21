package interfaces

import "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/WordCounter/wordcounterapi/structs"

type WordCounterInterface interface {
	CountWords(text string) []structs.WordCount
}
