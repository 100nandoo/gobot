package saaham

import "testing"

func TestTickerSnapshotContainsNormalizedSymbols(t *testing.T) {
	snapshot := loadTickerSnapshot("AAPL\ngoog\nBRK.B\n\n")

	tests := []struct {
		symbol string
		want   bool
	}{
		{symbol: "aapl", want: true},
		{symbol: "GOOG", want: true},
		{symbol: "brk.b", want: true},
		{symbol: "MSFT", want: false},
		{symbol: "$AAPL", want: false},
	}

	for _, tc := range tests {
		if got := snapshot.Contains(tc.symbol); got != tc.want {
			t.Fatalf("symbol %q: expected %v, got %v", tc.symbol, tc.want, got)
		}
	}
}

func TestSP100SnapshotIncludesRepresentativeMembers(t *testing.T) {
	tests := []string{"AAPL", "GOOG", "GOOGL", "BRK.B"}

	for _, symbol := range tests {
		if !sp100Snapshot.Contains(symbol) {
			t.Fatalf("expected snapshot to include %q", symbol)
		}
	}
}
