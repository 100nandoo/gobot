package saaham

import "strings"

type BatchQuoteResult struct {
	Query  string
	Result *QuoteResult
	Err    error
}

func parseBatchCommandLookups(args []string) []string {
	symbols := make([]string, 0, len(args))
	for _, arg := range args {
		symbol, ok := normalizeTickerToken(arg)
		if ok {
			symbols = append(symbols, symbol)
			continue
		}

		symbols = append(symbols, strings.ToUpper(strings.TrimSpace(arg)))
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
