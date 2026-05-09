package saaham

import (
	"errors"
	"regexp"
	"strings"

	yfclient "github.com/wnjoon/go-yfinance/pkg/client"
)

type LookupErrorKind string

const (
	LookupErrorInvalidSymbol  LookupErrorKind = "invalid_symbol"
	LookupErrorProviderFailed LookupErrorKind = "provider_failure"
)

var tickerTokenPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9.\-=]*$`)

type LookupError struct {
	Kind   LookupErrorKind
	Symbol string
	Cause  error
}

func (e *LookupError) Error() string {
	if e.Cause == nil {
		return string(e.Kind)
	}
	return string(e.Kind) + ": " + e.Cause.Error()
}

func (e *LookupError) Unwrap() error {
	return e.Cause
}

func normalizeTickerToken(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.Trim(trimmed, "[](){}<>\"'`.,!?;:")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return "", false
	}
	if strings.HasPrefix(trimmed, "$") {
		return "", false
	}
	if len(strings.Fields(trimmed)) != 1 {
		return "", false
	}

	symbol := strings.ToUpper(trimmed)
	if !tickerTokenPattern.MatchString(symbol) {
		return "", false
	}

	return symbol, true
}

func parseCommandLookup(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}

	symbol, ok := normalizeTickerToken(strings.Join(args, " "))
	if !ok {
		return "", &LookupError{Kind: LookupErrorInvalidSymbol}
	}

	return symbol, nil
}

func parsePlainTextLookup(text string) (string, bool) {
	return normalizeTickerToken(text)
}

func classifyLookupError(symbol string, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, yfclient.ErrInvalidSymbol) || errors.Is(err, yfclient.ErrNotFound) || yfclient.IsInvalidSymbolError(err) || yfclient.IsNotFoundError(err) {
		return &LookupError{
			Kind:   LookupErrorInvalidSymbol,
			Symbol: symbol,
			Cause:  err,
		}
	}

	return &LookupError{
		Kind:   LookupErrorProviderFailed,
		Symbol: symbol,
		Cause:  err,
	}
}
