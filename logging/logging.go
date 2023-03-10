/* MIT License
*
* Copyright (c) 2018 Mike Taghavi <mitghi[at]gmail.com>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
 */

package logging

import (
	"fmt"

  "github.com/romana/rlog"
)

const logformat = "%s <- (*%s).%s"

// Logging is the default logging mechanism
type Logging struct {
	Package string
	debug   bool
}

func formatArgs(args ...interface{}) []interface{} {
	var (
		ns []interface{} = make([]interface{}, 0, len(args)+2)
	)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	return ns
}

// NewLogger returns a pointer to a compatible  `protobase.LoggingInterface`.
func NewLogger(pkg string) *Logging {
	return &Logging{pkg, false}
}

// SetDebug sets debug flag to `status`.
func (l *Logging) SetDebug(status bool) {
	l.debug = status
}

// FInfo logs a entry with caller name prepended.
func (l *Logging) FInfo(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Info(ns...)
}

// FInfo logs a entry with caller name prepended.
func (l *Logging) FInfof(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("<- (*%s).%s", l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Infof(msg, ns...)
}

// FDebug logs a entry with caller name prepended.
func (l *Logging) FDebug(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Debug(ns...)
}

func (l *Logging) FDebugf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("%s <- (*%s).%s", msg, l.Package, fn)
	rlog.Debugf(format, args...)
}

// FWarn logs a entry with caller name prepended.
func (l *Logging) FWarn(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Warn(ns...)
}

// FWarn logs a entry with caller name prepended.
func (l *Logging) FWarnf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Warnf(msg, ns...)
}

// FError logs a entry with caller name prepended.
func (l *Logging) FError(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Error(ns...)
}

// FError logs a entry with caller name prepended.
func (l *Logging) FErrorf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Errorf(msg, ns...)
}

// FFatal logs a entry with caller name prepended.
func (l *Logging) FFatal(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns[0] = format
	copy(ns[1:], args)
	rlog.Critical(ns...)
}

// FFatal logs a entry with caller name prepended.
func (l *Logging) FFatalf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("<- (*%s).%s", l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Criticalf(msg, ns...)
}

// Trace logs a trace entry.
func (l *Logging) FTrace(level int, fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Trace(level, ns...)
}

// Trace logs a trace entry.
func (l *Logging) FTracef(level int, fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, l.Package, fn)
	rlog.Tracef(level, format, args...)
}

// Trace logs a trace entry.
func (l *Logging) Trace(level int, args ...interface{}) {
	rlog.Trace(level, args...)
}

// Trace logs a trace entry.
func (l *Logging) Tracef(level int, msg string, args ...interface{}) {
	rlog.Tracef(level, msg, args...)
}

// Debug logs a debug entry.
func (l *Logging) Debug(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Debug(ns...)
}

// Debug logs a debug entry.
func (l *Logging) Debugf(msg string, args ...interface{}) {
	rlog.Debugf(msg, args...)
}

// Info logs a debug entry.
func (l *Logging) Info(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Info(ns...)
}

// Info logs a debug entry.
func (l *Logging) Infof(msg string, args ...interface{}) {
	rlog.Infof(msg, args...)
}

// Warn logs a debug entry.
func (l *Logging) Warn(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Warn(ns...)
}

// Warn logs a debug entry.
func (l *Logging) Warnf(msg string, args ...interface{}) {
	rlog.Warnf(msg, args...)
}

// Error logs a debug entry.
func (l *Logging) Error(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Error(ns...)
}

// Error logs a debug entry.
func (l *Logging) Errorf(msg string, args ...interface{}) {
	rlog.Errorf(msg, args...)
}

// Fatal logs a debug entry.
func (l *Logging) Fatal(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Critical(ns...)
}

// Fatal logs a debug entry.
func (l *Logging) Fatalf(msg string, args ...interface{}) {
	rlog.Criticalf(msg, args...)
}

// Log logs a debug entry.
func (l *Logging) Log(lvl int, msf string, args []interface{}) {

}

// IsDebug returns wether logger running in a debugging environment.
func (l *Logging) IsDebug() bool {
	return l.debug == true
}

// FInfo is a function that  logs a debug entry.
func FInfo(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Info(ns...)
}

// FDebug is a function that  logs a debug entry.
func FDebug(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Debug(ns...)
}

// FWarn is a function that  logs a debug entry.
func FWarn(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns := make([]interface{}, 0, len(args)+2)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Warn(ns...)
}

// FError is a function that  logs a debug entry.
func FError(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Error(ns...)
}

// FFatal is a function that  logs a debug entry.
func FFatal(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns[0] = format
	copy(ns[1:], args)
	rlog.Critical(ns...)
}

// Trace is a function that logs a debug entry.
func Trace(msg string, args ...interface{}) {
}

// Debug is a function that  logs a debug entry.
func Debug(msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Debug(ns...)
}

// Info is a function that  logs a debug entry.
func Info(msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Info(ns...)
}

// Warn is a function that  logs a debug entry.
func Warn(msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Warn(ns...)
}

// Error is a function that  logs a debug entry.
func Error(msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Critical(ns...)
}

// Fatal is a function that  logs a debug entry.
func Fatal(msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Critical(ns...)
}

// Log is a function that  logs a debug entry.
func Log(lvl int, msf string, args []interface{}) {
}
