package pkg

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/go-co-op/gocron"
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

func EverySaturdaySundayThreeHour(operation func()) {
    _, errJob := scheduler.Every(3).Hours().Saturday().Sunday().Do(operation)
    
    if errJob != nil {
        fmt.Println("Error scheduling gocron job:", errJob)
        return
    }
    
    // Start the scheduler
    scheduler.StartAsync()
}

func StartBlocking() {
	scheduler.StartBlocking()
}
