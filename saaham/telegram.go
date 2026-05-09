package saaham

import (
	"log"
	"os"
	"time"

	"gobot/config"

	tele "gopkg.in/telebot.v3"
)

const helpMessage = `SAAHAM Bot

Fast stock and ETF quote lookup.

*Commands:*
- /q AAPL - Get the latest quote for one symbol
- /help - Show this message

*Notes:*
- Use exact Yahoo Finance ticker symbols
- SAAHAM Bot is a quote-only bot`

func Run() {
	pref := tele.Settings{
		Token:  os.Getenv(config.SaahamBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	helpHandler := func(c tele.Context) error {
		return c.Send(helpMessage, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}

	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)

	bot.Start()
}
