package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/maps-lib/wordcounter"
	"github.com/stretchr/testify/assert"
)

func Test_WordCountCaseInsensitive_ThreeWordsNoRepeat(t *testing.T) {
	testPhrase := "Just three words."
	testResult := map[string]int{
		"just":  1,
		"words": 1,
		"three": 1}

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	assert.EqualValues(t, testResult, wordCount)
}
func Test_WordCountCaseInsensitive_EmptyString(t *testing.T) {
	testPhrase := ""

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	assert.Empty(t, wordCount)
}

func Test_WordCountCaseInsensitive_WordRepeatSameCase(t *testing.T) {
	testPhrase, resultKey, resultValue := "The word \"word\" appears twice", "word", 2

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	assert.EqualValues(t, resultValue, wordCount[resultKey])
}

func Test_WordCountCaseInsensitive_WordRepeatDifferentCase(t *testing.T) {
	testPhrase, resultKey, resultValue := "word is The WORD. WoRd appears three times", "word", 3

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	assert.EqualValues(t, resultValue, wordCount[resultKey])
}

func Test_WordCount_ThreeWordsNoRepeat(t *testing.T) {
	testPhrase, expectedValue := "Just three words", 1

	wordCount := wordcounter.WordCount(testPhrase)

	for _, val := range wordCount {
		assert.EqualValues(t, expectedValue, val)
	}
}

func Test_WordCount_SingleWordRepeatedSameCase(t *testing.T) {
	testPhrase, expectedKey, expectedValue := "The word repeated is repeated twice", "repeated", 2

	wordCount := wordcounter.WordCount(testPhrase)

	assert.EqualValues(t, expectedValue, wordCount[expectedKey])
}

func Test_WordCount_SingleWordRepeatedDifferentCase(t *testing.T) {
	testPhrase, expectedKey1, expectedKey2, expectedValue := "The word repeated is REPEATED twice", "repeated", "REPEATED", 1

	wordCount := wordcounter.WordCount(testPhrase)

	assert.EqualValues(t, expectedValue, wordCount[expectedKey1])
	assert.EqualValues(t, expectedValue, wordCount[expectedKey2])
}

func Test_WordCount_EmptyString(t *testing.T) {
	var testPhrase string

	wordCount := wordcounter.WordCount(testPhrase)

	assert.Empty(t, wordCount)
}
