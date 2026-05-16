package saaham

import (
	"fmt"
	"strings"

	"github.com/wnjoon/go-yfinance/pkg/models"
	"github.com/wnjoon/go-yfinance/pkg/ticker"
)

type QuoteResult struct {
	Symbol         string
	InstrumentName string
	Currency       string
	Price          float64
	Change         float64
	ChangePercent  float64
	QuoteType      string
}

type QuoteService interface {
	Lookup(symbol string) (*QuoteResult, error)
}

type YahooQuoteService struct{}

func (s YahooQuoteService) Lookup(symbol string) (*QuoteResult, error) {
	t, err := ticker.New(symbol)
	if err != nil {
		return nil, classifyLookupError(symbol, fmt.Errorf("create ticker: %w", err))
	}
	defer t.Close()

	quote, err := t.Quote()
	if err != nil {
		return nil, classifyLookupError(symbol, fmt.Errorf("fetch quote: %w", err))
	}
	if !isSupportedInstrumentType(quote.QuoteType) {
		return nil, &LookupError{
			Kind:   LookupErrorUnsupported,
			Symbol: strings.ToUpper(strings.TrimSpace(quote.Symbol)),
		}
	}

	return quoteResultFromYahoo(quote), nil
}

func isSupportedInstrumentType(quoteType string) bool {
	switch strings.ToUpper(strings.TrimSpace(quoteType)) {
	case "EQUITY", "ETF", "INDEX", "CURRENCY":
		return true
	default:
		return false
	}
}

func quoteResultFromYahoo(quote *models.Quote) *QuoteResult {
	name := strings.TrimSpace(quote.LongName)
	if name == "" {
		name = strings.TrimSpace(quote.ShortName)
	}
	if name == "" {
		name = quote.Symbol
	}

	return &QuoteResult{
		Symbol:         strings.ToUpper(strings.TrimSpace(quote.Symbol)),
		InstrumentName: name,
		Currency:       quote.Currency,
		Price:          quote.RegularMarketPrice,
		Change:         quote.RegularMarketChange,
		ChangePercent:  quote.RegularMarketChangePercent,
		QuoteType:      quote.QuoteType,
	}
}
