package main

import (
	"fmt"
	"gobot/freegames"
	"gobot/freegames/app"
)

func main() {
	fmt.Println("Gobot v1.2 started...")
	app.Scouting()
	app.Cleaning()
	freegames.StartBlocking()
}
