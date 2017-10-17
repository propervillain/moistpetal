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

package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StdFilename returns an appropriate log filename for the running process.
func StdFilename() string {
	bin := os.Args[0]
	return fmt.Sprintf("%s_%s.log",
		strings.TrimSuffix(filepath.Base(bin), filepath.Ext(bin)),
		time.Now().Format("2006-01-02_15h04m05s"),
	)
}

type file struct {
	f *os.File
	w Writer
}

// NewFile opens a regular file as a log Writer.
func NewFile(filename string) (Writer, error) {

	// create directory, if required
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return nil, err
	}

	// create log file
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	// create logrus json writer
	flags := Ltimestamp | Ljson
	w := newLogrus(f, flags)

	return &file{f, w}, nil
}

// SetFlags implements Writer.
func (f *file) SetFlags(fl Flag) {
	f.w.SetFlags(fl)
}

// SetOutput implements Writer.
func (f *file) SetOutput(io.Writer) {
	// nop
}

// Log implements Writer.
func (f *file) Log(level Level, fields Fields, msg string) {
	f.w.Log(level, fields, msg)
}

// Flush implements Writer.
func (f *file) Flush() {
	if f.f != nil {
		f.f.Sync()
	}
}
