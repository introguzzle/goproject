package utils

import (
	_ "golang.org/x/text/cases"
	"strings"
)

func ToSnakeCase(s string) string {
	var result []string
	for i, c := range s {
		if i > 0 && (c >= 'A' && c <= 'Z') {
			result = append(result, "_")
		}
		result = append(result, string(c))
	}
	return strings.ToLower(strings.Join(result, ""))
}

func ToCamelCase(snake string) string {
	parts := strings.Split(snake, "_")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	return strings.Join(parts, "")
}
