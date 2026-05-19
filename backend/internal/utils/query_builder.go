package utils

import (
	"fmt"
	"strings"
)

// Placeholders generates PostgreSQL placeholders like $1, $2, $3...
// count: number of placeholders to generate
// start: starting number (usually 1)
func Placeholders(count, start int) string {
	if count <= 0 {
		return ""
	}
	
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = fmt.Sprintf("$%d", start+i)
	}
	return strings.Join(result, ", ")
}

// JoinPlaceholders joins multiple placeholder sets with commas
// e.g., for 3 fields starting at 1: "$1, $2, $3"
func JoinPlaceholders(fields []string, start int) string {
	if len(fields) == 0 {
		return ""
	}
	return Placeholders(len(fields), start)
}