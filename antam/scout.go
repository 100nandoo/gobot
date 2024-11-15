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
			pricePluang, errPluang := getPluangGoldPricesFromHTML()
			if errPluang != nil {
				pkg.LogWithTimestamp(fmt.Sprintf("Error fetching Pluang gold prices at %s: %v", timeStr, err))
				return
			}
			goldPrices := []GoldPrice{*price, *pricePluang}

			SendPrice(goldPrices...)
		}

		if immediate {
			execute()
		} else {
			pkg.EverydayOnWeekdaysAt(execute, timeStr)
		}
	}

	if now {
		pkg.LogWithTimestamp("Running scouting logic immediately")
		sendGoldPriceAt("now", true)
	} else {
		pkg.LogWithTimestamp("Scheduling scouting antam price")
		sendGoldPriceAt("10:05", false)
	}
}