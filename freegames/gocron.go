package freegames

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

var Location *time.Location

// Set location to Singapore
func init() {
	var err error
	Location, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		fmt.Println("errJob load location", err)
		return
	}
}

/*
Every11am

Start go cron job that is scheduled every day at 11.00 am
*/
func Every11am(operation func()) {
	s := gocron.NewScheduler(Location)

	_, errJob := s.Every(1).Day().At("11:00").Do(operation)

	if errJob != nil {
		fmt.Println("Error doing gocron job", errJob)
		return
	}

	s.StartBlocking()
}
