package antam

import (
	"fmt"
	"net/http"
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

func getPluangGoldPricesFromHTML() (*GoldPrice, error) {
	resp, err := http.Get("https://pluang.com/widgets/price-graph/desktop-vertical")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var sellPrice, buyPrice string

	doc.Find(".halfwidth").Each(func(i int, s *goquery.Selection) {
		// Find the <p> element within the current <div>
		s.Find("p").Each(func(j int, p *goquery.Selection) {
			text := strings.TrimSpace(p.Text())
			// Remove the "/g" suffix from the prices
			text = strings.ReplaceAll(text, "/g", "")
			text = strings.ReplaceAll(text, "Rp", "")
			if i == 0 {
				sellPrice = text
			} else if i == 1 {
				buyPrice = text
			}
		})
	})

	return &GoldPrice{
		Buy:    buyPrice,
		Sell:   sellPrice,
		Source: "Pluang.com",
	}, nil
}
