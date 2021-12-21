package interfaces

import "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/structs"

type WordCounterInterface interface {
	CountWords(text string) []structs.WordCount
}
