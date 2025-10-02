package main

import (
	"gobot/antam"
	"gobot/freegames"
	"gobot/pkg"
	"gobot/reddit"
	"gobot/spotifytube"
)

func main() {
	pkg.LogWithTimestamp("Gobot v1.9.20 started...")
	freegames.Scouting(false)
	freegames.Cleaning(false)

	antam.Scouting(false)

	reddit.Scouting(false)

	go antam.Run()
	go spotifytube.Run()

	pkg.StartBlocking()
}
