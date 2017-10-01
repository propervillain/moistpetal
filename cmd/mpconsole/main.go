// Copyright 2017, The Moistpetal Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/propervillain/moistpetal/subcmd"
	"github.com/propervillain/moistpetal/version"
)

func main() {

	// parse command line flags
	flagVersion := flag.Bool("version", false, "display moistpetal version")
	flag.Parse()

	// handle version
	if *flagVersion {
		fmt.Printf("mpconsole %s\n", version.Version())
		return
	}

	// display banner
	fmt.Printf(subcmd.Banner() + "\n")
	color.New(color.Bold).Printf(" + -- --={ %s %s }\n\n",
		filepath.Base(os.Args[0]),
		version.Version(),
	)

	// display farewell
	fmt.Printf(subcmd.Farewell() + "\n")
}
