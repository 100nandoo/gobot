package main

import (
	"fmt"
	"gobot/freegames/app"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/summarizer"
)

func main() {
	fmt.Println("Gobot v1.6 started...")
	app.Scouting()
	app.Cleaning()

	remoteok.Scouting()
	remoteok.Cleaning()

	go summarizer.Run()
	pkg.StartBlocking()
}
