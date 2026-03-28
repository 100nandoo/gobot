package finance

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

const helpMessage = `ETF Entry Reminder Bot

*Commands:*
- /etf - Analyze both VWRA and CSPX
- /vwra - Analyze VWRA (Vanguard FTSE All-World)
- /cspx - Analyze CSPX (iShares Core S&P 500)
- /help - Show this message`

func scoreEmoji(score int) string {
	switch {
	case score >= 5:
		return "\xF0\x9F\x9F\xA2" // green circle
	case score >= 2:
		return "\xF0\x9F\x94\xB5" // blue circle
	case score >= -1:
		return "\xE2\x9A\xAA"     // white circle
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

func analyzeAndRespond(c tele.Context, etfs []ETFSymbol) error {
	var response string
	for _, etf := range etfs {
		result, err := AnalyzeETF(etf.Ticker)
		if err != nil {
			pkg.LogWithTimestamp("Error analyzing %s: %v", etf.Ticker, err)
			response += fmt.Sprintf("`Error analyzing %s`\n\n", etf.Ticker)
			continue
		}
		response += formatAnalysis(etf, result) + "\n\n"
	}

	return c.Send(response, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
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

	b.Handle("/etf", func(c tele.Context) error {
		pkg.LogWithTimestamp("Received /etf command")
		return analyzeAndRespond(c, ETFs)
	})

	b.Handle("/vwra", func(c tele.Context) error {
		return analyzeAndRespond(c, []ETFSymbol{ETFs[0]})
	})

	b.Handle("/cspx", func(c tele.Context) error {
		return analyzeAndRespond(c, []ETFSymbol{ETFs[1]})
	})

	b.Handle("/start", func(c tele.Context) error {
		return analyzeAndRespond(c, ETFs)
	})

	b.Handle("/help", func(c tele.Context) error {
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
	silent := result.CompositeScore < 5

	_, sendErr := b.Send(tele.ChatID(num), message, &telebot.SendOptions{
		ParseMode:           telebot.ModeMarkdown,
		DisableNotification: silent,
	})

	if sendErr != nil {
		pkg.LogWithTimestamp("Error sending ETF analysis: %v", sendErr)
		return
	}
}
