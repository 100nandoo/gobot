package spotifytube

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"gobot/config"
	"gobot/pkg"
	"gobot/spotifytube/api"

	"google.golang.org/api/youtube/v3"
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
	currentToken *api.SpotifyTokenResponse
	youtubeService *youtube.Service
	showLog = true
)

// Run initializes the bot and starts listening for messages
func Run() {
	service, err := api.InitYoutube()
	youtubeService = service
	if err != nil {
		log.Fatal(err)
		return
	}

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
	id, err := ExtractSpotifyTrackID(messageText)
	if err != nil {
		return sendErrorMessage(c, "Invalid Spotify URL")
	}

	tokenResp, err := getValidSpotifyAccessToken()
	if err != nil {
		return sendErrorMessage(c, "Unable to retrieve Spotify token")
	}

	track, err := api.GetSpotifyTrack(tokenResp.AccessToken, id)
	if err != nil {
		return sendErrorMessage(c, "Unable to retrieve Spotify track")
	}

	if showLog {
		pkg.LogWithTimestamp("Spotify Track: %s - %s", track.Name, track.Artists[0].Name)
	}

	query := fmt.Sprintf("%s %s", track.Name, track.Artists[0].Name)
	searchResp, err := api.SearchYoutube(*youtubeService, query)
	if err != nil {
		return sendErrorMessage(c, "No matching YouTube results found")
	}

	return sendYoutubeURLs(c, searchResp)
}

func handleYoutubeURL(c tele.Context, messageText string) error {
	id, err := ExtractYoutubeVideoID(messageText)
	if err != nil {
		return sendErrorMessage(c, "Invalid Youtube URL")
	}

	resp, err := api.GetVideo(*youtubeService, id)
	if err != nil && len(resp.Items) == 0 {
		return sendErrorMessage(c, "Unable to retrieve Youtube video")
	}

	if showLog {
		pkg.LogWithTimestamp("Youtube: %s - %s", resp.Items[0].Snippet.Title, resp.Items[0].Snippet.ChannelTitle)
	}

	query := fmt.Sprintf("%s %s", resp.Items[0].Snippet.Title, resp.Items[0].Snippet.ChannelTitle)

	tokenResp, err := getValidSpotifyAccessToken()
	if err != nil {
		return sendErrorMessage(c, "Unable to retrieve Spotify token")
	}
	searchResp, err := api.SearchSpotify(tokenResp.AccessToken, query)
	if err != nil {
		return sendErrorMessage(c, "No matching YouTube results found")
	}

	return sendSpotifyURLs(c, searchResp)
}
func getValidSpotifyAccessToken() (*api.SpotifyTokenResponse, error) {
	// If there's no current token or the token is expired, get a new one
	if currentToken == nil || time.Now().After(currentToken.RequestedAt.Add(time.Duration(currentToken.ExpiresIn)*time.Second)) {
		// Fetch a new token
		token, err := api.GetSpotifyAccessToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get new access token: %v", err)
		}
		// Update the global current token
		currentToken = token
	}

	// Return the current (valid) token
	return currentToken, nil
}

func sendErrorMessage(c tele.Context, message string) error {
	return c.Send(message)
}

func sendYoutubeURLs(c tele.Context, searchResp *youtube.SearchListResponse) error {
	if len(searchResp.Items) == 0 {
		return c.Send("No items found.")
	}

	for _, item := range searchResp.Items {
		err := c.Send(YOUTUBE_MUSIC_PREFIX+ *&item.Id.VideoId)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendSpotifyURLs(c tele.Context, searchResp *api.SpotifySearchResponse) error {
	if len(searchResp.Tracks.Items) == 0 {
		return c.Send("No Spotify tracks found.")
	}

	for _, item := range searchResp.Tracks.Items {
		err := c.Send(item.ExternalURLs.Spotify)
		if err != nil {
			return err
		}
	}
	return nil
}
