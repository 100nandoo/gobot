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

// EverydayOnWeekdaysAt schedules a job to run at a specific hour on weekdays (Monday to Friday)
func EverydayOnWeekdaysAt(operation func(), hour string) {
	// Parse the hour to ensure it’s valid
	if _, err := time.Parse("15:04", hour); err != nil {
		fmt.Println("Error parsing hour:", err)
		return
	}

	// Schedule the job for the specified hour
	_, errJob := scheduler.Every(1).Day().At(hour).Do(func() {
		// Check if today is a weekday before executing the operation
		if isWeekday() {
			operation()
		} else {
			fmt.Println("Job skipped; today is not a weekday.")
		}
	})

	if errJob != nil {
		fmt.Println("Error scheduling gocron job:", errJob)
		return
	}
}

// isWeekday checks if the current day is a weekday
func isWeekday() bool {
	weekday := time.Now().Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
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
