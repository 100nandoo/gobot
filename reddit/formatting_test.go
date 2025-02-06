package reddit

import (
	"testing"
)

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Text with underscore",
			input:    "Hello_World",
			expected: "Hello\\_World",
		},
		{
			name:     "Text with asterisk",
			input:    "Hello*World",
			expected: "Hello\\*World",
		},
		{
			name:     "Text with backtick",
			input:    "Hello`World",
			expected: "Hello\\`World",
		},
		{
			name:     "Text with multiple special characters",
			input:    "Hello_World*with `backtick",
			expected: "Hello\\_World\\*with \\`backtick",
		},
		{
			name:     "Summer in Paris [OC @tombaenre]",
			input:    "Summer in Paris [OC @tombaenre]",
			expected: "Summer in Paris [OC @tombaenre]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := escapeMarkdown(test.input)
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}