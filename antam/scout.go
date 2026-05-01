package antam

import (
	"fmt"
	"gobot/pkg"
	"time"
)

/*
Scouting
Run once a week at the specified time.
*/
func Scouting(now bool) {
	// Helper function to log and send price
	sendGoldPriceAt := func(hour, minute uint, immediate bool) {
		execute := func() {
			// pkg.LogWithTimestamp(fmt.Sprintf("Scouting Antam at %s", timeStr))
			// price, err := getGoldPrices()
			// if err != nil {
			// 	pkg.LogWithTimestamp(fmt.Sprintf("Error fetching gold prices at %s: %v", timeStr, err))
			// 	return
			// }
			pricePluang, errPluang := getPluangGoldPrices()
			if errPluang != nil {
				pkg.LogWithTimestamp("%s", fmt.Sprintf("Error fetching Pluang gold prices at %d:%d: %v", hour, minute, errPluang))
				return
			}
			goldPrices := []GoldPrice{*pricePluang}

			SendPrice(goldPrices...)
		}

		if immediate {
			execute()
		} else {
			pkg.SpecificDayAtThisHour(execute, time.Wednesday, hour, minute)
		}
	}

	if now {
		pkg.LogWithTimestamp("Running scouting logic immediately")
		sendGoldPriceAt(10, 05, true)
	} else {
		pkg.LogWithTimestamp("Scheduling weekly scouting antam price")
		sendGoldPriceAt(10, 05, false)
	}
}
