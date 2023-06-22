package rss

import (
	"github.com/mmcdole/gofeed"
	"gobot/config"
	tele "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"time"
)

/*
SendRssItem

Send Rss item url to config.TelegramRemoteOk channel
*/
func SendRssItem(item gofeed.Item) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := tele.NewBot(pref)

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.TelegramRss), 10, 64)
	if chatIdErr != nil {
		println("Error from send items chatIdErr", chatIdErr)
		return
	}
	_, err := b.Send(tele.ChatID(num), item.Link)
	if err != nil {
		println("Error from send items", err)
		return
	}
}
