package antam

import (
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"log"
	"os"
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
		prices, err := getGoldPricesFromHTML()
		if err != nil {
			pkg.LogWithTimestamp("Error fetching gold prices: %v", err)
			return c.Send("Sorry, I couldn't fetch the gold prices right now.", &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
		}

		// Prepare the response message with prices
		responseMessage := fmt.Sprintf("`Harga Emas Antam:\n\nBeli: %s\nJual: %s`", prices.Buy, prices.Sell)

		return c.Send(responseMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/p", func(c tele.Context) error {
		prices, err := getPluangGoldPricesFromHTML() // Fetch gold prices
		if err != nil {
			pkg.LogWithTimestamp("Error fetching gold prices: %v", err)
			return c.Send("Sorry, I couldn't fetch the gold prices right now.", &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
		}

		// Prepare the response message with prices
		responseMessage := fmt.Sprintf("`Harga di Pluang:\n\nBeli: %s\nJual: %s`", prices.Buy, prices.Sell)

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
