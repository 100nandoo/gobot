package spotifytube

import (
	"log"
	"os"
	"regexp"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

const (
	regexSpotifyUrl = `^(https?\:\/\/)?((www|open|play)\.spotify\.com)\/.+$`
)

func Run() {
	spotifyClient := initClient()

	startMessage := `Hello this is SpotifyTubeBot
*Features:*
- Convert youtube music url to spotify url and vice versa
- Search song using title/artist

*2 ways to convert url:*
I. Send url to the bot
II. Type @spotifytubebot follow by url in chat box on any conversation

*Example:*
@spotifytubebot https://music.youtube.com/watch?v=ezVbN7e-L7Y


SpotifyTubeBot made with ❤️ by @crossix`
	info := "SpotifyTube made with ❤️ by @crossix"

	pref := tele.Settings{
		Token:  os.Getenv("SPOTIFYTUBE_BOT"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(startMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send(startMessage, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})
	})

	b.Handle("/info", func(c tele.Context) error {
		return c.Send(info)
	})

	b.Handle(telebot.OnText, func(c tele.Context) error {
		pattern := regexp.MustCompile(regexSpotifyUrl)

		if pattern.MatchString(c.Text()) {
			song := urlToSong(c.Text(), spotifyClient)
			res := song.Name + " - " + song.Artists[0].Name
			return c.Send(res)
		} else {
			return c.Send("Other regex")
		}
	})

	b.Start()
}
