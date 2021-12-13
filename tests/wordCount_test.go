package tests

import (
	"reflect"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/maps-lib/wordcounter"
)

func Test_WordCountCaseInsensitive_ThreeWordsNoRepeat(t *testing.T) {
	testPhrase := "Just three words."
	testResult := map[string]int{
		"just":  1,
		"words": 1,
		"three": 1}

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	if v := reflect.DeepEqual(wordCount, testResult); !v {
		t.Fatalf("The map %v isn't equeal to the expected map: %v\n", wordCount, testResult)
	}
}
func Test_WordCountCaseInsensitive_EmptyString(t *testing.T) {
	testPhrase := ""
	testResult := make(map[string]int)

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	if len(wordCount) != 0 {
		t.Fatalf("The map %v is not empty, expected: %v\n", wordCount, testResult)
	}
}

func Test_WordCountCaseInsensitive_WordRepeatSameCase(t *testing.T) {
	testPhrase, resultKey, resultValue := "The word \"word\" appears twice", "word", 2

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	if v := wordCount[resultKey]; v != resultValue {
		t.Fatalf("The value of the key \"%s\" is different from the expected. Expected: %d, Got: %d\n", resultKey, resultValue, v)
	}
}

func Test_WordCountCaseInsensitive_WordRepeatDifferentCase(t *testing.T) {
	testPhrase, resultKey, resultValue := "word is The WORD. WoRd appears three times", "word", 3

	wordCount := wordcounter.WordCountCaseInsensitive(testPhrase)

	if v := wordCount[resultKey]; v != resultValue {
		t.Fatalf("The value of the key \"%s\" is different from the expected. Expected: %d, Got: %d\n", resultKey, resultValue, v)
	}
}

func Test_WordCount_ThreeWordsNoRepeat(t *testing.T) {
	testPhrase, expectedValue := "Just three words", 1

	wordCount := wordcounter.WordCount(testPhrase)

	t.Log(wordCount)
	for key, val := range wordCount {
		if val != expectedValue {
			t.Fatalf("The key \"%s\" has an incorrect value. Expected: %d. Got: %d\n", key, expectedValue, val)
		}
	}
}

func Test_WordCount_SingleWordRepeatedSameCase(t *testing.T) {
	testPhrase, expectedKey, expectedValue := "The word repeated is repeated twice", "repeated", 2

	wordCount := wordcounter.WordCount(testPhrase)

	if v := wordCount[expectedKey]; v > expectedValue {
		t.Fatalf("The key \"%s\" has an incorrect value. Expected: %d. Got: %d\n", expectedKey, expectedValue, v)
	}
}

func Test_WordCount_SingleWordRepeatedDifferentCase(t *testing.T) {
	testPhrase, expectedKey1, expectedKey2, expectedValue := "The word repeated is REPEATED twice", "repeated", "REPEATED", 1

	wordCount := wordcounter.WordCount(testPhrase)

	if v1, v2 := wordCount[expectedKey1], wordCount[expectedKey2]; v1 != expectedValue || v2 != expectedValue {
		t.Fatal()
	}
}

func Test_WordCount_EmptyString(t *testing.T) {
	var testPhrase string

	wordCount := wordcounter.WordCount(testPhrase)

	if val, expr := len(wordCount), 0; val != expr {
		t.Fatalf("The word count is incorrect. Expected: %d. Got: %d}n", expr, val)
	}
}
