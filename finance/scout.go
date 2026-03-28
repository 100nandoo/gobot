package finance

import (
	"fmt"
	"gobot/pkg"
)

func Scouting(now bool) {
	execute := func() {
		pkg.LogWithTimestamp("Scouting finance analysis started")
		etf := ETFs[0] // VWRA.L
		result, err := AnalyzeETF(etf.Ticker)
		if err != nil {
			pkg.LogWithTimestamp("%s", fmt.Sprintf("Error analyzing %s: %v", etf.Ticker, err))
			return
		}
		SendAnalysis(etf, result)
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
