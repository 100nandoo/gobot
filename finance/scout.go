package finance

import (
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"os"
	"strconv"
)

func Scouting(now bool) {
	execute := func() {
		pkg.LogWithTimestamp("Scouting finance analysis started")
		channelID, err := strconv.ParseInt(os.Getenv(config.ChannelFinance), 10, 64)
		if err != nil {
			pkg.LogWithTimestamp("Error parsing finance channel ID: %v", err)
			return
		}

		settings := GetFinanceSettings(channelID)
		for _, etf := range settings.Symbols {
			result, analyzeErr := AnalyzeETF(etf.Ticker)
			if analyzeErr != nil {
				pkg.LogWithTimestamp("%s", fmt.Sprintf("Error analyzing %s: %v", etf.Ticker, analyzeErr))
				continue
			}

			previousScore, hasPreviousScore := settings.LastScores[etf.Ticker]

			if result.CompositeScore < settings.AlertThreshold {
				pkg.LogWithTimestamp("Skipping %s because score %+d is below alert threshold %+d", etf.Ticker, result.CompositeScore, settings.AlertThreshold)
				if err := UpdateLastScore(channelID, etf.Ticker, result.CompositeScore); err != nil {
					pkg.LogWithTimestamp("Error updating last score for %s: %v", etf.Ticker, err)
				}
				continue
			}

			if hasPreviousScore && result.CompositeScore <= previousScore {
				pkg.LogWithTimestamp("Skipping %s because score %+d did not improve from previous %+d", etf.Ticker, result.CompositeScore, previousScore)
				if err := UpdateLastScore(channelID, etf.Ticker, result.CompositeScore); err != nil {
					pkg.LogWithTimestamp("Error updating last score for %s: %v", etf.Ticker, err)
				}
				continue
			}

			SendAnalysis(etf, result)
			if err := UpdateLastScore(channelID, etf.Ticker, result.CompositeScore); err != nil {
				pkg.LogWithTimestamp("Error updating last score for %s: %v", etf.Ticker, err)
			}
		}
		pkg.LogWithTimestamp("Scouting finance analysis finished")
	}

	if now {
		pkg.LogWithTimestamp("Running finance scouting immediately")
		execute()
	} else {
		pkg.LogWithTimestamp("Scheduling finance scouting")
		// Market open: LSE opens 8:00 AM GMT = 4:00 PM SGT
		pkg.EverydayOnWeekdaysAt(execute, 16, 0)
	}
}
