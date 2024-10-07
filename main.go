package main

import (
	"gobot/freegames"
	"gobot/pkg"
	"gobot/remoteok"
	"gobot/summarizer"
	"gobot/warta"
)

func main() {
	pkg.LogWithTimestamp("Gobot v1.9.3 started...")
	freegames.Scouting()
	freegames.Cleaning(false)

	remoteok.Scouting()
	remoteok.Cleaning(false)

	warta.Scouting(false)

	go summarizer.Run()
	pkg.StartBlocking()
}
