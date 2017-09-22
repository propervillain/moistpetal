package main

import (
	"fmt"

	"github.com/propervillain/moistpetal/subcmd"
	"github.com/ttacon/chalk"
)

func main() {

	// display banner
	fmt.Printf(subcmd.Banner() + "\n")
	b := chalk.White.NewStyle().WithTextStyle(chalk.Bold)
	fmt.Printf(b.Style(" + -- --={ mpconsole }\n\n"))

	// display farewell
	fmt.Printf(subcmd.Farewell() + "\n")
}
