package saaham

import (
	"errors"
	"testing"
)

type stubQuoteService struct {
	results map[string]*QuoteResult
	errors  map[string]error
}

func (s stubQuoteService) Lookup(symbol string) (*QuoteResult, error) {
	if err := s.errors[symbol]; err != nil {
		return nil, err
	}
	if result := s.results[symbol]; result != nil {
		return result, nil
	}
	return nil, errors.New("unexpected symbol")
}

func TestLookupBatchQuotesPreservesInputOrderAndPartialResults(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"AAPL": {Symbol: "AAPL", InstrumentName: "Apple Inc."},
			"MSFT": {Symbol: "MSFT", InstrumentName: "Microsoft"},
		},
		errors: map[string]error{
			"BAD": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "BAD"},
		},
	}

	results := lookupBatchQuotes(service, []string{"AAPL", "BAD", "MSFT"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Query != "AAPL" || results[0].Result == nil {
		t.Fatal("expected first result to be AAPL success")
	}
	if results[1].Query != "BAD" || results[1].Err == nil {
		t.Fatal("expected second result to be BAD failure")
	}
	if results[2].Query != "MSFT" || results[2].Result == nil {
		t.Fatal("expected third result to be MSFT success")
	}
}

func TestParseBatchCommandLookupsDeduplicatesNormalizedSymbols(t *testing.T) {
	got := parseBatchCommandLookups([]string{"aapl", "(AAPL)", "msft!", "MSFT", "vwra.l"})

	if len(got) != 3 {
		t.Fatalf("expected 3 unique symbols, got %d", len(got))
	}
	if got[0] != "AAPL" || got[1] != "MSFT" || got[2] != "VWRA.L" {
		t.Fatalf("unexpected order/content: %#v", got)
	}
}
