package main

import (
	"fmt"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/maps-lib/wordcounter"
)

func main() {
	s := "Test of word counter, this is a test"
	wordMap := wordcounter.WordCountCaseInsensitive(s)

	fmt.Print(wordMap)
	fmt.Println()
}
