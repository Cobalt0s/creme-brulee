package gutils

import (
	"log"
	"sort"
	"strings"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func UniqueSetToList(uniqueSet map[string]bool) []string {
	result := make([]string, len(uniqueSet))
	i := 0
	for k, _ := range uniqueSet {
		result[i] = k
		i++
	}
	return result
}

func SortStrings(input []string) {
	sort.Slice(input, func(i, j int) bool {
		return input[i] < input[j]
	})
}

func NoSpaces(text string) string {
	for strings.Contains(text, " ") {
		text = strings.Replace(text, " ", "", -1)
	}
	return text
}
