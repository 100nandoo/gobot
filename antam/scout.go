package antam

import (
	"fmt"
	"gobot/pkg"
)

/*
Scouting
Run every day at specified times.
*/
func Scouting(now bool) {
	// Helper function to log and send price
	sendGoldPriceAt := func(timeStr string, immediate bool) {
		execute := func() {
			pkg.LogWithTimestamp(fmt.Sprintf("Scouting Antam at %s", timeStr))
			price, err := getGoldPricesFromHTML()
			if err != nil {
				pkg.LogWithTimestamp(fmt.Sprintf("Error fetching gold prices at %s: %v", timeStr, err))
				return
			}
			SendPrice(*price)
		}

		if immediate {
			execute()
		} else {
			pkg.EverydayAtThisHour(execute, timeStr)
		}
	}

	if now {
		pkg.LogWithTimestamp("Running scouting logic immediately")
		sendGoldPriceAt("now", true)
	} else {
		pkg.LogWithTimestamp("Scheduling scouting antam price")
		sendGoldPriceAt("10:05", false)
		sendGoldPriceAt("15:05", false)
		sendGoldPriceAt("16:30", false)
	}
}