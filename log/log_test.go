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

package log_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/propervillain/moistpetal/log"
)

// Example demonstrates basic usage of the log package.
func Example() {

	// initialise new global logger
	log.SetDefault(log.New().
		WithLevel(log.DebugLevel).
		WithOutput(os.Stdout).
		WithFlags(0),
	)

	// log messages with different levels
	log.Debug("Use debug to report conditions to developers.")
	log.Info("Use info to report conditions to users.")

	// log messages with different fields
	log.WithScope("pkg").Info("WithScope is used to categorise messages.")
	log.WithError(errors.New("test")).
		Info("WithError is a helper method to log error conditions.")
	log.WithContext(context.WithValue(context.Background(), 0, "testID")).
		Info("WithContext is a helper method to log context's.")
	log.WithFields(log.Fields{
		"num": 97,
		"str": "more years Morty!",
	}).Info("WithFields can be used to add relevant key value pairs to messages.")

	// Output:
	// level=debug msg="Use debug to report conditions to developers."
	// level=info msg="Use info to report conditions to users."
	// level=info msg="WithScope is used to categorise messages." scope=pkg
	// level=info msg="WithError is a helper method to log error conditions." error=test
	// level=info msg="WithContext is a helper method to log context's." ctx=testID
	// level=info msg="WithFields can be used to add relevant key value pairs to messages." num=97 str="more years Morty!"
}

// ExampleNewFile demonstrates how to use a regular file as an external Writer.
func ExampleNewFile() {

	// setup external logger file service
	f, err := log.NewFile("testdata/test.log")
	if err != nil {
		panic(err)
	}
	log.SetDefault(log.New().
		WithFlags(0).
		WithOutput(os.Stdout).
		WithExternal(f),
	)
	defer log.Flush()

	// perform logging as normal
	log.Info("Log to output normally, as well as external file in JSON.")

	// Output:
	// level=info msg="Log to output normally, as well as external file in JSON."
}

// TestDiscard verifies that log messages are ignored at DiscardLevel.
func TestDiscard(t *testing.T) {
	m := "captain's log"
	tlog, verify := newMockLogger(m)
	defer verify(t)
	tlog.Info(m)
	tlog.WithLevel(log.DiscardLevel)
	tlog.Info("discard this msg")
}

// TestExternal tests that an external writer receives events correctly.
func TestExternal(t *testing.T) {
	m := "captain's log"
	e, extVerify := newMockWriter(m)
	tlog, logVerify := newMockLogger(m)
	tlog.WithExternal(e)
	defer logVerify(t)
	defer extVerify(t)
	tlog.Info(m)
	tlog.WithLevel(log.DiscardLevel)
	tlog.Info("discard this msg")
}

// TestLevel verifies that log messages are filtered correctly by level.
func TestLevel(t *testing.T) {
	m := []string{
		"foo\n",
		"bar\n",
	}
	var e string
	for _, msg := range m {
		e += msg
	}
	tlog, verify := newMockLogger(e)
	defer verify(t)
	tlog.Debug("discard this msg")
	tlog.WithLevel(log.DebugLevel)
	tlog.Debug(m[0])
	tlog.Info(m[1])
	tlog.WithLevel(log.DiscardLevel)
	tlog.Info("discard this msg")
}

// TestFile verifies that the log file Writer works as intended.
func TestFile(t *testing.T) {
	filename := filepath.Join("testdata", log.StdFilename())
	f, err := log.NewFile(filename)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	tlog := log.New().
		WithExternal(f).
		WithOutput(buf).
		WithFlags(log.Ljson)
	tlog.Debug("ignored debug message")
	tlog.WithLevel(log.DebugLevel)
	tlog.Debug("debug message")
	tlog.Info("info message")
	tlog.WithLevel(log.DiscardLevel)
	tlog.Info("discard message")
	tlog.Flush()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(data)
	expected := buf.String()
	if actual != expected {
		t.Errorf("\ngot:\n  %q\nwant:\n  %q\n", actual, expected)
	}
}

func newMockLogger(expected string) (log.State, func(*testing.T)) {
	w, fn := newMockWriter(expected)
	l := log.New().
		WithOutput(nil).
		WithFlags(0).
		WithInternal(w)
	return l, fn
}

type mockWriter struct {
	w io.Writer
}

func newMockWriter(expected string) (log.Writer, func(*testing.T)) {
	out := new(bytes.Buffer)
	fn := func(t *testing.T) {
		actual := strings.TrimSpace(out.String())
		expected = strings.TrimSpace(expected)
		if actual != expected {
			t.Errorf("\ngot:\n  %q\nwant:\n  %q\n", actual, expected)
		}
	}
	return &mockWriter{out}, fn
}

func (mw *mockWriter) SetFlags(f log.Flag) {
	// nop
}

func (mw *mockWriter) SetOutput(w io.Writer) {
	mw.w = w
}

func (mw *mockWriter) Log(level log.Level, fields log.Fields, msg string) {
	mw.w.Write([]byte(msg))
}

func (mw *mockWriter) Flush() {
	// nop
}
