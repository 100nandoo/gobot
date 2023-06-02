package telegram

import (
	"gobot/config"
	"gobot/freegames/reddit"
	"time"

	tele "gopkg.in/telebot.v3"
)

/*
SendPost

Send reddit post url to config.TelegramFreeGames channel
*/
func SendPost(post reddit.Post) {
	pref := tele.Settings{
		Token:  config.TelegramAaron,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := tele.NewBot(pref)
	_, err := b.Send(tele.ChatID(config.TelegramFreeGames), post.URL)
	if err != nil {
		println("Error from send post", err)
		return
	}
}
