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
	"fmt"
)

// Subscribe error messages
var (
	ESNotFound error = errors.New("no such topic or client.")
)

// Subscribe is the second subscriptions struct.
type Subscribe struct {
	topic string
	subs  []string
	qos   []byte
	next  map[string]*Subscribe
}

type substorage struct {
	subs []string
	qos  []byte
}

// NewSub returns a pointer to a new `Subscribe` struct.
func NewSub() *Subscribe {
	result := &Subscribe{
		next: make(map[string]*Subscribe),
	}
	return result
}

// rinesrt recursively inserts a new topic and qos to appropirate node.
func (self *Subscribe) rinsert(topic []byte, qos byte, sub string) error {
	if len(topic) == 0 {
		for i, v := range self.subs {
			if v == string(topic) {
				self.qos[i] = qos
				return nil
			}
		}
		self.subs = append(self.subs, sub)
		self.qos = append(self.qos, qos)

		return nil
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return err
	}
	lvl := string(nt)
	n, ok := self.next[lvl]
	if !ok {
		n = NewSub()
		n.topic = lvl
		self.next[lvl] = n
	}

	return n.rinsert(rem, qos, sub)
}

// rremove recursively deletes entries associated to `client`.
func (self *Subscribe) rremove(topic []byte, client string) error {
	if len(topic) == 0 {
		for i := range self.subs {
			if self.subs[i] == client {
				self.subs = append(self.subs[:i], self.subs[i+1:]...)
				self.qos = append(self.qos[:i], self.qos[i+1:]...)
				return nil
			}
		}
		return ESNotFound
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return err
	}
	lvl := string(nt)
	n, ok := self.next[lvl]
	if !ok {
		return fmt.Errorf("no such topic or client. current: %s nlevel: %v", lvl, rem)
	}
	if err := n.rremove(rem, client); err != nil {
		return err
	}
	if len(n.subs) == 0 && len(n.next) == 0 {
		delete(self.next, lvl)
	}

	return nil
}

// rfind recursively finds all subscribers of a particular topic with
// certain quality of service.
func (self *Subscribe) rfind(topic []byte, qos byte, storage *substorage) error {
	if len(topic) == 0 {
		for i, v := range self.subs {
			if qos <= self.qos[i] {
				storage.subs = append(storage.subs, v)
				storage.qos = append(storage.qos, qos)
			}
		}
		return nil
	}
	nt, rem, err := DNextLevelP(topic)
	if err != nil {
		return err
	}
	lvl := string(nt)
	for k, n := range self.next {
		if k == "*" {
			for i, v := range n.subs {
				if qos <= n.qos[i] {
					storage.subs = append(storage.subs, v)
					storage.qos = append(storage.qos, qos)
				}
			}
		} else if k == lvl {
			if err := n.rfind(rem, qos, storage); err != nil {
				return err
			}
		}
	}
	return nil
}
