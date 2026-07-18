package reddit

import (
	"bytes"
	"fmt"
	"gobot/config"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type TelegramClient struct {
	bot    *tele.Bot
	chatID int64
}

var telegramMediaHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

func NewTelegramClient() (*TelegramClient, error) {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.ChannelReddit), 10, 64)
	if chatIdErr != nil {
		return nil, fmt.Errorf("failed to parse chat ID: %w", chatIdErr)
	}

	return &TelegramClient{
		bot:    b,
		chatID: num,
	}, nil
}

func (t *TelegramClient) SendRedditPost(post *Post, isSilent bool) error {
	caption := fmt.Sprintf(
		"*%s*\n"+
			"by %s ⬆️ %s\n"+
			"[View on Reddit](https://old.reddit.com%s)",
		escapeMarkdown(post.Title),
		escapeMarkdown(post.Author),
		escapeMarkdown(strconv.Itoa(post.Score)),
		post.Permalink,
	)

	if !post.IsGallery {
		return t.sendSingleImage(post.Image, caption, isSilent)
	}
	return t.sendGallery(post.Images, caption, isSilent)
}

func (t *TelegramClient) sendSingleImage(imageURL, caption string, isSilent bool) error {
	photo, err := photoFromURL(imageURL)
	if err != nil {
		return fmt.Errorf("failed to prepare single image: %w", err)
	}

	photo.Caption = caption

	_, err = t.bot.Send(tele.ChatID(t.chatID), photo, &tele.SendOptions{
		ParseMode:           tele.ModeMarkdown,
		DisableNotification: isSilent,
	})
	if err != nil {
		return fmt.Errorf("failed to send single image: %w", err)
	}
	return nil
}

func (t *TelegramClient) sendGallery(images []string, caption string, isSilent bool) error {
	const maxAlbumSize = 10
	const maxRetries = 3
	const baseDelay = 6 * time.Second

	for i := 0; i < len(images); i += maxAlbumSize {
		end := i + maxAlbumSize
		if end > len(images) {
			end = len(images)
		}

		var album tele.Album
		for _, imageURL := range images[i:end] {
			if imageURL == "" {
				log.Printf("empty image URL found in gallery")
				continue
			}

			photo, err := photoFromURL(imageURL)
			if err != nil {
				log.Printf("failed to prepare gallery image %q: %v", imageURL, err)
				continue
			}

			album = append(album, photo)
		}

		if len(album) == 0 {
			log.Printf("no valid images in this album chunk")
			continue
		}

		album[0].(*tele.Photo).Caption = caption

		if err := t.sendAlbumWithRetry(album, maxRetries, baseDelay, isSilent); err != nil {
			return err
		}

		if i+maxAlbumSize < len(images) {
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func (t *TelegramClient) sendAlbumWithRetry(album tele.Album, maxRetries int, baseDelay time.Duration, isSilent bool) error {
	var lastErr error
	for retry := 0; retry < maxRetries; retry++ {
		_, err := t.bot.SendAlbum(tele.ChatID(t.chatID), album, &tele.SendOptions{
			ParseMode:           tele.ModeMarkdown,
			DisableNotification: isSilent,
		})
		if err == nil {
			return nil
		}

		lastErr = err
		if strings.Contains(err.Error(), "retry after") {
			retryAfter := (1 << retry) * int(baseDelay.Seconds())
			log.Printf("Rate limited. Waiting %d seconds before retry %d/%d", retryAfter, retry+1, maxRetries)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}
		break
	}
	return fmt.Errorf("failed to send gallery images after %d retries: %w", maxRetries, lastErr)
}

func photoFromURL(imageURL string) (*tele.Photo, error) {
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := telegramMediaHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download image: unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read image body: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("read image body: empty response")
	}

	return &tele.Photo{
		File: tele.FromReader(bytes.NewReader(body)),
	}, nil
}
