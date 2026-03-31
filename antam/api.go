package antam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GoldPrice represents the sell and buy prices for gold
type GoldPrice struct {
	Buy    string
	Sell   string
	Source string
}

// Get gold prices from the website. It returns a GoldPrice struct.
func getGoldPrices() (*GoldPrice, error) {
	resp, err := http.Get("https://harga-emas.org/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var perGramPrice string

	doc.Find("table.ComprehensiveTable_table__NjmlD tbody tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		unit := strings.TrimSpace(row.Find("td").First().Text())
		if unit != "Gram (gr)" {
			return true
		}

		priceCell := row.Find("td").Eq(2).Clone()
		priceCell.Find("span").Remove()
		perGramPrice = strings.TrimSpace(priceCell.Text())
		return false
	})

	if perGramPrice == "" {
		return nil, fmt.Errorf("failed to find IDR per gram price in harga-emas table")
	}

	return &GoldPrice{
		Buy:    perGramPrice,
		Sell:   perGramPrice,
		Source: "Harga-Emas.org",
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
