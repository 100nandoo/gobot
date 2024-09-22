package pkg

import (
	"fmt"
	"time"
)

// Define the time format as a constant
const TimeFormat = "2006-01-02 15:04:05"

// LogWithTimestamp prints a log message with a timestamp
func LogWithTimestamp(message string, args ...interface{}) {
	fmt.Println(time.Now().Format(TimeFormat), fmt.Sprintf(message, args...))
}