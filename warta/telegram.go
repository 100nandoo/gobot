package warta

import (
	"fmt"
	"gobot/config"
	tele "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"time"
)

func (w SupabaseWarta) formatOutput() string {
	return fmt.Sprintf(
		"%s\n\n*Pengkhotbah pagi:*\n%s\n\n*Pengkhotbah sore:*\n%s",
		w.BulletinDate,
		w.Preacher1,
		w.Preacher2,
	)
}

/*
SendWarta

Send warta to config.TelegramWarta channel
*/
func SendWarta(warta SupabaseWarta) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := tele.NewBot(pref)

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.TelegramWarta), 10, 64)
	if chatIdErr != nil {
		println("Error from send warta chatIdErr", chatIdErr)
		return
	}

	_, err := b.Send(tele.ChatID(num), warta.formatOutput(), tele.ModeMarkdown)
	if err != nil {
		println("Error from send warta", err)
		return
	}
}
