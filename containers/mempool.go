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

package containers

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

// State constants
const (
	MPFixAlloc byte = iota
	MPNone
)

// BPDefaultRate is default recycle rate in which unused data will be
// either freed or recycled.
const BPDefaultRate = 500 * time.Millisecond

// BufferPool is the pool structure. It communicates via two channels,
// `Get` for new memory and `Release` for retaining/releasing unused memory.
// It starts a recycling round at each `recrate` tick.
type BufferPool struct {
	stat struct {
		mempf *runtime.MemStats
		start time.Time
	}
	Get     chan []byte
	Release chan []byte
	recrate time.Duration
	counter int
	frees   int
	blen    int
	bcap    int
	started bool
}

type (
	// BPNode is the container for individual memory chunks. It writes current
	// time in `when` which is used by `BufferPool` for recycling.
	BPNode struct {
		when time.Time
		data []byte
	}
	// q is internal memory container structure implemented as `BPNode` slice.
	q []BPNode
)

// BufferPoolContext is a context manager which holds into `size` items before releasing
// them to `BufferPool`. It is useful for providing memory chunks during the burst times.
type BufferPoolContext struct {
	Context  [][]byte
	get      <-chan []byte
	release  chan<- []byte
	size     int
	courser  int
	count    int
	occupied int
}

// newBufferPoolContext creates a context with capacity of `size` and connects it
// to the given BufferPool `bp` channels.
func newBufferPoolContext(size int, bp *BufferPool) *BufferPoolContext {
	return &BufferPoolContext{
		Context:  make([][]byte, size),
		get:      bp.Get,
		release:  bp.Release,
		size:     size,
		courser:  0,
		count:    0,
		occupied: 0,
	}
}

// TotalBytes returns total bytes 'used' by the nodes in the `Context`.
func (self *BufferPoolContext) TotalBytes() (inuse int) {
	for i := 0; i < self.size; i++ {
		if slot := self.Context[i]; slot != nil {
			inuse += len(slot)
		}
	}
	return inuse
}

// RefreshStat refreshes total number of occupied slots in `Context` and returns it.
func (self *BufferPoolContext) RefreshStat() (occupied int) {
	for i := 0; i < self.size; i++ {
		if self.Context[i] != nil {
			occupied++
		}
	}
	self.occupied = occupied
	return occupied
}

// Get returns a new memory chunk from the current Context. When no memory is allocated or
// either released, it fetches a new chunk from the `BufferPool` channels and returns it.
func (self *BufferPoolContext) Get() []byte {
	for i := 0; i < self.size; i++ {
		if self.Context[i] != nil {
			chunk := self.Context[i]
			self.Context[i] = nil
			self.count--
			return chunk
		}
	}
	chunk := <-self.get
	return chunk
}

// Release either retains a chunk by writing it to the Context slice and when
// Context is full, it releases the memory to the `BufferPool`.
func (self *BufferPoolContext) Release(chunk []byte) {
	if ok := self.insert(chunk); ok == false {
		self.release <- chunk
	} else {
		self.count++
	}
}

// Flush releases all occupied slots in the current Context to the `BufferPool`.
func (self *BufferPoolContext) Flush() {
	for i := 0; i < self.size; i++ {
		if self.Context[i] != nil {
			self.count--
			self.release <- self.Context[i]
		}
	}
}

// insert finds an empty slot and inserts buff into it. It returns `false` when
// all slots are occupied.
func (self *BufferPoolContext) insert(buff []byte) bool {
	for i := 0; i < self.size; i++ {
		if self.Context[i] == nil {
			self.Context[i] = buff
			self.count++
			return true
		}
	}
	return false
}

// add pushes a new `BPNode` struct into the slice.
func (self *q) add(b []byte) {
	*self = append(*self, BPNode{
		when: time.Now(),
		data: b,
	})
}

// deleteNewest removes most recent nodes from the slice.
func (self *q) deleteNewest() {
	q := *self
	q[len(q)-1] = BPNode{}
	*self = q[0 : len(q)-1]
}

// newest returns the chunk from newest node in the slice.
func (self q) newest() []byte {
	return self[len(self)-1].data
}

// deleteOlderThan removes all node.when<t from the slice.
func (self *q) deleteOlderThan(t time.Time) {
	q := *self
	inactiveCount := len(q)
	for i, e := range q {
		if e.when.After(t) {
			inactiveCount = i
			break
		}
	}
	// Copy all active elements to the start of the slice.
	copy(q, q[inactiveCount:])
	activeCount := len(q) - inactiveCount
	// Zero out all inactive elements of the q
	for j := activeCount; j < len(q); j++ {
		q[j] = BPNode{}
	}
	*self = q[0:activeCount]
}

// NewBufferPool allocates and returns a new `BufferPool`. It will configure
// it with recycle rate `recrate`, chunk allocation len `blen` and chunk
// capacity `bcap`.
func NewBufferPool(recrate time.Duration, blen int, bcap int) *BufferPool {
	result := &BufferPool{
		stat: struct {
			mempf *runtime.MemStats
			start time.Time
		}{mempf: nil, start: time.Now()},
		Get:     make(chan []byte),
		Release: make(chan []byte),
		recrate: recrate,
		counter: 0,
		frees:   0,
		blen:    blen,
		bcap:    bcap,
		started: false,
	}
	return result
}

// alloc allocates a new chunk of bytes with parameters given during
// initialization. It returns newly allocated memory.
func (self *BufferPool) alloc() []byte {
	self.counter += 1
	return make([]byte, self.blen, self.bcap)
}

// Run is a function that starts the main algorithm for memory pool and
// buffer recycling. It returns its Get and Release channels as well as
// an error. If `Run(...)` called twice, it will return an error.
func (self *BufferPool) Run() (<-chan []byte, chan<- []byte, error) {
	if self.started == true {
		return nil, nil, errors.New("BufferPool: already started.")
	}
	self.started = true

	go func() {
		timer := time.NewTimer(self.recrate)
		var q q
		for {
			if len(q) == 0 {
				q.add(self.alloc())
			}
			select {
			case b := <-self.Release:
				q.add(b)

			case self.Get <- q.newest():
				q.deleteNewest()

			case <-timer.C:
				q.deleteOlderThan(time.Now().Add(-self.recrate))
				timer.Reset(self.recrate)
			}
		}
	}()

	return self.Get, self.Release, nil
}

// ResetStat refreshes memory stats.
func (self *BufferPool) ResetStat() {
	self.stat.mempf = nil
	self.stat.mempf = &runtime.MemStats{}
	self.stat.start = time.Now()
}

// MemStat returns memory usage details provided by `runtime` package. Prior to
// a call to this function, memory profile must be refreshes by using `ResetStat`.
// It returns an string with SystemHeap, IdleHeap, ReleaseHeap, current counter and
// current free nodes.
func (self *BufferPool) MemStat() (*string, error) {
	if self.stat.mempf == nil {
		return nil, errors.New("BufferPool: memory stat is not set.")
	}
	runtime.ReadMemStats(self.stat.mempf)
	mempf := self.stat.mempf
	stat := fmt.Sprintf("%d,%d,%d,%d,%d,%d\n", mempf.HeapSys, mempf.HeapAlloc, mempf.HeapIdle, mempf.HeapReleased, self.counter, self.frees)

	return &stat, nil
}

// CreateNewContext allocates and returns a new `BufferPoolContext` from
// the current `BufferPool`.
func (self *BufferPool) CreateNewContext(size int) *BufferPoolContext {
	result := newBufferPoolContext(size, self)
	return result
}
