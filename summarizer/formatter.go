package summarizer

import (
	"fmt"
	"strings"
)

func (s SmmryResponse) ToMarkdownString() string {
	var builder strings.Builder

	if len(s.SmAPITitle) > 0 {
		builder.WriteString("*")
		builder.WriteString(s.SmAPITitle)
		builder.WriteString("*\n\n")
	}

	sentences := strings.Split(s.SmAPIContent, ". ")
	for _, sentence := range sentences {
		builder.WriteString(sentence)
		builder.WriteString(".\n\n")
	}

	builder.WriteString(fmt.Sprintf("\nContent reduced by *%s*", s.SmAPIContentReduced))

	return builder.String()
}
