package saaham

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gobot/config"
	"gobot/pkg"

	tele "gopkg.in/telebot.v3"
)

const helpMessage = `SAAHAM Bot

Fast stock, ETF, index, and currency pair quote lookup.

*Commands:*
- /q AAPL - Get the latest quote for one symbol
- /fx 100 SGD IDR - Convert an amount using the latest FX rate
- /help - Show this message

*Notes:*
- You can use exact Yahoo Finance ticker symbols or supported shortcuts like STI, CSPX, VWRA, BJBR, IHSG, and USDSGD
- In groups, you can send /q AAPL, /fx 100 SGD IDR, or a single ticker like AAPL
- Group replies quote the triggering message
- SAAHAM Bot is a quote-only bot`

var quoteService QuoteService = YahooQuoteService{}

func allowsCommandLookup(chat *tele.Chat) bool {
	return chat != nil
}

func allowsPlainTextLookup(chat *tele.Chat) bool {
	if chat == nil {
		return false
	}

	switch chat.Type {
	case tele.ChatPrivate, tele.ChatGroup, tele.ChatSuperGroup:
		return true
	default:
		return false
	}
}

func usesQuotedGroupReply(chat *tele.Chat) bool {
	if chat == nil {
		return false
	}

	switch chat.Type {
	case tele.ChatGroup, tele.ChatSuperGroup:
		return true
	default:
		return false
	}
}

func shouldIgnoreTrigger(c tele.Context) bool {
	sender := c.Sender()
	return sender != nil && sender.IsBot
}

func sendQuoteMessage(c tele.Context, text string) error {
	if usesQuotedGroupReply(c.Chat()) {
		return c.Reply(text, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}

	return c.Send(text, &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})
}

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

func formatBatchQuoteReply(results []BatchQuoteResult) string {
	blocks := make([]string, 0, len(results))
	for _, item := range results {
		if item.Err != nil {
			blocks = append(blocks, formatBatchLookupFailure(item.Query, item.Err))
			continue
		}
		blocks = append(blocks, formatQuoteReply(item.Result))
	}

	return strings.Join(blocks, "\n\n")
}

func formatBatchLookupFailure(symbol string, err error) string {
	var lookupErr *LookupError
	switch {
	case errors.As(err, &lookupErr) && lookupErr.Kind == LookupErrorInvalidSymbol:
		return fmt.Sprintf("*%s* - invalid symbol", symbol)
	case errors.As(err, &lookupErr) && lookupErr.Kind == LookupErrorUnsupported:
		return fmt.Sprintf("*%s* - unsupported instrument", symbol)
	default:
		return fmt.Sprintf("*%s* - provider failure", symbol)
	}
}

func quoteCommand(c tele.Context) error {
	if !allowsCommandLookup(c.Chat()) {
		return nil
	}
	if shouldIgnoreTrigger(c) {
		return nil
	}

	if len(c.Args()) == 0 {
		return sendQuoteMessage(c, "Usage: `/q AAPL`")
	}

	symbols := parseBatchCommandLookups(c.Args())
	if len(symbols) == 0 {
		return sendQuoteMessage(c, "Usage: `/q AAPL`")
	}
	if len(symbols) > maxBatchSymbols {
		return sendQuoteMessage(c, "Usage: `/q AAPL MSFT NVDA TSLA AMZN` (up to 5 symbols)")
	}
	if len(symbols) == 1 {
		result, err := lookupQuote(quoteService, symbols[0], map[string]cachedLookupResult{})
		if err != nil {
			pkg.LogWithTimestamp("SAAHAM quote lookup failed for %s: %v", symbols[0], err)
			return sendLookupFailure(c, err)
		}

		return sendQuoteMessage(c, formatQuoteReply(result))
	}

	results := lookupBatchQuotes(quoteService, symbols)
	for _, item := range results {
		if item.Err != nil {
			pkg.LogWithTimestamp("SAAHAM batch quote lookup failed for %s: %v", item.Query, item.Err)
		}
	}

	return sendQuoteMessage(c, formatBatchQuoteReply(results))
}

func conversionCommand(c tele.Context) error {
	if !allowsCommandLookup(c.Chat()) {
		return nil
	}
	if shouldIgnoreTrigger(c) {
		return nil
	}

	request, err := parseConversionCommand(c.Args())
	if err != nil || request == nil {
		return sendQuoteMessage(c, conversionUsage)
	}

	result, convertedAmount, err := lookupConversion(quoteService, request)
	if err != nil {
		pkg.LogWithTimestamp("SAAHAM conversion lookup failed for %s: %v", request.PairSymbol(), err)
		return sendConversionFailure(c, err)
	}

	return sendQuoteMessage(c, formatConversionReply(request, result, convertedAmount))
}

func plainTextQuoteLookup(c tele.Context) error {
	if !allowsPlainTextLookup(c.Chat()) {
		return nil
	}
	if shouldIgnoreTrigger(c) {
		return nil
	}
	if c.Chat() != nil && c.Chat().Type == tele.ChatPrivate {
		if request, ok := parsePlainTextConversion(c.Text()); ok {
			result, convertedAmount, err := lookupConversion(quoteService, request)
			if err != nil {
				pkg.LogWithTimestamp("SAAHAM plain-text conversion lookup failed for %s: %v", request.PairSymbol(), err)
				return sendConversionFailure(c, err)
			}

			return sendQuoteMessage(c, formatConversionReply(request, result, convertedAmount))
		}
	}

	symbol, ok := parsePlainTextLookup(c.Text())
	if !ok {
		return nil
	}

	result, err := lookupQuote(quoteService, symbol, map[string]cachedLookupResult{})
	if err != nil {
		pkg.LogWithTimestamp("SAAHAM plain-text quote lookup failed for %s: %v", symbol, err)
		return sendLookupFailure(c, err)
	}

	return sendQuoteMessage(c, formatQuoteReply(result))
}

func sendConversionFailure(c tele.Context, err error) error {
	var lookupErr *LookupError
	if errors.As(err, &lookupErr) {
		switch lookupErr.Kind {
		case LookupErrorInvalidSymbol:
			return sendQuoteMessage(c, "I couldn't find that currency pair for conversion.")
		case LookupErrorUnsupported:
			return sendQuoteMessage(c, "That currency pair is valid, but SAAHAM Bot couldn't convert it right now.")
		}
	}

	return sendQuoteMessage(c, "Sorry, I couldn't fetch that conversion rate right now.")
}

func sendLookupFailure(c tele.Context, err error) error {
	var lookupErr *LookupError
	if errors.As(err, &lookupErr) {
		switch lookupErr.Kind {
		case LookupErrorInvalidSymbol:
			return sendQuoteMessage(c, "I couldn't find that exact ticker symbol.")
		case LookupErrorUnsupported:
			return sendQuoteMessage(c, "That ticker is valid, but SAAHAM Bot currently supports only stocks, ETFs, indices, and currency pairs.")
		}
	}

	return sendQuoteMessage(c, "Sorry, I couldn't fetch that quote right now.")
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
	bot.Handle("/fx", conversionCommand)
	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)
	bot.Handle(tele.OnText, plainTextQuoteLookup)

	bot.Start()
}
