package main

import (
	"fmt"
	"gobot/freegames/app"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/rss"
)

func main() {
	fmt.Println("Gobot v1.3 started...")
	app.Scouting()
	app.Cleaning()

	remoteok.Scouting()
	remoteok.Cleaning()

	rss.Scouting()

	pkg.StartBlocking()
}
