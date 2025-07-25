package main

import (
	"gobot/antam"
	"gobot/freegames"
	"gobot/pkg"
	"gobot/reddit"
	"gobot/remoteok"
	"gobot/spotifytube"
	"gobot/warta"
)

func main() {
	pkg.LogWithTimestamp("Gobot v1.9.19 started...")
	freegames.Scouting(false)
	freegames.Cleaning(false)

	remoteok.Scouting()
	remoteok.Cleaning(false)

	warta.Scouting(false)

	antam.Scouting(false)

	reddit.Scouting(false)

	go antam.Run()
	go spotifytube.Run()

	pkg.StartBlocking()
}
