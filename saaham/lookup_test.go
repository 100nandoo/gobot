package saaham

import "testing"

func TestNormalizeTickerToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		ok    bool
	}{
		{name: "plain symbol", input: "aapl", want: "AAPL", ok: true},
		{name: "wrapper punctuation", input: "(vwra.l)", want: "VWRA.L", ok: true},
		{name: "trailing punctuation", input: "msft!", want: "MSFT", ok: true},
		{name: "mixed prose", input: "show me aapl", ok: false},
		{name: "dollar prefix", input: "$aapl", ok: false},
		{name: "multiple symbols", input: "AAPL MSFT", ok: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := normalizeTickerToken(tc.input)
			if ok != tc.ok {
				t.Fatalf("expected ok=%v, got %v", tc.ok, ok)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestParseCommandLookup(t *testing.T) {
	got, err := parseCommandLookup([]string{"(aapl)"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "AAPL" {
		t.Fatalf("expected AAPL, got %q", got)
	}

	if _, err := parseCommandLookup([]string{"show", "me", "aapl"}); err == nil {
		t.Fatal("expected invalid symbol error for mixed prose command")
	}
}
