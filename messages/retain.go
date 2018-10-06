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

package messages

import (
	"errors"

	"github.com/mitghi/protox/protobase"
)

// Retain error messages
var (
	ERNotFound error = errors.New("subscribe: node not found")
	ERINVNode  error = errors.New("subscribe: inconsistent / invalid node")
)

// Retain is the container for retained messages.
type Retain struct {
	topic  string
	packet protobase.EDProtocol
	next   map[string]*Retain
}

// NewRetain allocates and initializes a new `Retain`
// struct and returns a pointer to it. It also allocates
// the internal `map`.
func NewRetain() *Retain {
	return &Retain{next: make(map[string]*Retain)}
}

// Insert inserts/replace a node and its associated `packet`
// at `topic` path.
func (self *Retain) Insert(topic []byte, packet protobase.EDProtocol) (err error) {
	return self.rinsert(topic, packet)
}

// Find finds a node associated with `topic` and returns
// its associated packet.
func (self *Retain) Find(topic []byte) (packet protobase.EDProtocol, err error) {
	return self.rfind(topic)
}

// Remove removes a node associated to `topic`.
func (self *Retain) Remove(topic []byte) (err error) {
	return self.rremove(topic)
}

// rinsert is a receiver method that recursively traverse
// the tree and insert `packet` argument into the appropirate
// node. It creates missing levels during recursion.
func (self *Retain) rinsert(topic []byte, packet protobase.EDProtocol) (err error) {
	if len(topic) == 0 {
		self.packet = packet
		return nil
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return err
	}
	lvl := string(nt)
	n, ok := self.next[lvl]
	if !ok {
		n = NewRetain()
		n.topic = lvl
		n.packet = nil
		self.next[lvl] = n
	}

	return n.rinsert(rem, packet)
}

// rfind is a receiver method that recursively traverse
// the tree to find edge node ( i.e. len(topic) == 0 ). It
// returns an error to indicate unsuccessfull operation.
func (self *Retain) rfind(topic []byte) (packet protobase.EDProtocol, err error) {
	if len(topic) == 0 {
		if self.packet == nil {
			return nil, ERINVNode
		}
		return self.packet, nil
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return nil, err
	}
	lvl := string(nt)
	n, ok := self.next[lvl]
	if !ok {
		return nil, ERNotFound
	}
	return n.rfind(rem)
}

// rremove is a receiver method that recursively traverse
// the tree to find appropirate node and removes it. It
// returns an error to indicate unsuccessfull operation.
// When a node becomes a leaf ( 0 children ), it gets
// removed as well.
func (self *Retain) rremove(topic []byte) error {
	if len(topic) == 0 {
		if self.packet == nil {
			return ERNotFound
		}
		self.packet = nil
		self.next = nil
		return nil
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return err
	}
	lvl := string(nt)
	n, ok := self.next[lvl]
	if !ok {
		return ERNotFound
	}
	if err := n.rremove(rem); err != nil {
		return err
	}
	if len(self.next) == 0 {
		delete(self.next, lvl)
	}

	return nil
}
