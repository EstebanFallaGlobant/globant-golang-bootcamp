package wordcounter

import (
	"fmt"
	"regexp"
	"strings"
)

func WordCountCaseInsensitive(s string) map[string]int {
	result := make(map[string]int)
	var word string
	fields := splitByRegex(s)

	for i := 0; i < len(fields); i++ {
		word = strings.ToLower(fields[i])
		if value, ok := result[word]; ok {
			result[word] = value + 1
		} else {
			result[word] = 1
		}
	}
	return result
}

func WordCount(s string) map[string]int {
	result := make(map[string]int)
	var word string
	fields := splitByRegex(s)

	fmt.Printf("%v\n", fields)

	for i := 0; i < len(fields); i++ {
		word = fields[i]
		if value, ok := result[word]; ok {
			result[word] = value + 1
		} else {
			result[word] = 1
		}
	}
	return result
}

func splitByRegex(text string) []string {
	re := regexp.MustCompile(`\w+`)
	result := re.FindAllString(text, -1)
	return result
}
