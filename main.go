package main

import (
	"gobot/antam"
	"gobot/freegames"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/summarizer"
	"gobot/warta"
)

func main() {
	pkg.LogWithTimestamp("Gobot v1.9.4 started...")
	freegames.Scouting()
	freegames.Cleaning(false)

	remoteok.Scouting()
	remoteok.Cleaning(false)

	warta.Scouting(false)

	antam.Scouting(false)

	go summarizer.Run()
	go antam.Run()

	pkg.StartBlocking()
}
