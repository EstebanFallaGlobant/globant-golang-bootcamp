package wordcounter

import "strings"

func WordCount(s string) map[string]int {
	result := make(map[string]int)
	var word string
	fields := strings.Fields(s)

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
