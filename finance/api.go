package finance

import (
	"fmt"

	"github.com/wnjoon/go-yfinance/pkg/models"
	"github.com/wnjoon/go-yfinance/pkg/ticker"
)

type ETFSymbol struct {
	Ticker string
	Name   string
}

var DefaultSymbols = []ETFSymbol{
	{Ticker: "VWRA.L", Name: "Vanguard FTSE All-World"},
	{Ticker: "CSPX.L", Name: "iShares Core S&P 500"},
}

type PriceData struct {
	Close     []float64
	Currency  string
	Price     float64
	PrevClose float64
}

func fetchPriceData(symbol string) (*PriceData, error) {
	t, err := ticker.New(symbol)
	if err != nil {
		return nil, fmt.Errorf("creating ticker for %s: %w", symbol, err)
	}
	defer t.Close()

	quote, err := t.Quote()
	if err != nil {
		return nil, fmt.Errorf("fetching quote for %s: %w", symbol, err)
	}

	bars, err := t.History(models.HistoryParams{
		Period:     "1y",
		Interval:   "1d",
		AutoAdjust: false,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching history for %s: %w", symbol, err)
	}

	closes := make([]float64, 0, len(bars))
	for _, bar := range bars {
		if bar.Close == 0 {
			continue
		}
		closes = append(closes, bar.Close)
	}

	if len(closes) < 200 {
		return nil, fmt.Errorf("insufficient price data for %s: got %d days, need at least 200", symbol, len(closes))
	}

	currency := quote.Currency
	if currency == "" {
		if meta := t.GetHistoryMetadata(); meta != nil {
			currency = meta.Currency
		}
	}

	return &PriceData{
		Close:     closes,
		Currency:  currency,
		Price:     quote.RegularMarketPrice,
		PrevClose: quote.RegularMarketPreviousClose,
	}, nil
}
