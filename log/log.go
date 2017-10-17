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

/*
Package log exports logging primitives for the moistpetal framework. By default it will log to os.Stderr using ANSI colored text. Writers can be specified to support custom log formatting.

This package is inspired from the following sources:

	- https://dave.cheney.net/2015/11/05/lets-talk-about-logging
	- github.com/sirupsen/logrus
	- upspin.io/log
*/
package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Global default logger.
var defaultLogger = New()

// Flag represents a logging option.
type Flag uint32

// Bits or'ed together to control how log message format.
const (
	Ltimestamp Flag         = 1 << iota // include log timestamp
	Ljson                               // log messages in json format
	LstdFlags  = Ltimestamp             // initial values for the default logger
)

// Level represents the level of logging.
type Level int

// Structured levels of logging.
const (
	DebugLevel Level = iota
	InfoLevel
	DiscardLevel
)

// Different logging field keys.
const (
	scopeKey   = "scope"
	errorKey   = "error"
	contextKey = "ctx"
)

// Fields represents structured logging context.
type Fields map[string]interface{}

// Logger describes the high-level interface for logging messages.
type Logger interface {

	// WithScope sets the scope key value pair.
	WithScope(string) Logger

	// WithError sets the error key value pair.
	WithError(error) Logger

	// WithContext sets the context key value pair.
	// If context.Value(0) returns a string value then this
	// will be shown for logging, otherwise the context address is used.
	WithContext(context.Context) Logger

	// WithFields sets a collection of key value pairs.
	WithFields(Fields) Logger

	// Debug writes a message to the log, at Debug level.
	Debug(msg string)

	// Info writes a message to the log, at Info level.
	Info(msg string)
}

type logger struct {
	state  *state
	fields Fields
}

func newLogger(state *state) *logger {
	return &logger{
		state:  state,
		fields: make(Fields),
	}
}

func (l *logger) clone() *logger {
	c := newLogger(l.state)
	for k, v := range l.fields {
		c.fields[k] = v
	}
	return c
}

// WithScope implements Logger.
func (l *logger) WithScope(scope string) Logger {
	c := l.clone()
	c.fields[scopeKey] = scope
	return c
}

// WithError implements Logger.
func (l *logger) WithError(err error) Logger {
	c := l.clone()
	c.fields[errorKey] = err
	return c
}

// WithContext implements Logger.
func (l *logger) WithContext(ctx context.Context) Logger {

	// ignore nil context
	if ctx == nil {
		return l
	}
	c := l.clone()

	// consider key '0', with string value to be ID by convention
	if v, ok := ctx.Value(0).(string); ok {
		c.fields[contextKey] = fmt.Sprintf("%v", v)
	} else { // otherwise we can only differentiate by address
		c.fields[contextKey] = fmt.Sprintf("%p", ctx)
	}

	return c
}

// WithFields implements Logger.
func (l *logger) WithFields(fields Fields) Logger {
	c := l.clone()
	for k, v := range fields {
		c.fields[k] = v
	}
	return c
}

// Info implements Logger.
func (l *logger) Info(msg string) {
	if InfoLevel < l.state.level {
		return
	}
	l.state.log(InfoLevel, l.fields, msg)
}

// Debug implements Logger.
func (l *logger) Debug(msg string) {
	if DebugLevel < l.state.level {
		return
	}
	l.state.log(DebugLevel, l.fields, msg)
}

// Writer describes a service that writes formatted log messages.
type Writer interface {

	// SetFlags configures the writer for log format options.
	SetFlags(Flag)

	// SetOutput sets the underlying write to write to.
	SetOutput(io.Writer)

	// Log processes a log event and writes a formatted log message.
	Log(Level, Fields, string)

	// Flush finalises any pending log messages.
	Flush()
}

// State describes the underlying configuration for high-level loggers.
type State interface {

	// WithFlags sets the log format options.
	WithFlags(Flag) State

	// WithLevel sets the logging level.
	WithLevel(Level) State

	// WithOutput sets the underlying writer to output to.
	WithOutput(io.Writer) State

	// WithInternal sets the internal log formatter.
	WithInternal(Writer) State

	// WithExternal sets the external log formatter.
	WithExternal(Writer) State

	// Flush finalises any pending log messages.
	Flush()

	// Base logger interface for the configuration state.
	Logger
}

type state struct {
	flags    Flag
	level    Level
	output   io.Writer
	internal Writer
	external Writer
	*logger
}

func (s *state) log(level Level, fields Fields, msg string) {
	s.internal.Log(level, fields, msg)
	if s.external != nil {
		s.external.Log(level, fields, msg)
	}
}

// WithFlags configures logging options.
func (s *state) WithFlags(f Flag) State {
	s.flags = f
	s.internal.SetFlags(f)
	if s.external != nil {
		s.external.SetFlags(f)
	}
	return s
}

// WithLevel enables logging at the specified level and above.
func (s *state) WithLevel(level Level) State {
	s.level = level
	return s
}

// WithOutput sets the logger to write to w.
// If w is nil, then logging is disabled.
func (s *state) WithOutput(w io.Writer) State {
	s.output = w
	s.internal.SetOutput(w)
	return s
}

// WithInternal sets the internal writer for log messages.
func (s *state) WithInternal(w Writer) State {
	s.internal = w
	return s
}

// WithExternal sets an external service to process log messages.
func (s *state) WithExternal(w Writer) State {
	s.external = w
	return s
}

// Flush flushes the logging writers.
func (s *state) Flush() {
	s.internal.Flush()
	if s.external != nil {
		s.external.Flush()
	}
}

// New returns a new default log state.
func New() State {
	flags := LstdFlags
	s := &state{
		flags:    flags,
		level:    InfoLevel,
		internal: newLogrus(os.Stderr, flags),
	}
	s.logger = newLogger(s)
	return s
}

// SetDefault overrides the default logger.
func SetDefault(s State) {
	defaultLogger = s
}

// Flush for default global logger.
func Flush() {
	defaultLogger.Flush()
}

// WithScope for default global logger.
func WithScope(scope string) Logger {
	return defaultLogger.WithScope(scope)
}

// WithError for default global logger.
func WithError(err error) Logger {
	return defaultLogger.WithError(err)
}

// WithContext for default global logger.
func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}

// WithFields for default global logger.
func WithFields(fields Fields) Logger {
	return defaultLogger.WithFields(fields)
}

// Debug for default global logger.
func Debug(msg string) {
	defaultLogger.Debug(msg)
}

// Info for default global logger.
func Info(msg string) {
	defaultLogger.Info(msg)
}
