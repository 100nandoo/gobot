package reddit

import (
	"html"
	"strings"
)

func escapeMarkdown(text string) string {
	// First decode HTML entities
	decoded := html.UnescapeString(text)

	specialChars := []string{"_", "*", "`"}
	escaped := decoded
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	println(escaped)
	return escaped
}
