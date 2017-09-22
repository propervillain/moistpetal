package main

import (
	"fmt"

	"github.com/propervillain/moistpetal/subcmd"
	"github.com/ttacon/chalk"
)

func main() {

	// display banner
	fmt.Printf(subcmd.Banner() + "\n")
	b := chalk.Magenta.NewStyle().WithTextStyle(chalk.Bold)
	fmt.Printf(b.Style(" + -- --={ mpd }\n\n"))
}
