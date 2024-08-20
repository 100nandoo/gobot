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
	fmt.Println("Gobot v1.8 started...")
	app.Scouting()
	app.Cleaning()

	remoteok.Scouting()
	remoteok.Cleaning()

	go summarizer.Run()

	warta.Scouting()
	pkg.StartBlocking()
}
