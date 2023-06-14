package main

import (
	"fmt"
	"gobot/freegames/app"
	"gobot/pkg"
	"gobot/remoteok"
)

func main() {
	fmt.Println("Gobot v1.3 started...")
	app.Scouting()
	app.Cleaning()

	remoteok.Scouting()
	remoteok.Cleaning()

	pkg.StartBlocking()
}
