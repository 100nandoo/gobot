package saaham

import (
	"testing"

	"github.com/wnjoon/go-yfinance/pkg/models"
)

func TestQuoteResultFromYahooUsesCanonicalFields(t *testing.T) {
	quote := &models.Quote{
		Symbol:                     "aapl",
		LongName:                   "Apple Inc.",
		Currency:                   "USD",
		RegularMarketPrice:         210.15,
		RegularMarketChange:        3.14,
		RegularMarketChangePercent: 1.52,
		QuoteType:                  "EQUITY",
	}

	got := quoteResultFromYahoo(quote)

	if got.Symbol != "AAPL" {
		t.Fatalf("expected symbol AAPL, got %q", got.Symbol)
	}
	if got.InstrumentName != "Apple Inc." {
		t.Fatalf("expected instrument name Apple Inc., got %q", got.InstrumentName)
	}
	if got.Currency != "USD" {
		t.Fatalf("expected currency USD, got %q", got.Currency)
	}
	if got.Price != 210.15 {
		t.Fatalf("expected price 210.15, got %v", got.Price)
	}
	if got.Change != 3.14 {
		t.Fatalf("expected change 3.14, got %v", got.Change)
	}
	if got.ChangePercent != 1.52 {
		t.Fatalf("expected change percent 1.52, got %v", got.ChangePercent)
	}
}

func TestIsSupportedInstrumentType(t *testing.T) {
	tests := []struct {
		quoteType string
		want      bool
	}{
		{quoteType: "EQUITY", want: true},
		{quoteType: "ETF", want: true},
		{quoteType: "INDEX", want: true},
		{quoteType: "CURRENCY", want: true},
		{quoteType: "MUTUALFUND", want: false},
		{quoteType: "CRYPTOCURRENCY", want: false},
		{quoteType: "", want: false},
	}

	for _, tc := range tests {
		if got := isSupportedInstrumentType(tc.quoteType); got != tc.want {
			t.Fatalf("quoteType %q: expected %v, got %v", tc.quoteType, tc.want, got)
		}
	}
}
