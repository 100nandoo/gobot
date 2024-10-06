package main

import (
	"fmt"
	"gobot/freegames/app"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/summarizer"
	"gobot/warta"
)

func main() {
	fmt.Println("Gobot v1.9.3 started...")
	app.Scouting()
	app.Cleaning(false)

	remoteok.Scouting()
	remoteok.Cleaning(false)

	warta.Scouting(false)

	go summarizer.Run()
	pkg.StartBlocking()
}
