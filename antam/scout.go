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
			var buyPrices []GoldPrice

			priceHargaEmas, errHargaEmas := getHargaEmasComPrices()
			if errHargaEmas != nil {
				pkg.LogWithTimestamp("%s", fmt.Sprintf("Error fetching hargaemas.com buy price at %d:%d: %v", hour, minute, errHargaEmas))
			} else {
				buyPrices = append(buyPrices, *priceHargaEmas)
			}

			pricePluang, errPluang := getPluangGoldPrices()
			if errPluang != nil {
				pkg.LogWithTimestamp("%s", fmt.Sprintf("Error fetching Pluang gold prices at %d:%d: %v", hour, minute, errPluang))
			} else {
				buyPrices = append(buyPrices, *pricePluang)
			}

			if len(buyPrices) == 0 {
				return
			}

			SendBuyPrice(buyPrices...)
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
