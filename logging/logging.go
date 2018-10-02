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

	"github.com/mitghi/rlog"
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
func (self *Logging) SetDebug(status bool) {
	self.debug = status
}

// FInfo logs a entry with caller name prepended.
func (self *Logging) FInfo(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Info(ns...)
}

// FInfo logs a entry with caller name prepended.
func (self *Logging) FInfof(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("<- (*%s).%s", self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Infof(msg, ns...)
}

// FDebug logs a entry with caller name prepended.
func (self *Logging) FDebug(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Debug(ns...)
}

func (self *Logging) FDebugf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("%s <- (*%s).%s", msg, self.Package, fn)
	rlog.Debugf(format, args...)
}

// FWarn logs a entry with caller name prepended.
func (self *Logging) FWarn(fn string, msg string, args ...interface{}) error {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Warn(ns...)
	return nil
}

// FWarn logs a entry with caller name prepended.
func (self *Logging) FWarnf(fn string, msg string, args ...interface{}) error {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Warnf(msg, ns...)
	return nil
}

// FError logs a entry with caller name prepended.
func (self *Logging) FError(fn string, msg string, args ...interface{}) error {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Error(ns...)
	return nil
}

// FError logs a entry with caller name prepended.
func (self *Logging) FErrorf(fn string, msg string, args ...interface{}) error {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Errorf(msg, ns...)
	return nil
}

// FFatal logs a entry with caller name prepended.
func (self *Logging) FFatal(fn string, msg string, args ...interface{}) {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns[0] = format
	copy(ns[1:], args)
	rlog.Critical(ns...)
}

// FFatal logs a entry with caller name prepended.
func (self *Logging) FFatalf(fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf("<- (*%s).%s", self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Criticalf(msg, ns...)
}

// Trace logs a trace entry.
func (self *Logging) FTrace(level int, fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	ns := formatArgs(append([]interface{}{format}, args...)...)
	rlog.Trace(level, ns...)
}

// Trace logs a trace entry.
func (self *Logging) FTracef(level int, fn string, msg string, args ...interface{}) {
	format := fmt.Sprintf(logformat, msg, self.Package, fn)
	rlog.Tracef(level, format, args...)
}

// Trace logs a trace entry.
func (self *Logging) Trace(level int, args ...interface{}) {
	rlog.Trace(level, args...)
}

// Trace logs a trace entry.
func (self *Logging) Tracef(level int, msg string, args ...interface{}) {
	rlog.Tracef(level, msg, args...)
}

// Debug logs a debug entry.
func (self *Logging) Debug(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Debug(ns...)
}

// Debug logs a debug entry.
func (self *Logging) Debugf(msg string, args ...interface{}) {
	rlog.Debugf(msg, args...)
}

// Info logs a debug entry.
func (self *Logging) Info(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Info(ns...)
}

// Info logs a debug entry.
func (self *Logging) Infof(msg string, args ...interface{}) {
	rlog.Infof(msg, args...)
}

// Warn logs a debug entry.
func (self *Logging) Warn(msg string, args ...interface{}) error {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Warn(ns...)
	return nil
}

// Warn logs a debug entry.
func (self *Logging) Warnf(msg string, args ...interface{}) error {
	rlog.Warnf(msg, args...)
	return nil
}

// Error logs a debug entry.
func (self *Logging) Error(msg string, args ...interface{}) error {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Error(ns...)
	return nil
}

// Error logs a debug entry.
func (self *Logging) Errorf(msg string, args ...interface{}) error {
	rlog.Errorf(msg, args...)
	return nil
}

// Fatal logs a debug entry.
func (self *Logging) Fatal(msg string, args ...interface{}) {
	ns := formatArgs(append([]interface{}{msg}, args...)...)
	rlog.Critical(ns...)
}

// Fatal logs a debug entry.
func (self *Logging) Fatalf(msg string, args ...interface{}) {
	rlog.Criticalf(msg, args...)
}

// Log logs a debug entry.
func (self *Logging) Log(lvl int, msf string, args []interface{}) {

}

// IsDebug returns wether logger running in a debugging environment.
func (self *Logging) IsDebug() bool {
	return self.debug == true
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
func FWarn(fn string, msg string, args ...interface{}) error {
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns := make([]interface{}, 0, len(args)+2)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Warn(ns...)
	return nil
}

// FError is a function that  logs a debug entry.
func FError(fn string, msg string, args ...interface{}) error {
	ns := make([]interface{}, 0, len(args)+2)
	format := fmt.Sprintf("%s -> %s", fn, msg)
	ns = append(ns, format)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Error(ns...)
	return nil
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
func Warn(msg string, args ...interface{}) error {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Warn(ns...)
	return nil
}

// Error is a function that  logs a debug entry.
func Error(msg string, args ...interface{}) error {
	ns := make([]interface{}, 0, len(args)+1)
	ns = append(ns, msg)
	if len(args) > 0 {
		ns = append(ns, args...)
	}
	rlog.Critical(ns...)
	return nil
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
