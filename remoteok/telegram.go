package remoteok

import (
	"gobot/config"
	tele "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"time"
)

/*
SendJob

Send remoteOk job url to config.TelegramRemoteOk channel
*/
func SendJob(job Job) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := tele.NewBot(pref)

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.TelegramRemoteOk), 10, 64)
	if chatIdErr != nil {
		println("Error from send job chatIdErr", chatIdErr)
		return
	}
	_, err := b.Send(tele.ChatID(num), job.URL)
	if err != nil {
		println("Error from send job", err)
		return
	}
}
