package saaham

import "strings"

const maxBatchSymbols = 5

type BatchQuoteResult struct {
	Query  string
	Result *QuoteResult
	Err    error
}

func parseBatchCommandLookups(args []string) []string {
	symbols := make([]string, 0, len(args))
	seen := make(map[string]struct{}, len(args))
	for _, arg := range args {
		symbol, ok := normalizeTickerToken(arg)
		if ok {
			if _, exists := seen[symbol]; exists {
				continue
			}
			seen[symbol] = struct{}{}
			symbols = append(symbols, symbol)
			continue
		}

		fallback := strings.ToUpper(strings.TrimSpace(arg))
		if _, exists := seen[fallback]; exists {
			continue
		}
		seen[fallback] = struct{}{}
		symbols = append(symbols, fallback)
	}
	return symbols
}

func lookupBatchQuotes(service QuoteService, symbols []string) []BatchQuoteResult {
	results := make([]BatchQuoteResult, 0, len(symbols))
	for _, symbol := range symbols {
		result, err := service.Lookup(symbol)
		results = append(results, BatchQuoteResult{
			Query:  symbol,
			Result: result,
			Err:    err,
		})
	}
	return results
}
