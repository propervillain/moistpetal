package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/propervillain/moistpetal/subcmd"
)

func main() {

	// display banner
	fmt.Printf(subcmd.Banner() + "\n")
	color.New(color.Bold).Printf(" + -- --={ mpd }\n\n")
}
