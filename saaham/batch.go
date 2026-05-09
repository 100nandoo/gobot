package saaham

import (
	"errors"
	"strings"
)

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

type cachedLookupResult struct {
	result *QuoteResult
	err    error
}

func lookupQuote(service QuoteService, symbol string, cache map[string]cachedLookupResult) (*QuoteResult, error) {
	candidates := resolveTickerCandidates(symbol)
	var providerErr error
	var unsupportedErr error
	var invalidErr error

	for _, candidate := range candidates {
		result, err := lookupCandidate(service, candidate, cache)
		if err == nil {
			return result, nil
		}

		var lookupErr *LookupError
		if errors.As(err, &lookupErr) {
			switch lookupErr.Kind {
			case LookupErrorProviderFailed:
				if providerErr == nil {
					providerErr = err
				}
			case LookupErrorUnsupported:
				if unsupportedErr == nil {
					unsupportedErr = err
				}
			case LookupErrorInvalidSymbol:
				if invalidErr == nil {
					invalidErr = err
				}
			}
			continue
		}

		if providerErr == nil {
			providerErr = err
		}
	}

	switch {
	case providerErr != nil:
		return nil, providerErr
	case unsupportedErr != nil:
		return nil, unsupportedErr
	default:
		return nil, invalidErr
	}
}

func lookupCandidate(service QuoteService, candidate string, cache map[string]cachedLookupResult) (*QuoteResult, error) {
	if cached, ok := cache[candidate]; ok {
		return cached.result, cached.err
	}

	result, err := service.Lookup(candidate)
	cache[candidate] = cachedLookupResult{
		result: result,
		err:    err,
	}
	return result, err
}

func lookupBatchQuotes(service QuoteService, symbols []string) []BatchQuoteResult {
	cache := make(map[string]cachedLookupResult, len(symbols))
	seen := make(map[string]struct{}, len(symbols))
	results := make([]BatchQuoteResult, 0, len(symbols))
	for _, symbol := range symbols {
		result, err := lookupQuote(service, symbol, cache)
		key := dedupeBatchKey(symbol, result, err)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		results = append(results, BatchQuoteResult{
			Query:  symbol,
			Result: result,
			Err:    err,
		})
	}
	return results
}
