package saaham

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gobot/config"
	"gobot/pkg"

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

var quoteService QuoteService = YahooQuoteService{}

func formatQuoteReply(result *QuoteResult) string {
	return fmt.Sprintf(
		"*%s* - %s\n`%.2f %s` `%+.2f (%+.2f%%)`",
		result.Symbol,
		result.InstrumentName,
		result.Price,
		result.Currency,
		result.Change,
		result.ChangePercent,
	)
}

func quoteCommand(c tele.Context) error {
	symbol, err := parseCommandLookup(c.Args())
	if len(c.Args()) == 0 {
		return c.Send("Usage: `/q AAPL`", &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}
	if err != nil {
		return c.Send("I couldn't find that exact ticker symbol.")
	}

	result, err := quoteService.Lookup(symbol)
	if err != nil {
		pkg.LogWithTimestamp("SAAHAM quote lookup failed for %s: %v", symbol, err)
		return sendLookupFailure(c, err)
	}

	return c.Send(formatQuoteReply(result), &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})
}

func plainTextQuoteLookup(c tele.Context) error {
	chat := c.Chat()
	if chat == nil || chat.Type != tele.ChatPrivate {
		return nil
	}

	symbol, ok := parsePlainTextLookup(c.Text())
	if !ok {
		return nil
	}

	result, err := quoteService.Lookup(symbol)
	if err != nil {
		pkg.LogWithTimestamp("SAAHAM plain-text quote lookup failed for %s: %v", symbol, err)
		return sendLookupFailure(c, err)
	}

	return c.Send(formatQuoteReply(result), &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})
}

func sendLookupFailure(c tele.Context, err error) error {
	var lookupErr *LookupError
	if errors.As(err, &lookupErr) {
		switch lookupErr.Kind {
		case LookupErrorInvalidSymbol:
			return c.Send("I couldn't find that exact ticker symbol.")
		case LookupErrorUnsupported:
			return c.Send("That ticker is valid, but SAAHAM Bot currently supports only stocks and ETFs.")
		}
	}

	return c.Send("Sorry, I couldn't fetch that quote right now.")
}

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

	pkg.LogWithTimestamp("SAAHAM bot started, listening for commands...")

	helpHandler := func(c tele.Context) error {
		return c.Send(helpMessage, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}

	bot.Handle("/q", quoteCommand)
	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)
	bot.Handle(tele.OnText, plainTextQuoteLookup)

	bot.Start()
}
