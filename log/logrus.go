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
	"io"

	"github.com/sirupsen/logrus"
)

// abstracts differences of Logger and Entry
type logrusLogger interface {
	Info(...interface{})
	Debug(...interface{})
}

type logrusWriter struct {
	lr *logrus.Logger
}

func newLogrus(w io.Writer, flags Flag) *logrusWriter {
	lrw := &logrusWriter{
		lr: &logrus.Logger{
			Out:   w,
			Level: logrus.DebugLevel,
			Hooks: make(logrus.LevelHooks),
		},
	}
	lrw.SetFlags(flags)
	return lrw
}

// SetFlags implements Writer.
func (lrw *logrusWriter) SetFlags(f Flag) {
	showTimestamp := f&Ltimestamp > 0
	if f&Ljson > 0 {
		lrw.lr.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: !showTimestamp,
		}
	} else {
		lrw.lr.Formatter = &logrus.TextFormatter{
			DisableTimestamp: !showTimestamp,
		}
	}
}

// SetOutput implements Writer.
func (lrw *logrusWriter) SetOutput(w io.Writer) {
	lrw.lr.Out = w
}

// Log implements Writer.
func (lrw *logrusWriter) Log(level Level, fields Fields, msg string) {

	// construct logging context
	var lr logrusLogger
	if len(fields) > 0 {
		lr = lrw.lr.WithFields(logrus.Fields(fields))
	} else {
		lr = lrw.lr
	}

	// determine logging function for level
	var logFn func(...interface{})
	switch level {
	case InfoLevel:
		logFn = lr.Info
	default:
		logFn = lr.Debug
	}

	// write log message
	logFn(msg)
}

// Flush implements Writer.
func (lrw *logrusWriter) Flush() {
	// nop
}
