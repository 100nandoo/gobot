package reddit

import (
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Scouting(now bool) {
	// Wrapper function to create standardized scouting logic
	createScoutingLogic := func(subreddit string, score int) func() {
		return func() {
			pkg.LogWithTimestamp("Scouting reddit %s started", subreddit)
			RedditTopPosts(subreddit, Week, score)
			pkg.LogWithTimestamp("Scouting reddit %s finished successfully", subreddit)
		}
	}

	// Create scouting logic for each subreddit
	pixelographyLogic := createScoutingLogic("pixelography", 100)
	mobileLogic := createScoutingLogic("mobilephotography", 180)
	itapLogic := createScoutingLogic("itookapicture", 1000)
	postprocessingLogic := createScoutingLogic("postprocessing", 450)

	if now {
		pkg.LogWithTimestamp("Running all scouting logic immediately")
		pixelographyLogic()
		mobileLogic()
		itapLogic()
		postprocessingLogic()
	} else {
		pkg.LogWithTimestamp("Scheduling scouting logic across different days")
		pkg.SpecificDayAtThisHour(pixelographyLogic, time.Monday, "10:10")
		pkg.SpecificDayAtThisHour(mobileLogic, time.Tuesday, "10:10")
		pkg.SpecificDayAtThisHour(itapLogic, time.Wednesday, "10:10")
		pkg.SpecificDayAtThisHour(postprocessingLogic, time.Thursday, "10:10")
	}
}

func RedditTopPosts(subreddit string, timeFilter TimeFilter, score int) error {
	pref := tele.Settings{
		Token:  os.Getenv(config.TelegramBot),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	num, chatIdErr := strconv.ParseInt(os.Getenv(config.ChannelReddit), 10, 64)
	if chatIdErr != nil {
		return fmt.Errorf("failed to parse chat ID: %w", chatIdErr)
	}

	response, err := FetchTopPosts(subreddit, timeFilter, score)
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	if len(response.Data.Children) == 0 {
		return fmt.Errorf("no posts found")
	}

	var sendErrors []error
	for _, child := range response.Data.Children {
		post := child.Data

		// Create markdown message (caption)
		caption := fmt.Sprintf(
			"*%s*\n"+"by %s\n"+
				"[View on Reddit](https://old.reddit.com%s)",
			escapeMarkdown(post.Title),
			escapeMarkdown(post.Author),
			post.Permalink,
		)

		if !post.IsGallery {
			// Send the image with the caption
			_, err = b.Send(tele.ChatID(num), &tele.Photo{
				File:    tele.FromURL(post.Image),
				Caption: caption,
			}, &tele.SendOptions{
				ParseMode: tele.ModeMarkdown,
			})
			if err != nil {
				log.Printf("failed to send single image: %v", err)
				sendErrors = append(sendErrors, fmt.Errorf("failed to send image with caption: %w", err))
				continue // Skip to next post instead of returning
			}
		} else {
			const maxAlbumSize = 10
			const maxRetries = 3
			const baseDelay = 6 * time.Second

			// Split images into chunks of maxAlbumSize
			for i := 0; i < len(post.Images); i += maxAlbumSize {
				end := i + maxAlbumSize
				if end > len(post.Images) {
					end = len(post.Images)
				}

				var album tele.Album
				for _, imageURL := range post.Images[i:end] {
					if imageURL == "" {
						log.Printf("empty image URL found in gallery")
						continue
					}
					album = append(album, &tele.Photo{
						File: tele.FromURL(imageURL),
					})
				}

				if len(album) == 0 {
					log.Printf("no valid images in this album chunk")
					continue
				}

				// Add caption to first image in each album chunk
				album[0].(*tele.Photo).Caption = caption

				// Implement retry logic with exponential backoff
				var lastErr error
				for retry := 0; retry < maxRetries; retry++ {
					_, err = b.SendAlbum(tele.ChatID(num), album, &tele.SendOptions{
						ParseMode: tele.ModeMarkdown,
					})
					if err == nil {
						break
					}

					lastErr = err
					if strings.Contains(err.Error(), "retry after") {
						// Extract retry after duration if available
						retryAfter := (1 << retry) * int(baseDelay.Seconds())
						log.Printf("Rate limited. Waiting %d seconds before retry %d/%d", retryAfter, retry+1, maxRetries)
						time.Sleep(time.Duration(retryAfter) * time.Second)
						continue
					}
					// If it's not a rate limit error, break the retry loop
					break
				}

				if lastErr != nil {
					log.Printf("failed to send gallery images after %d retries: %v", maxRetries, lastErr)
					sendErrors = append(sendErrors, fmt.Errorf("failed to send gallery images: %w", lastErr))
					break
				}

				// Add delay between sending chunks
				if i+maxAlbumSize < len(post.Images) {
					time.Sleep(2 * time.Second) // Increased base delay between chunks
				}
			}
		}

		time.Sleep(1 * time.Second) // Increased delay between posts
	}

	// Return combined errors if any occurred
	if len(sendErrors) > 0 {
		return fmt.Errorf("encountered %d errors while sending posts: %v", len(sendErrors), sendErrors)
	}

	return nil
}

func escapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "`", "["}
	escaped := text
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}
