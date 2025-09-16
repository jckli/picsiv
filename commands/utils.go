package commands

import (
	"strings"
)

func fuzzySearch(arr []string, searchStr string) []string {
	var result []string
	for _, str := range arr {
		if strings.Contains(str, searchStr) {
			result = append(result, str)
		}
	}
	return result
}
