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

// Package path exports filepath primitives for the moistpetal framework.
package path

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// Global default path base.
var defaultBase = NewBase(Locate("", ""))

var (
	// HomeSymbol represents user home directory.
	HomeSymbol = "~"
)

// Base represents a base directory for a path tree.
type Base struct {
	dir string
}

// NewBase returns a new default path base.
func NewBase(dir string) *Base {
	return &Base{
		dir: dir,
	}
}

// Dir returns the filepath of the base dir.
func (b *Base) Dir() string {
	return b.dir
}

// SetDefault overrides the default base path.
func SetDefault(base *Base) {
	defaultBase = base
}

// Dir for the default path base.
func Dir() string {
	return defaultBase.Dir()
}

/*
Locate selects a base directory in following order:
 1. Directory from Environment Variable
 2. Provided default directory
 3. Current working directory
*/
func Locate(pathEnv string, pathDefault string) string {
	path := func() string {

		// is the path environment variable present?
		if p := os.Getenv(pathEnv); p != "" {
			return p
		}

		// is default path set?
		if pathDefault != "" {
			return pathDefault
		}

		// use working directory as last resort
		p, _ := os.Getwd()
		return p
	}()
	return Clean(path)
}

// Expand will expand environment variables and home symbol.
func Expand(path string) string {
	if strings.HasPrefix(path, HomeSymbol) {
		if home := HomeDir(); home != "" {
			path = filepath.Join(home, path[len(HomeSymbol):])
		}
	}
	return os.ExpandEnv(path)
}

// Clean canonicalize's a path by expanding home symbols and
// environment variables, evaluating symlinks and converting
// to an absolute path.
func Clean(path string) string {
	path = Expand(path)
	if p, err := filepath.EvalSymlinks(path); err == nil {
		path = p
	}
	if p, err := filepath.Abs(path); err == nil {
		path = p
	}
	return path
}

// Exists returns whether the file/directory exists or not.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// HomeDir returns the current user's home directory.
func HomeDir() string {

	// create validation lambda function
	valid := func(dir string) error {

		fi, err := os.Stat(dir)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			return errors.New("not a directory")
		}

		return nil
	}

	// query os for current user
	if u, _ := user.Current(); u != nil && valid(u.HomeDir) == nil {
		return u.HomeDir
	}

	// macOS, linux
	if s := os.Getenv("HOME"); valid(s) == nil {
		return s
	}

	// windows (preferred)
	if s := os.Getenv("USERPROFILE"); valid(s) == nil {
		return s
	}

	// windows
	if s := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH"); valid(s) == nil {
		return s
	}

	panic("unable to locate user home directory")
}

// ProgramDir returns the program directory for the current OS.
func ProgramDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.ExpandEnv("%PROGRAMDATA%")
	default:
		return "/usr/local"
	}
}
