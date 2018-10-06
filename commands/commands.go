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

package commands

import (
	"bytes"
	"io"
)

// Ensure interface (protocol) conformance
var _ io.Writer = (*Command)(nil)

type position struct {
	start int
	end   int
}

// Command is a buffer for creating
// and serializing messages.
type Command struct {
	*bytes.Buffer
	parent   *Command
	totalcnt int
	valcnt   int
	size     int
	commited bool
}

// NewCommand allocate and initializes a new
// `Command` and returns a pointer to it.
func NewCommand() *Command {
	return &Command{
		Buffer: &bytes.Buffer{},
	}
}

// Write writes `p` bytes into the buffer. It returns
// number of written bytes and an error in case of
// unsuccessfull operation.
func (c *Command) Write(p []byte) (n int, err error) {
	n, err = c.Buffer.Write(p)
	c.size += n
	return n, err
}

// WriteString writes an string to the buffer. It returns
// number of written bytes and an error in case of
// unsuccessfull operation.
func (c *Command) WriteString(s string) (n int, err error) {
	n, err = c.Buffer.WriteString(addString(s))
	c.size += n
	c.totalcnt += len(s)
	c.valcnt++
	return n, err
}

// WriteArrayHeader is a function that creates an array header
// with its number of values set to `nval` and total length set
// to `tl`. It returns number of written bytes and an error in
// in case of unsuccessful attempt.
func (c *Command) WriteArrayHeader(nval int, tl int) (n int, err error) {
	n, err = c.Buffer.WriteString(addArrayHeader(nval, tl))
	c.size += n
	c.valcnt++
	return n, err
}

// WriteArrayItem is a function that writes an item associated
// to underlaying array. It returns number of written bytes
// and an error in case of unsuccesfull operation.
func (c *Command) WriteArrayItem(s string) (n int, err error) {
	n, err = c.Buffer.WriteString(addString(s))
	c.size += n
	c.totalcnt += len(s)
	if c.parent != nil {
		c.valcnt++
	}
	return n, err
}

// WriteArrayString creates an array from variadic strings
// `args` and writes it to the buffer. It returns number of
// written bytes and an error in case of unsuccessfull operation.
func (c *Command) WriteArrayString(args ...string) (n int, err error) {
	n, err = c.Buffer.WriteString(addArray(args...))
	c.totalcnt += n
	c.size += n
	c.valcnt++
	return n, err
}

// WriteArray takes a pointer to a `Command` struct and creates
// an array from its content. It returns number of written bytes
// and an error in case of unsuccesfull error.
func (c *Command) WriteArray(a *Command) (n int, err error) {
	var (
		rn int64
	)
	n, err = c.WriteArrayHeader(a.valcnt, a.totalcnt)
	if err != nil {
		a.commited = false
		return n, err
	}
	rn, err = c.ReadFrom(a.Buffer)
	if err != nil {
		a.commited = false
		return int(rn), err
	}
	a.commited = true
	return int(rn), err
}

// MakeArray allocate and initializes a new buffer for
// array and returns a pointer to it. It sets the parent
// pointer to itself.
func (c *Command) MakeArray() *Command {
	var (
		nc *Command = NewCommand()
	)
	nc.parent = c
	return nc
}

// CommitArray commits its content by calling a receiver
// function on its parent. The parent creates a new array
// header with statistics from caller and reads the caller
// bytes into its own buffer.
func (c *Command) CommitArray() bool {
	if c.parent == nil {
		return false
	}
	c.parent.WriteArray(c)
	return c.commited
}

// Reset releases underlaying byte buffer and
// zeros the statistics.
func (c *Command) Reset() {
	c.size = 0
	c.totalcnt = 0
	c.valcnt = 0
	c.commited = false
	c.Buffer.Reset()
}
