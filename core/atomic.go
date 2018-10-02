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

package core

import (
	"errors"
	"sync/atomic"
)

const (
	ACActive byte = iota
	ACExpired
)

/**
*
* This package includes a set of atomic tools.
*
*
**/

// AtomicCounter is function signature for the counter function
type AtomicCounter func() int64

// ExpiringAtomicCounter is the function signature for expiring counter function
type ExpiringAtomicCounter func() (int64, byte)

// CTCallback is function signature for the callback invoked when counter reaches the
// threshold
type CTCallback func(int64)

// CTExpCallback is function signature for the callback invoked when a already expired
// counter is reused
type CTExpCallback func(int64)

// NewAtomicCounter is a function that creates a new counter, counting from `start` to
// `threshold` and returns it as a `AtomicCounter` function. It resets the counter
// once `threshold` is reached, therefore can be reused.
func NewAtomicCounter(start, threshold int64, cb CTCallback) (AtomicCounter, error) {
	if threshold < 1 || threshold < start {
		return nil, errors.New("AtomicCounter: invalid threshold value.")
	}
	var (
		counter int64 = start
		ret     AtomicCounter
		curr    int64
	)
	ret = func() int64 {
		curr = atomic.AddInt64(&counter, 1)
		if curr >= counter && curr%threshold == 0 {
			atomic.AddInt64(&counter, -threshold)
			cb(threshold)
		}
		return curr
	}
	return ret, nil
}

// NewExpiringAtomicCounter is a function that creates a new expiring counter,
// counting from `start` to `threashold`. It expires the counter once `threshold`
// value is reaches and any further calls to the counter function results in invokation
// of second callback `ecb` with the status code set to `ACExpired` in returned values.
func NewExpiringAtomicCounter(start, threshold int64, cb CTCallback, ecb CTExpCallback) (ExpiringAtomicCounter, error) {
	if threshold < 1 || threshold < start {
		return nil, errors.New("AtomicCounter: invalid threshold value.")
	}
	var (
		counter int64 = start
		ret     ExpiringAtomicCounter
		curr    int64
		hasExp  bool
	)
	ret = func() (int64, byte) {
		if hasExp == true {
			ecb(threshold)
			return threshold, ACExpired
		}
		curr = atomic.AddInt64(&counter, 1)
		if curr >= counter && curr%threshold == 0 {
			cb(threshold)
			hasExp = true
		}
		return curr, ACActive
	}
	return ret, nil
}
