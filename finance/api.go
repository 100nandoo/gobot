package finance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency           string  `json:"currency"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				PreviousClose      float64 `json:"previousClose"`
			} `json:"meta"`
			Indicators struct {
				Quote []struct {
					Close []any `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

func fetchPriceData(symbol string) (*PriceData, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=1y&interval=1d", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo finance returned status %d for %s", resp.StatusCode, symbol)
	}

	var chart yahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&chart); err != nil {
		return nil, fmt.Errorf("decoding response for %s: %w", symbol, err)
	}

	if chart.Chart.Error != nil {
		return nil, fmt.Errorf("yahoo finance error for %s: %s", symbol, chart.Chart.Error.Description)
	}

	if len(chart.Chart.Result) == 0 || len(chart.Chart.Result[0].Indicators.Quote) == 0 {
		return nil, fmt.Errorf("no data returned for %s", symbol)
	}

	result := chart.Chart.Result[0]
	rawClose := result.Indicators.Quote[0].Close

	var closes []float64
	for _, v := range rawClose {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case float64:
			closes = append(closes, val)
		case json.Number:
			f, err := val.Float64()
			if err == nil {
				closes = append(closes, f)
			}
		}
	}

	if len(closes) < 200 {
		return nil, fmt.Errorf("insufficient price data for %s: got %d days, need at least 200", symbol, len(closes))
	}

	return &PriceData{
		Close:     closes,
		Currency:  result.Meta.Currency,
		Price:     result.Meta.RegularMarketPrice,
		PrevClose: result.Meta.PreviousClose,
	}, nil
}
