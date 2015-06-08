// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Tobias Schottdorf

package log

import (
	"strings"

	"golang.org/x/net/context"
)

// Level Enum for log severity levels.
type Level int

// TODO(Bram): combine this with clog.severity, and rename severity.

// These constants identify the log levels in order of increasing severity.
const (
	INFO Level = iota
	WARNING
	ERROR
	FATAL
)

var levels = [...]string{
	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
}

// String returns the string representation of the log level.
func (level Level) String() string {
	return levels[level]
}

// LevelFromString returns the log level if it can match it.
func LevelFromString(level string) (Level, bool) {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "FATAL":
		return FATAL, true
	case "ERROR":
		return ERROR, true
	case "WARNING":
		return WARNING, true
	case "INFO":
		return INFO, true
	default:
		return INFO, false
	}
}

func init() {
	// TODO(tschottdorf) this should go to our logger. Currently this will log
	// with clog (=glog) format.
	CopyStandardLogTo("INFO")
}

// FatalOnPanic recovers from a panic and exits the process with a
// Fatal log. This is useful for avoiding a panic being caught through
// a CGo exported function or preventing HTTP handlers from recovering
// panics and ignoring them.
func FatalOnPanic() {
	if r := recover(); r != nil {
		Fatalf("unexpected panic: %s", r)
	}
}

// EnableLogFileOutput turns on logging using the specified directory.
// For unittesting only.
func EnableLogFileOutput(dir string) {
	*logDir = dir
	logging.toStderr = false
	logging.alsoToStderr = true
}

// DisableLogFileOutput turns off logging. For unittesting only.
func DisableLogFileOutput() {
	*logDir = ""
	logging.toStderr = true
	logging.alsoToStderr = false
}

// logDepth uses the PrintWith to format the output string and
// formulate the context information into the machine-readable
// dictionary for separate binary-log output.
func logDepth(ctx context.Context, depth int, sev severity, format string, args []interface{}) {
	// TODO(tschottdorf): logging hooks should have their entry point here.
	AddStructured(ctx, sev, depth+1, format, args)
}

// Infoc logs to the WARNING and INFO logs. It extracts values from the context
// using the Field keys specified in this package and logs them along with the
// given message and any additional pairs specified as consecutive elements in
// kvs.
func Infoc(ctx context.Context, format string, args ...interface{}) {
	logDepth(ctx, 1, infoLog, format, args)
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended.
func Info(args ...interface{}) {
	logDepth(nil, 1, infoLog, "", args)
}

// Infof logs to the INFO log. Don't use it; use Info or Infoc instead.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	logDepth(nil, 1, infoLog, format, args)
}

// InfoDepth logs to the INFO log, offsetting the caller's stack frame by
// 'depth'.
func InfoDepth(depth int, args ...interface{}) {
	logDepth(nil, depth+1, infoLog, "", args)
}

// Warningc logs to the WARNING and INFO logs. It extracts values from the
// context using the Field keys specified in this package and logs them along
// with the given message and any additional pairs specified as consecutive
// elements in kvs.
func Warningc(ctx context.Context, format string, args ...interface{}) {
	logDepth(ctx, 1, warningLog, format, args)
}

// Warning logs to the WARNING and INFO logs.
// Warningf logs to the WARNING and INFO logs. Don't use it; use Warning or
// Arguments are handled in the manner of fmt.Print; a newline is appended.
func Warning(args ...interface{}) {
	logDepth(nil, 1, warningLog, "", args)
}

// Warningf logs to the WARNING and INFO logs. Don't use it; use Warning or
// Warningc instead. Arguments are handled in the manner of fmt.Printf; a
// newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	logDepth(nil, 1, warningLog, format, args)
}

// WarningDepth logs to the WARNING and INFO logs, offsetting the caller's
// stack frame by 'depth'.
func WarningDepth(depth int, args ...interface{}) {
	logDepth(nil, depth+1, warningLog, "", args)
}

// Errorc logs to the ERROR, WARNING, and INFO logs. It extracts values from
// Field keys specified in this package and logs them along with the given
// message and any additional pairs specified as consecutive elements in kvs.
func Errorc(ctx context.Context, format string, args ...interface{}) {
	logDepth(ctx, 1, errorLog, format, args)
}

// Error logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended.
func Error(args ...interface{}) {
	logDepth(nil, 1, errorLog, "", args)
}

// Errorf logs to the ERROR, WARNING, and INFO logs. Don't use it; use Error
// Info or Errorc instead. Arguments are handled in the manner of fmt.Printf;
// a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	logDepth(nil, 1, errorLog, format, args)
}

// ErrorDepth logs to the ERROR, WARNING, and INFO logs, offsetting the
// caller's stack frame by 'depth'.
func ErrorDepth(depth int, args ...interface{}) {
	logDepth(nil, depth+1, errorLog, "", args)
}

// Fatalc logs to the INFO, WARNING, ERROR, and FATAL logs, including a stack
// trace of all running goroutines, then calls os.Exit(255). It extracts values
// from the context using the Field keys specified in this package and logs
// them along with the given message and any additional pairs specified as
// consecutive elements in kvs.
func Fatalc(ctx context.Context, format string, args ...interface{}) {
	logDepth(ctx, 1, fatalLog, format, args)
}

// Fatal logs to the INFO, WARNING, ERROR, and FATAL logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Print; a newline is appended.
func Fatal(args ...interface{}) {
	logDepth(nil, 1, fatalLog, "", args)
}

// Fatalf logs to the INFO, WARNING, ERROR, and FATAL logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Printf; a newline is appended.
func Fatalf(format string, args ...interface{}) {
	logDepth(nil, 1, fatalLog, format, args)
}

// FatalDepth logs to the INFO, WARNING, ERROR, and FATAL logs,
// including a stack trace of all running goroutines, then calls os.Exit(255),
// offsetting the caller's stack frame by 'depth'.
func FatalDepth(depth int, args ...interface{}) {
	logDepth(nil, depth+1, fatalLog, "", args)
}

// V returns true if the logging verbosity is set to the specified level or
// higher.
func V(level level) bool {
	return VDepth(level, 1)
}
