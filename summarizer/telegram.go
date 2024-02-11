package summarizer

import (
	"fmt"
	"gobot/config"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

const (
	regexUrl = `(?i)\b(?:https?://|www\.)\S+\b`

	helpMessage = `Hello this is Captain Kidd Bot

*Features:*
- Summarise web article
- Summarise paragraph text

*How to:*
Simply send an url or paragraph text to this bot

Captain Kidd bot made with ❤️ by @crossix`
)

func Run() {
	println("Run Summarizer")
	pref := tele.Settings{
		Token:  os.Getenv(config.CaptainKiddBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/s", func(c tele.Context) error {
		return c.Send(helpMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle(telebot.OnText, func(c tele.Context) error {
		pattern := regexp.MustCompile(regexUrl)
		if pattern.MatchString(c.Text()) {
			fmt.Println("url pattern")
			smmryResponse, errResponse := SummarizeURL(c.Text())
			if errResponse != nil {
				fmt.Println("Error:", errResponse.SmAPIMessage)
				return c.Send("Something went wrong...")
			}
			return c.Send(smmryResponse.ToMarkdownString(), &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
		} else {
			words := strings.Fields(c.Text())
			wordCount := len(words)

			if wordCount > 100 {
				fmt.Println("text pattern")
				smmryResponse, errResponse := SummarizeText(c.Text())
				if errResponse != nil {
					fmt.Println("Error:", errResponse.SmAPIMessage)
					return c.Send("Something went wrong...")
				}
				return c.Send(smmryResponse.ToMarkdownString(), &telebot.SendOptions{
					ParseMode: telebot.ModeMarkdown,
				})
			}
			return c.Send("Other regex")
		}
	})

	b.Start()
}
