package pkg

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/go-co-op/gocron/v2"
)

var Location *time.Location
var scheduler gocron.Scheduler

// Set location to Singapore
func init() {
	var err error
	Location, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		fmt.Println("errJob load location", err)
		return
	}
	scheduler, err = gocron.NewScheduler(gocron.WithLocation(Location))
}

/*
EverydayAtThisHour

Start go cron job that is scheduled every day at specific hour define on the parameter
*/
func EverydayAtThisHour(operation func(), hour, minute uint) {
	atTimes := gocron.NewAtTimes(gocron.NewAtTime(hour, minute, 0))
	jobDefinition := gocron.DailyJob(1, atTimes)
	task := gocron.NewTask(operation)
	_, err := scheduler.NewJob(jobDefinition, task)
	if err != nil {
		fmt.Println("Error scheduling gocron job:", err)
		return
	}
}

// EverydayOnWeekdaysAt schedules a job to run at a specific hour on weekdays (Monday to Friday)
func EverydayOnWeekdaysAt(operation func(), hour, minute uint) {
	atTimes := gocron.NewAtTimes(gocron.NewAtTime(hour, minute, 0))
	jobDefinition := gocron.DailyJob(1, atTimes)
	task := gocron.NewTask(func() {
		// Check if today is a weekday before executing the operation
		if isWeekday() {
			operation()
		} else {
			fmt.Println("Job skipped; today is not a weekday.")
		}
	})
	// Schedule the job for the specified hour
	_, errJob := scheduler.NewJob(jobDefinition, task)

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
SpecificDayAtThisHour

Start go cron job that is scheduled every specific day at specific hour define on the parameter
*/
func SpecificDayAtThisHour(operation func(), day time.Weekday, hour, minute uint) {
	jobDefinition := gocron.WeeklyJob(
		1,
		gocron.NewWeekdays(time.Weekday(day)),
		gocron.NewAtTimes(gocron.NewAtTime(hour, minute, 0)),
	)

	task := gocron.NewTask(operation)

	_, errJob := scheduler.NewJob(jobDefinition, task)

	if errJob != nil {
		fmt.Println("Error doing gocron job", errJob)
		return
	}
}

func StartBlocking() {
	scheduler.Start()
	select {}
}
