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

func TestLookupQuoteUsesShortcutFallbackOrder(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"CSPX.L": {Symbol: "CSPX.L", InstrumentName: "iShares Core S&P 500 UCITS ETF"},
		},
		errors: map[string]error{
			"^CSPX":   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "^CSPX"},
			"CSPX.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "CSPX.JK"},
		},
	}

	result, err := lookupQuote(service, "CSPX", map[string]cachedLookupResult{})
	if err != nil {
		t.Fatalf("expected fallback success, got %v", err)
	}
	if result == nil || result.Symbol != "CSPX.L" {
		t.Fatalf("expected CSPX.L result, got %#v", result)
	}
}

func TestLookupQuoteUsesCurrencyPairShortcutFirst(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"USDSGD=X": {Symbol: "USDSGD=X", InstrumentName: "USD/SGD"},
		},
		errors: map[string]error{
			"^USDSGD":   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "^USDSGD"},
			"USDSGD.L":  &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "USDSGD.L"},
			"USDSGD.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "USDSGD.JK"},
			"USDSGD":    &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "USDSGD"},
		},
	}

	result, err := lookupQuote(service, "USDSGD", map[string]cachedLookupResult{})
	if err != nil {
		t.Fatalf("expected fx shortcut success, got %v", err)
	}
	if result == nil || result.Symbol != "USDSGD=X" {
		t.Fatalf("expected USDSGD=X result, got %#v", result)
	}
}

func TestLookupQuoteUsesBareFallbackForNonSP100Symbols(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"SNDK": {Symbol: "SNDK", InstrumentName: "Sandisk Corp."},
		},
		errors: map[string]error{
			"^SNDK":   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "^SNDK"},
			"SNDK.L":  &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "SNDK.L"},
			"SNDK.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "SNDK.JK"},
		},
	}

	result, err := lookupQuote(service, "SNDK", map[string]cachedLookupResult{})
	if err != nil {
		t.Fatalf("expected bare fallback success, got %v", err)
	}
	if result == nil || result.Symbol != "SNDK" {
		t.Fatalf("expected SNDK result, got %#v", result)
	}
}

func TestLookupQuoteUsesBareFirstForSP100Members(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"AAPL": {Symbol: "AAPL", InstrumentName: "Apple Inc."},
		},
		errors: map[string]error{
			"^AAPL":   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "^AAPL"},
			"AAPL.L":  &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "AAPL.L"},
			"AAPL.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "AAPL.JK"},
		},
	}

	result, err := lookupQuote(service, "AAPL", map[string]cachedLookupResult{})
	if err != nil {
		t.Fatalf("expected bare symbol success, got %v", err)
	}
	if result == nil || result.Symbol != "AAPL" {
		t.Fatalf("expected AAPL result, got %#v", result)
	}
}

func TestLookupQuotePrefersProviderFailureOverUnsupportedAndInvalid(t *testing.T) {
	service := stubQuoteService{
		errors: map[string]error{
			"STI":    &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "STI"},
			"^STI":   &LookupError{Kind: LookupErrorUnsupported, Symbol: "^STI"},
			"STI.L":  &LookupError{Kind: LookupErrorProviderFailed, Symbol: "STI.L"},
			"STI.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "STI.JK"},
		},
	}

	_, err := lookupQuote(service, "STI", map[string]cachedLookupResult{})
	var lookupErr *LookupError
	if !errors.As(err, &lookupErr) {
		t.Fatalf("expected lookup error, got %v", err)
	}
	if lookupErr.Kind != LookupErrorProviderFailed || lookupErr.Symbol != "STI.L" {
		t.Fatalf("expected STI.L provider failure, got %#v", lookupErr)
	}
}

func TestLookupBatchQuotesDeduplicatesResolvedCanonicalSymbols(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"^STI":   {Symbol: "^STI", InstrumentName: "STI Index"},
			"CSPX.L": {Symbol: "CSPX.L", InstrumentName: "iShares Core S&P 500 UCITS ETF"},
			"^JKSE":  {Symbol: "^JKSE", InstrumentName: "Jakarta Composite Index"},
		},
		errors: map[string]error{
			"STI":     &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "STI"},
			"^CSPX":   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "^CSPX"},
			"CSPX.JK": &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "CSPX.JK"},
		},
	}

	results := lookupBatchQuotes(service, []string{"sti", "^sti", "cspx", "cspx.l", "ihsg"})
	if len(results) != 3 {
		t.Fatalf("expected 3 canonical results, got %d", len(results))
	}
	if results[0].Result == nil || results[0].Result.Symbol != "^STI" {
		t.Fatalf("expected first result to resolve to ^STI, got %#v", results[0])
	}
	if results[1].Result == nil || results[1].Result.Symbol != "CSPX.L" {
		t.Fatalf("expected second result to resolve to CSPX.L, got %#v", results[1])
	}
	if results[2].Result == nil || results[2].Result.Symbol != "^JKSE" {
		t.Fatalf("expected third result to resolve to ^JKSE, got %#v", results[2])
	}
}
