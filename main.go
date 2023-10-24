package main

import (
	"github.com/dohnj0e/snagtag/cmd"
	"github.com/dohnj0e/snagtag/platforms/tiktok"
	"github.com/dohnj0e/snagtag/platforms/youtube"
)

func main() {
	youtube.Init()
	tiktok.Init()
	cmd.Execute()
}
