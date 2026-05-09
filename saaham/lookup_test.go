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
		{name: "caret prefixed index", input: "^sti", want: "^STI", ok: true},
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

func TestResolveTickerCandidates(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		want   []string
	}{
		{name: "s&p 100 member keeps bare symbol first", symbol: "AAPL", want: []string{"AAPL", "^AAPL", "AAPL.L", "AAPL.JK"}},
		{name: "non-member uses shortcut-first order with bare fallback", symbol: "CSPX", want: []string{"^CSPX", "CSPX.L", "CSPX.JK", "CSPX"}},
		{name: "canonical symbol skips expansion", symbol: "VWRA.L", want: []string{"VWRA.L"}},
		{name: "index symbol skips expansion", symbol: "^STI", want: []string{"^STI"}},
		{name: "alias resolves before expansion", symbol: "IHSG", want: []string{"^JKSE"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveTickerCandidates(tc.symbol)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d candidates, got %d (%#v)", len(tc.want), len(got), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("candidate %d: expected %q, got %q", i, tc.want[i], got[i])
				}
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
