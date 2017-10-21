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

package path_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/propervillain/moistpetal/path"
)

// Example demonstrates basic usage of the path package.
func Example() {

	// configure path vars
	pathEnv := "MOISTPETAL_PATH"
	pathDefault := filepath.Join(path.ProgramDir(), "moistpetal")

	// locate path base directory
	path.SetDefault(path.NewBase(path.Locate(pathEnv, pathDefault)))

	// open file in program directory
	filename := filepath.Join(path.Dir(), "moistpetal.yml")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// do something with data
	_ = data
}

// TestExpand tests expansion of home symbol and env vars.
func TestExpand(t *testing.T) {

	home := path.Expand(path.HomeSymbol)
	if home == path.HomeSymbol || home == "" {
		t.Errorf("failed to expand home symbol")
	}

	v := "wubba lubba dub dub!"
	os.Setenv("TEST_VAR", v)
	check(t, path.Expand("$TEST_VAR"), v, "expand environment var")
}

// TestLocate tests the various methods of finding the base path.
func TestLocate(t *testing.T) {

	pathEnv := "TEST_PATH"
	pathDir := "testdata"
	pathDirClean := path.Clean(pathDir)

	os.Unsetenv(pathEnv)
	base := path.NewBase(path.Locate(pathEnv, pathDir))
	check(t, base.Dir(), pathDirClean, "set base with default dir")

	os.Setenv(pathEnv, pathDir)
	base = path.NewBase(path.Locate(pathEnv, pathDir))
	check(t, base.Dir(), pathDirClean, "set base with default dir (empty env)")

	os.Setenv(pathEnv, "testdata")
	base = path.NewBase(path.Locate(pathEnv, ""))
	check(t, base.Dir(), pathDirClean, "set base with environment var")
}

// TestClean tests path canonicalization.
func TestClean(t *testing.T) {
	var p string

	p = path.Clean("~/")
	if p == "" || p != path.HomeDir() {
		t.Errorf("failed to clean home symbol path")
	}

	p = path.Clean("testdata")
	cwd, _ := os.Getwd()
	if p == "" || p != filepath.Join(cwd, "testdata") {
		t.Errorf("failed to clean relative path")
	}

	p = path.Clean(path.ProgramDir())
	if p == "" || p != path.ProgramDir() {
		t.Errorf("failed to clean absolute path")
	}

	// add symbolic link test
}

// TestExists tests checking for existing files/directories.
func TestExists(t *testing.T) {
	if !path.Exists("path.go") {
		t.Errorf("path.go exists but path.Exists() returned false")
	}
	if path.Exists("invalid.file") {
		t.Errorf("invalid.file does not exist but path.Exists() returned true")
	}
	if !path.Exists("../path") {
		t.Errorf("../path exists but path.Exists() returned false")
	}
	if path.Exists("../invalid") {
		t.Errorf("../invalid does not exist but path.Exists() returned true")
	}
}

func check(t *testing.T, actual string, expected string, desc string) {
	if actual != expected {
		t.Errorf("[%s]\ngot:\n  %q\nwant:\n  %q", desc, actual, expected)
	}
}
