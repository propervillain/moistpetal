//usr/bin/env go run "$0" "$@"; exit

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

// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// determine source directory
	srcDir, err := filepath.Abs(".githooks")
	if err != nil {
		log.Fatal(err)
	}

	// determine destination directory
	dstDir, err := filepath.Abs(".git/hooks")
	if err != nil {
		log.Fatal(err)
	}

	// fetch list of files
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	// iterate over source directory entries
	for _, file := range files {

		// determine filepaths
		f := file.Name()
		src := filepath.Join(srcDir, f)
		dst := filepath.Join(dstDir, f)

		// ignore if link already exists
		fi, err := os.Stat(dst)
		if err == nil {
			fmt.Printf("%s: found existing\n", f)
			continue
		}

		// ignore entries that are not regular files
		if fi != nil && !fi.Mode().IsRegular() {
			continue
		}

		// install symlink
		fmt.Printf("%s: installed link\n", f)
		if err := os.Symlink(src, dst); err != nil {
			log.Fatal(err)
		}
	}
}
