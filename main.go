package main

import (
	"gobot/antam"
	"gobot/finance"
	"gobot/freegames"
	"gobot/pkg"
	"gobot/reddit"
	"gobot/saaham"
	"gobot/spotifytube"
)

func main() {
	pkg.LogWithTimestamp("Gobot 2.0.0 started...")
	freegames.Scouting(false)
	freegames.Cleaning(false)

	antam.Scouting(true)

	reddit.Scouting(false)

	finance.Scouting(false)

	go antam.Run()
	go spotifytube.Run()
	go finance.Run()
	go saaham.Run()

	pkg.StartBlocking()
}
