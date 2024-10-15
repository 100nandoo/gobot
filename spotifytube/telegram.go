package spotifytube

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gobot/config"

	tele "gopkg.in/telebot.v3"
)

// Constants
const (
	helpMessage = `Hello this is SpotifyTubeBot

**Features:**
- Convert YouTube music URL to Spotify URL and vice versa
- Search song using title/artist

**2 ways to convert URL:**
I. Send URL to the bot
II. Type @spotifytubebot followed by URL in chat box on any conversation

Example:
@spotifytubebot https://music.youtube.com/watch?v=ezVbN7e-L7Y

SpotifyTubeBot made with ❤️ by @crossix`

	SPOTIFY_REGEX = `^(https?://)?((www|open)\.spotify\.com)/.+$`
	YOUTUBE_REGEX  = `^(https?\:\/\/)?((www|music)\.youtube\.com|youtu\.be)\/.+$`
)

// Compile regex patterns once
var (
	spotifyPattern = regexp.MustCompile(SPOTIFY_REGEX)
	youtubePattern  = regexp.MustCompile(YOUTUBE_REGEX)
)

// Run initializes the bot and starts listening for messages
func Run() {
	pref := tele.Settings{
		Token:  os.Getenv(config.SpotifytubeBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Define a common help handler
	helpHandler := func(c tele.Context) error {
		return c.Send(helpMessage, &tele.SendOptions{
			ParseMode: tele.ModeMarkdown,
		})
	}

	// Register the help handler for both commands
	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)

	// Handle incoming text messages
	bot.Handle(tele.OnText, handleTextMessage)

	// Start the bot
	bot.Start()
}

// handleTextMessage processes incoming text messages
func handleTextMessage(c tele.Context) error {
	messageText := c.Text()

	// Check for Spotify URL
	if spotifyPattern.MatchString(messageText) {
		return handleSpotifyURL(c, messageText)
	}

	// Check for YouTube URL
	if youtubePattern.MatchString(messageText) {
		return handleYoutubeURL(c, messageText)
	}

	return nil // No action taken for other messages
}

// handleSpotifyURL processes Spotify URLs
func handleSpotifyURL(c tele.Context, messageText string) error {
	resp, err := InspectUrl(messageText)
	if err != nil {
		return sendErrorMessage(c, "Invalid Spotify URL")
	}
	
	searchResp, err := Search(resp.Data.Name, resp.Data.ArtistNames[0], YOUTUBE_ONLY)
	if err != nil {
		return sendErrorMessage(c, "Can't convert")
	}
	return sendTrackURLs(c, false, searchResp)
}

// handleYoutubeURL processes YouTube URLs
func handleYoutubeURL(c tele.Context, messageText string) error {
	resp, err := InspectUrl(messageText)
	if err != nil {
		return sendErrorMessage(c, "Invalid YouTube URL")
	}

	searchResp, err := Search(resp.Data.Name, resp.Data.ArtistNames[0], SPOTIFY_ONLY)
	if err != nil {
		return sendErrorMessage(c, "Can't convert")
	}
	return sendTrackURLs(c, true, searchResp)
}

// sendErrorMessage sends an error message to the user
func sendErrorMessage(c tele.Context, message string) error {
	return c.Send(message)
}

// sendTrackURLs sends the track URLs to the user
func sendTrackURLs(c tele.Context, isSpotify bool, searchResp *SearchResponse) error {
	// Check if there are tracks in the response
	if len(searchResp.Tracks) == 0 {
		return c.Send("No tracks found.")
	}

	// Create a slice to hold the URLs
	var urls []string

	// Loop through each track and collect the URLs
	for _, track := range searchResp.Tracks {
		if isSpotify {
			if track.Data.URL != nil {
				urls = append(urls, *track.Data.URL)
			} else {
				log.Println("Warning: track.Data.URL is nil, skipping this entry.")
			}
		} else {
			urls = append(urls, YOUTUBE_MUSIC_PREFIX+ *&track.Data.ExternalID)
		}
	}

	// Check if any URLs were collected
	if len(urls) == 0 {
		return c.Send("No URLs found in the track data.")
	}

	// Format the URLs for sending
	responseMessage := strings.Join(urls, "\n")

	// Send the response message
	return c.Send(responseMessage)
}