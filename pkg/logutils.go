package pkg

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Define the time format as a constant
const TimeFormat = "2006-01-02 15:04:05"

// Max package name length for padding
const MaxPackageNameLength = 18 // Adjust as needed

// LogWithTimestamp prints a log message with a timestamp and package name
func LogWithTimestamp(message string, args ...interface{}) {
	packageName := getPackageName()
	paddedPackageName := padPackageName(packageName)
	fmt.Println(time.Now().Format(TimeFormat), paddedPackageName, fmt.Sprintf(message, args...))
}

// LogWithTimestampInt prints a log message with a timestamp and package name using interface{}
func LogWithTimestampInt(a ...any) {
	packageName := getPackageName()
	paddedPackageName := padPackageName(packageName)
	fmt.Println(time.Now().Format(TimeFormat), paddedPackageName, fmt.Sprint(a...))
}

// getPackageName retrieves the package name of the caller
func getPackageName() string {
	pc, _, _, _ := runtime.Caller(2)         // get the caller's PC, line, file, and error
	funcName := runtime.FuncForPC(pc).Name() // get the function name
	parts := strings.Split(funcName, ".")    // split by "."
	if len(parts) > 0 {
		return parts[len(parts)-2] // return the package name
	}
	return ""
}

// padPackageName adds padding to the package name to align log messages
func padPackageName(pkg string) string {
	pkgLen := len(pkg)
	if pkgLen >= MaxPackageNameLength {
		return pkg
	}
	padding := MaxPackageNameLength - pkgLen
	return pkg + strings.Repeat(" ", padding)
}

func main() {
	LogWithTimestamp("This is a test message: %s", "Hello, World!")
	LogWithTimestampInt("This is a test with integers:", 1, 2, 3)
}
