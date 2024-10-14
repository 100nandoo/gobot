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

	helpMessage = `Halo ini adalah bot cek harga emas antam

*Cara penggunaan:*
- Kirim /start untuk cek harga jual beli emas antam
- Kirim /p untuk cek harga jual beli emas antam di pluang

Emas Antam Bot dibuat dengan ❤️ oleh @crossix`
)

// Helper function to format the gold price response message
func formatGoldPriceResponse(price GoldPrice) string {
	return fmt.Sprintf("`Harga Emas:\n\nBeli: %s\nJual: %s`", price.Buy, price.Sell)
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

	b.Handle("/start", func(c tele.Context) error {
		price, err := getGoldPricesFromHTML()
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
	})

	b.Handle("/p", func(c tele.Context) error {
		price, err := getPluangGoldPricesFromHTML() // Fetch gold prices
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
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Start()
}

/*
SendPrice

Send antam gold price to config.ChannelAntam channel
*/
func SendPrice(price GoldPrice) {
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

	responseMessage := formatGoldPriceResponse(price)

	_, sendErr := b.Send(tele.ChatID(num), responseMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
	})
	
	if sendErr != nil {
		pkg.LogWithTimestamp("Error from send price: %v", err)
		return
	}
}
