package main

import (
	"gobot/antam"
	"gobot/freegames"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/warta"
)

func main() {
	pkg.LogWithTimestamp("Gobot v1.9.7 started...")
	freegames.Scouting()
	freegames.Cleaning(false)

	remoteok.Scouting()
	remoteok.Cleaning(false)

	warta.Scouting(false)

	antam.Scouting(false)

	go antam.Run()

	pkg.StartBlocking()
}
