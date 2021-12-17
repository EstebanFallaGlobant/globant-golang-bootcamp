package testapi

import (
	"testing"

	wc "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/lib"
	wcStructs "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/structs"
	"github.com/stretchr/testify/assert"
)

func Test_CountWords_ValidString_NoRepeats(t *testing.T) {

	testPhrase := "This is a test"
	expectedResult := []wcStructs.WordCount{
		{
			Word:  "this",
			Count: 1,
		},
		{
			Word:  "is",
			Count: 1,
		},
		{
			Word:  "a",
			Count: 1,
		},
		{
			Word:  "test",
			Count: 1,
		},
	}

	response := wc.CountWords(testPhrase)

	assert.NotEmpty(t, response)
	assert.ElementsMatch(t, expectedResult, response)
}

func Test_CountWords_SingleRepeated(t *testing.T) {
	testPhrase := "This is a test with \"this\" repeated"
	expectedResult := 2

	response := wc.CountWords(testPhrase)
	var elementCount int

	for i := range response {
		if element := response[i]; element.Word == "this" {
			elementCount = element.Count
		}
	}

	assert.NotEmpty(t, response)
	assert.EqualValues(t, expectedResult, elementCount)
}
