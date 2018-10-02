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

package server

import (
	"fmt"
	"sync/atomic"

	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/utils/strs"
)

/*
*
* NOTE: this package is not THREAD-SAFE, INTENTIONALLY.
* Lock must be explicitely acquired before using any
* of these subroutines.
*
 */

/* 
* TODO:
* . decouple into separate package
*/

// subcacheline is a simple struct that holds individual
// subscription information for each cache line. NOTE: not
// thread safe, should manually acquire the lock before using
// this receiver.
type subcacheline struct {
	line []*subscription
}

// subcache is the cache struct for subscriptions. It contains
// a map used to associate routes to cache lines and statistical
// information such as successfull cache hits and .... . NOTE: not
// protected by mutex therefore not thread safe, lock should be
// manually acquired before mutation attempts.
type subcache struct {
	// TODO
	// . check alignment
	cache   map[string]*subcacheline
	hits    uint64
	removes uint64
	inserts uint64
	matches uint64
}

// add is a receiver method that appends a new subscription
// in to its storage. NOTE: not thread safe, should manually
// acquire the lock before using this receiver.
func (scl *subcacheline) add(sb *subscription) {
	scl.line = append(scl.line, sb)
}

// clone is a receiver method that copies subscription data into
// a new slice and returns it. It is used to prevent concurrent
// access to slice data which can have unintended side effects,
// segmentation faults and data corruption. NOTE: not thread safe,
// should manually acquire the lock before using this receiver.
func (scl *subcacheline) clone() (ret []*subscription) {
	if len(scl.line) == 0 {
		return ret
	}
	nl := make([]*subscription, len(scl.line))
	copy(nl, scl.line)
	return nl
}

// GetCacheLine is a receiver method that returns an exact one-to-one
// match from the cache, if it exists along with a boolean value to
// indicate its existence. NOTE: not thread safe, should manually acquire
// the lock before using this receiver.
func (sc *subcache) GetCacheLine(route string) ([]*subscription, bool) {
	var (
		sb  *subcacheline
		ok  bool
		nsb []*subscription
	)
	sb, ok = sc.cache[route]
	if !ok || (sb == nil || (sb != nil && len(sb.line) == 0)) {
		return nil, false
	}
	atomic.AddUint64(&sc.hits, 1)
	nsb = sb.clone()

	return nsb, true
}

// GetCacheLines is a receiver method that returns all subscriptions
// that match `route` argument. NOTE: not thread safe, should manually
// acquire the lock before using this receiver.
func (sc *subcache) GetCacheLines(route string) (lines [][]*subscription) {
	var (
		nsb []*subscription
	)
	for k, v := range sc.cache {
		if !strs.Match(k, route, protobase.Sep, protobase.Wlcd) {
			continue
		}

		nsb = v.clone()
		lines = append(lines, nsb)
	}
	return lines
}

// AddCacheLine is a receiver method that adds a single subscription
// line into the cache NOTE: not thread safe, should manually acquire
// the lock before using this receiver .
func (sc *subcache) AddCacheLine(route string, sb *subscription) {
	if _, ok := sc.cache[route]; !ok {
		sc.cache[route] = &subcacheline{}
	}
	sc.cache[route].add(sb)
}

// NOTE: not thread safe, should manually acquire the lock before
// using this receiver.
func (sc *subcache) AddCacheLines(route string, sb *subscription) {
	sc.AppendToCacheLine(route, sb)
	if _, ok := sc.cache[route]; !ok {
		sc.cache[route] = &subcacheline{}
		sc.cache[route].add(sb)
	}
}

// AppendToCacheLine appends a subscription to all matching lines
// that match `route` argument. NOTE: not thread safe, should manually acquire
// the lock before using this receiver.
func (sc *subcache) AppendToCacheLine(route string, sb *subscription) {
	for k, _ := range sc.cache {
		if strs.Match(k, route, protobase.Sep, protobase.Wlcd) {
			sc.cache[k].add(sb)
		}
	}
}

// RemoveCacheLine removes a single cache line entry ( i.e. single
// one-to-one exact match ). NOTE: not thread safe, should manually acquire
// the lock before using this receiver.
func (sc *subcache) RemoveCacheLine(route string) {
	l, ok := sc.cache[route]
	if !ok || l == nil {
		return
	}
	delete(sc.cache, route)
}

// RemoveCacheLines removes all cache lines that match `route`
// argument. NOTE: not thread safe, should manually acquire
// the lock before using this receiver.
func (sc *subcache) RemoveCacheLines(route string) {
	for k, _ := range sc.cache {
		if strs.Match(k, route, protobase.Sep, protobase.Wlcd) {
			delete(sc.cache, route)
		}
	}
}

// CachePrintV prints cache lines. It is used for debugging purposes.
// NOTE: not thread safe, should manually acquire the lock before
// using this receiver.
func (sc *subcache) CachePrintV() {
	for k, v := range sc.cache {
		fmt.Println("\tTopic", k)
		for _, t := range v.line {
			fmt.Printf("\t\t. Topic(%s), QoS(%d), UID(%s)\n", t.topic, t.qos, t.uid)
		}
	}
}
