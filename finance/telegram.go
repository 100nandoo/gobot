package finance

import (
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

const helpMessage = `Finance Watchlist Bot

*Commands:*
- /a - Analyze saved list
- /w - Show watchlist
- /w add AAPL MSFT
- /w rm AAPL
- /w set AAPL MSFT NVDA
- /w reset
- /t - Show threshold
- /t set 5
- /t reset
- /s TSLA - Analyze one ticker
- /vwra
- /cspx
- /help - Show this message`

func scoreEmoji(score int) string {
	switch {
	case score >= 5:
		return "\xF0\x9F\x9F\xA2" // green circle
	case score >= 2:
		return "\xF0\x9F\x94\xB5" // blue circle
	case score >= -1:
		return "\xE2\x9A\xAA" // white circle
	case score >= -4:
		return "\xF0\x9F\x9F\xA0" // orange circle
	default:
		return "\xF0\x9F\x94\xB4" // red circle
	}
}

func changeArrow(change float64) string {
	if change >= 0 {
		return "\xE2\x96\xB2" // up triangle
	}
	return "\xE2\x96\xBC" // down triangle
}

func ptsEmoji(score int) string {
	if score > 0 {
		return "\xF0\x9F\x9F\xA2" // green
	}
	if score < 0 {
		return "\xF0\x9F\x94\xB4" // red
	}
	return "\xE2\x9A\xAA" // white
}

func fmtPts(s int) string {
	if s == 0 {
		return "0"
	}
	return fmt.Sprintf("%+d", s)
}

func shortSignal(s string) string {
	switch s {
	case "Bullish Crossover":
		return "Bullish Cross"
	case "Bearish Crossover":
		return "Bearish Cross"
	case "Below Both MAs":
		return "Below MAs"
	case "Above Both MAs":
		return "Above MAs"
	case "Strong Uptrend":
		return "Uptrend"
	case "Strong Downtrend":
		return "Downtrend"
	case "Near Lower Band":
		return "Near Lower"
	case "Near Upper Band":
		return "Near Upper"
	case "Below Lower Band":
		return "Below Lower"
	case "Above Upper Band":
		return "Above Upper"
	case "Within Bands":
		return "In Range"
	default:
		return s
	}
}

func formatAnalysis(etf ETFSymbol, r *IndicatorResult) string {
	sign := "+"
	if r.PriceChange < 0 {
		sign = ""
	}

	arrow := changeArrow(r.PriceChange)

	return fmt.Sprintf(
		"*%s* `%.2f %s` %s `%s%.2f%%`\n"+
			"\n"+
			"%s RSI `%.1f` %s `%s`\n"+
			"%s MACD `%+.2f` %s `%s`\n"+
			"%s SMA `%.1f / %.1f` %s `%s`\n"+
			"%s BB `%.1f / %.1f / %.1f` %s `%s`\n"+
			"\n"+
			"%s *%+d/10* _%s_",
		etf.Ticker, r.CurrentPrice, r.Currency, arrow, sign, r.PriceChange,
		ptsEmoji(r.RSIScore), r.RSI, shortSignal(r.RSISignal), fmtPts(r.RSIScore),
		ptsEmoji(r.MACDScore), r.MACDHistogram, shortSignal(r.MACDSignal), fmtPts(r.MACDScore),
		ptsEmoji(r.SMAScore), r.SMA50, r.SMA200, shortSignal(r.SMASignal), fmtPts(r.SMAScore),
		ptsEmoji(r.BBScore), r.BollingerLower, r.BollingerMid, r.BollingerUpper, shortSignal(r.BollingerSignal), fmtPts(r.BBScore),
		scoreEmoji(r.CompositeScore), r.CompositeScore, r.Recommendation,
	)
}

func formatWatchlist(symbols []ETFSymbol) string {
	if len(symbols) == 0 {
		return "Your watchlist is empty."
	}

	lines := make([]string, 0, len(symbols)+1)
	lines = append(lines, "*Current watchlist:*")
	for _, symbol := range symbols {
		lines = append(lines, "- `"+symbol.Ticker+"`")
	}

	return strings.Join(lines, "\n")
}

func formatThreshold(threshold int) string {
	return fmt.Sprintf("*Scout alert threshold:* `%+d`\nScheduled alerts are sent only when a ticker reaches this composite score or higher.", threshold)
}

func analyzeSymbols(symbols []ETFSymbol) string {
	var response string
	for _, etf := range symbols {
		result, err := AnalyzeETF(etf.Ticker)
		if err != nil {
			pkg.LogWithTimestamp("Error analyzing %s: %v", etf.Ticker, err)
			response += fmt.Sprintf("`Error analyzing %s`\n\n", etf.Ticker)
			continue
		}
		response += formatAnalysis(etf, result) + "\n\n"
	}

	return strings.TrimSpace(response)
}

func sendMarkdownChunks(c tele.Context, message string) error {
	const maxMessageSize = 3500

	if len(message) <= maxMessageSize {
		return c.Send(message, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}

	parts := strings.Split(message, "\n\n")
	var chunk string
	for _, part := range parts {
		candidate := part
		if chunk != "" {
			candidate = chunk + "\n\n" + part
		}

		if len(candidate) <= maxMessageSize {
			chunk = candidate
			continue
		}

		if chunk != "" {
			if err := c.Send(chunk, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown}); err != nil {
				return err
			}
		}
		chunk = part
	}

	if chunk == "" {
		return nil
	}

	return c.Send(chunk, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func analyzeAndRespond(c tele.Context, symbols []ETFSymbol) error {
	response := analyzeSymbols(symbols)
	return sendMarkdownChunks(c, response)
}

func thresholdCommand(c tele.Context) error {
	chat := c.Chat()
	if chat == nil {
		return c.Send("Unable to determine chat.")
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Send(formatThreshold(GetAlertThreshold(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}

	action := strings.ToLower(args[0])
	switch action {
	case "set":
		if len(args) < 2 {
			return c.Send("Usage: `/threshold set 5`", &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
		}

		value, err := strconv.Atoi(args[1])
		if err != nil {
			return c.Send("Threshold must be a whole number between -10 and 10.")
		}

		if err := SaveAlertThreshold(chat.ID, value); err != nil {
			return c.Send(err.Error())
		}

		return c.Send("Updated.\n\n"+formatThreshold(GetAlertThreshold(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	case "reset":
		if err := ResetAlertThreshold(chat.ID); err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Reset to default.\n\n"+formatThreshold(GetAlertThreshold(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	default:
		return c.Send("Unknown subcommand. Use `/t`, `/t set`, or `/t reset`.", &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}
}

func watchlistCommand(c tele.Context) error {
	chat := c.Chat()
	if chat == nil {
		return c.Send("Unable to determine chat.")
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Send(formatWatchlist(GetWatchlist(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}

	action := strings.ToLower(args[0])
	payload := args[1:]
	switch action {
	case "add":
		if len(payload) == 0 {
			return c.Send("Usage: `/w add AAPL MSFT`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		}
		if err := ValidateTickers(payload); err != nil {
			return c.Send(err.Error())
		}
		symbols, err := AddToWatchlist(chat.ID, payload)
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Updated.\n\n"+formatWatchlist(symbols), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	case "remove":
		if len(payload) == 0 {
			return c.Send("Usage: `/w rm AAPL`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		}
		symbols, err := RemoveFromWatchlist(chat.ID, payload)
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Updated.\n\n"+formatWatchlist(symbols), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	case "set":
		if len(payload) == 0 {
			return c.Send("Usage: `/w set AAPL MSFT NVDA`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		}
		if err := ValidateTickers(payload); err != nil {
			return c.Send(err.Error())
		}
		if err := SaveWatchlist(chat.ID, payload); err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Updated.\n\n"+formatWatchlist(GetWatchlist(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	case "reset":
		if err := ResetWatchlist(chat.ID); err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Reset to defaults.\n\n"+formatWatchlist(GetWatchlist(chat.ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	case "rm":
		if len(payload) == 0 {
			return c.Send("Usage: `/w rm AAPL`", &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		}
		symbols, err := RemoveFromWatchlist(chat.ID, payload)
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Updated.\n\n"+formatWatchlist(symbols), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	default:
		return c.Send("Unknown subcommand. Use `/w`, `/w add`, `/w rm`, `/w set`, or `/w reset`.", &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}
}

func analyzeTickerCommand(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Usage: `/s TSLA`", &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	}

	tickers := parseTickers(args)
	if len(tickers) == 0 {
		return c.Send("Please provide at least one valid ticker.")
	}

	if err := ValidateTickers(tickers); err != nil {
		return c.Send(err.Error())
	}

	return analyzeAndRespond(c, buildSymbols(tickers))
}

func Run() {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	pkg.LogWithTimestamp("Finance bot started, listening for commands...")

	b.Handle("/a", func(c tele.Context) error {
		pkg.LogWithTimestamp("Received /a command")
		return analyzeAndRespond(c, GetWatchlist(c.Chat().ID))
	})

	b.Handle("/vwra", func(c tele.Context) error {
		return analyzeAndRespond(c, []ETFSymbol{DefaultSymbols[0]})
	})

	b.Handle("/cspx", func(c tele.Context) error {
		return analyzeAndRespond(c, []ETFSymbol{DefaultSymbols[1]})
	})

	b.Handle("/w", watchlistCommand)

	b.Handle("/t", thresholdCommand)

	b.Handle("/s", analyzeTickerCommand)

	b.Handle("/stock", analyzeTickerCommand)

	b.Handle("/ticker", analyzeTickerCommand)

	b.Handle("/symbols", func(c tele.Context) error {
		return c.Send(formatWatchlist(GetWatchlist(c.Chat().ID)), &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Start()
}

func SendAnalysis(etf ETFSymbol, result *IndicatorResult) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.ChannelFinance), 10, 64)
	if chatIdErr != nil {
		pkg.LogWithTimestamp("Error parsing ETF channel ID: %v", chatIdErr)
		return
	}

	message := formatAnalysis(etf, result)

	_, sendErr := b.Send(tele.ChatID(num), message, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})

	if sendErr != nil {
		pkg.LogWithTimestamp("Error sending ETF analysis: %v", sendErr)
		return
	}
}
