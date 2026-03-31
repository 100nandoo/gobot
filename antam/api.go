package antam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// GoldPrice represents the sell and buy prices for gold
type GoldPrice struct {
	Buy    string
	Sell   string
	Source string
}

// Get gold prices from the website. It returns a GoldPrice struct.
func getGoldPricesFromHTML() (*GoldPrice, error) {
	resp, err := http.Get("https://harga-emas.org/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var sellPrice, buyPrice string

	// Find the first table with the class "in_table"
	table := doc.Find(".in_table").First()

	// Select the row which contains the prices (4th row)
	priceRow := table.Find("tr").Eq(3)

	// Get the last two <td> elements for sell and buy prices
	buyPrice = priceRow.Find("td").Eq(8).Text()
	sellPrice = priceRow.Find("td").Eq(9).Text()

	return &GoldPrice{
		Buy:    buyPrice,
		Sell:   sellPrice,
		Source: "Gedung Antam Jakarta",
	}, nil
}

func getPluangGoldPrices() (*GoldPrice, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api-pluang.pluang.com/api/v3/asset/gold/pricing?daysLimit=1", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:149.0) Gecko/20100101 Firefox/149.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://pluang.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pluang pricing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from pluang api: %d", resp.StatusCode)
	}

	var payload struct {
		StatusCode int `json:"statusCode"`
		Data       struct {
			Current struct {
				Sell int64 `json:"sell"`
				Buy  int64 `json:"buy"`
			} `json:"current"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode pluang pricing response: %w", err)
	}

	if payload.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pluang api returned statusCode %d", payload.StatusCode)
	}

	return &GoldPrice{
		Buy:    formatPriceIDR(payload.Data.Current.Buy),
		Sell:   formatPriceIDR(payload.Data.Current.Sell),
		Source: "Pluang.com",
	}, nil
}

func formatPriceIDR(price int64) string {
	digits := strconv.FormatInt(price, 10)
	if len(digits) <= 3 {
		return digits
	}

	result := make([]byte, 0, len(digits)+(len(digits)-1)/3)
	leading := len(digits) % 3
	if leading == 0 {
		leading = 3
	}

	result = append(result, digits[:leading]...)
	for i := leading; i < len(digits); i += 3 {
		result = append(result, '.')
		result = append(result, digits[i:i+3]...)
	}

	return string(result)
}
