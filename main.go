package main

import (
	"github.com/dohnj0e/snagtag/cmd"
	"github.com/dohnj0e/snagtag/logger"
	"github.com/dohnj0e/snagtag/platforms/tiktok"
	"github.com/dohnj0e/snagtag/platforms/youtube"
)

func main() {
	logger.Init()
	youtube.Init()
	tiktok.Init()

	cmd.Execute()
}
