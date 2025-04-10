package reddit

import (
	"fmt"
	"gobot/pkg"
	"log"
	"time"
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
		pkg.SpecificDayAtThisHour(pixelographyLogic, time.Monday, 10, 10)
		pkg.SpecificDayAtThisHour(mobileLogic, time.Tuesday, 10, 10)
		pkg.SpecificDayAtThisHour(itapLogic, time.Wednesday, 10, 10)
		pkg.SpecificDayAtThisHour(postprocessingLogic, time.Thursday, 10, 10)
	}
}

func RedditTopPosts(subreddit string, timeFilter TimeFilter, score int) error {
	client, err := NewTelegramClient()
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	response, err := FetchTopPosts(subreddit, timeFilter, score)
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	if len(response.Data.Children) == 0 {
		return fmt.Errorf("no posts found")
	}

	var sendErrors []error
	for i, child := range response.Data.Children {
		if err := client.SendRedditPost(&child.Data, i != 0); err != nil {
			log.Printf("failed to send post: %v", err)
			sendErrors = append(sendErrors, err)
			continue
		}
		time.Sleep(1 * time.Second)
	}

	if len(sendErrors) > 0 {
		return fmt.Errorf("encountered %d errors while sending posts: %v", len(sendErrors), sendErrors)
	}

	return nil
}
