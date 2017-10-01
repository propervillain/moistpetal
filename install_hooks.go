//usr/bin/env go run "$0" "$@"; exit

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
