package pkg

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
	_ "time/tzdata"
)

var Location *time.Location
var scheduler *gocron.Scheduler

// Set location to Singapore
func init() {
	var err error
	Location, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		fmt.Println("errJob load location", err)
		return
	}
	scheduler = gocron.NewScheduler(Location)
}

/*
EverydayAtThisHour

Start go cron job that is scheduled every day at specific hour define on the parameter
*/
func EverydayAtThisHour(operation func(), hour string) {
	_, errJob := scheduler.Every(1).Day().At(hour).Do(operation)

	if errJob != nil {
		fmt.Println("Error doing gocron job", errJob)
		return
	}
}

/*
EverySaturdayDayAtThisHour

Start go cron job that is scheduled every Saturday day at specific hour define on the parameter
*/
func EverySaturdayDayAtThisHour(operation func(), hour string) {
	_, errJob := scheduler.Every(1).Saturday().At(hour).Do(operation)

	if errJob != nil {
		fmt.Println("Error doing gocron job", errJob)
		return
	}
}

func StartBlocking() {
	scheduler.StartBlocking()
}
