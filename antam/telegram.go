package antam

import (
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

const (
	startMessage = `Halo, ini adalah bot cek harga emas Antam.

- /a untuk harga emas antam dari hargaemas.com
- /p untuk harga emas digital dari pluang.com
- /h untuk harga emas digital dari harga-emas.org
- /help untuk lihat bantuan ini lagi

Catatan:
- Beli = harga saat kamu membeli emas
- Jual = harga saat kamu menjual kembali emas`

	helpMessage = `Halo, ini adalah bot cek harga emas Antam.

- /a untuk harga emas antam dari hargaemas.com
- /p untuk harga emas digital dari pluang.com
- /h untuk harga emas digital dari harga-emas.org
- /help untuk lihat bantuan ini lagi

Catatan:
- Beli = harga saat kamu membeli emas
- Jual = harga saat kamu menjual kembali emas`
)

// Helper function to format the gold price response message
func formatGoldPriceResponse(price GoldPrice) string {
	buyLabel := price.BuyLabel
	if buyLabel == "" {
		buyLabel = "Beli"
	}

	sellLabel := price.SellLabel
	if sellLabel == "" {
		sellLabel = "Jual"
	}

	return fmt.Sprintf("`%s\n%s: %s\n%s: %s`", price.Source, buyLabel, price.Buy, sellLabel, price.Sell)
}

func formatGoldBuyPriceResponse(price GoldPrice) string {
	return fmt.Sprintf("%s\n*%s*", price.Source, price.Buy)
}

func Run() {
	pref := tele.Settings{
		Token:  os.Getenv(config.AntamTelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	sendPrice := func(c tele.Context, fetch func() (*GoldPrice, error)) error {
		price, err := fetch()
		if err != nil {
			pkg.LogWithTimestamp("Error fetching gold prices: %v", err)
			return c.Send("Sorry, I couldn't fetch the gold prices right now.", &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
		}

		responseMessage := formatGoldPriceResponse(*price)

		return c.Send(responseMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(startMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/h", func(c tele.Context) error {
		return sendPrice(c, getGoldPrices)
	})

	b.Handle("/a", func(c tele.Context) error {
		return sendPrice(c, getHargaEmasComPrices)
	})

	b.Handle("/p", func(c tele.Context) error {
		return sendPrice(c, getPluangGoldPrices)
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Start()
}

func sendPricesToChannel(formatter func(GoldPrice) string, prices ...GoldPrice) {
	pref := tele.Settings{
		Token:  os.Getenv(config.AntamTelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.ChannelAntam), 10, 64)
	if chatIdErr != nil {
		pkg.LogWithTimestamp("Error from send price chatIdErr: %v", chatIdErr)
		return
	}

	var responseMessage string

	for _, price := range prices {
		responseMessage += formatter(price) + "\n\n"
	}

	_, sendErr := b.Send(tele.ChatID(num), responseMessage, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	if sendErr != nil {
		pkg.LogWithTimestamp("Error from send price: %v", err)
		return
	}
}

/*
SendPrice

Send antam gold price to config.ChannelAntam channel
*/
func SendPrice(prices ...GoldPrice) {
	sendPricesToChannel(formatGoldPriceResponse, prices...)
}

/*
SendBuyPrice

Send antam buy-only gold price to config.ChannelAntam channel
*/
func SendBuyPrice(prices ...GoldPrice) {
	sendPricesToChannel(formatGoldBuyPriceResponse, prices...)
}
