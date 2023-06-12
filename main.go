package main

import (
	"fmt"
	"gobot/freegames"
	"gobot/freegames/app"
)

func main() {
	fmt.Println("Gobot v1.1 started...")
	go app.Scouting()
	go app.Cleaning()
	freegames.StartBlocking()
}
