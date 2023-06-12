package freegames

import (
	"gobot/config"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

/*
SendPost

Send reddit post url to config.TelegramFreeGames channel
*/
func SendPost(post Post) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := tele.NewBot(pref)

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.TelegramFreeGames), 10, 64)
	if chatIdErr != nil {
		println("Error from send post chatIdErr", chatIdErr)
		return
	}
	_, err := b.Send(tele.ChatID(num), post.URL)
	if err != nil {
		println("Error from send post", err)
		return
	}
}
