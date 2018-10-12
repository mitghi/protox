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
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/mitghi/protox/containers"
	"github.com/mitghi/protox/protobase"
)

type MessageBox struct {
	sync.RWMutex
	in  *MsgEntry
	out *MsgEntry
}

type QueueBoxEntry struct {
	ids *QueueId
}

type QueueBox struct {
	sync.RWMutex
	in  *containers.OrderedSet
	out *containers.OrderedSet
}

func NewMessageBox() *MessageBox {
	var result *MessageBox = &MessageBox{
		in: &MsgEntry{
			messages: make(map[string]protobase.EDProtocol),
			order:    make(map[string]int),
			ids:      NewMessageId(),
		},
		out: &MsgEntry{
			messages: make(map[string]protobase.EDProtocol),
			order:    make(map[string]int),
			ids:      NewMessageId(),
		},
	}
	return result
}

var uidstr func(protobase.EDProtocol) string = func(msg protobase.EDProtocol) string {
	return uuid.UUID(msg.UUID()).String()
}

func (self *MessageBox) AddInbound(msg protobase.EDProtocol) bool {
	self.RLock()
	defer self.RUnlock()

	var (
		cid string = uidstr(msg)
	)

	self.in.Lock()

	if _, ok := self.in.messages[cid]; ok {
		self.in.Unlock()
		return false
	}
	self.in.messages[cid] = msg

	self.in.Unlock()

	return true
}

func (self *MessageBox) AddOutbound(msg protobase.EDProtocol) bool {
	self.RLock()
	var cid string = uidstr(msg)
	self.out.Lock()

	if _, ok := self.out.messages[cid]; ok == true {
		self.out.Unlock()
		self.RUnlock()

		return false
	}
	self.out.messages[cid] = msg

	self.out.Unlock()
	self.RUnlock()

	return true
}

func (self *MessageBox) DeleteIn(msg protobase.EDProtocol) bool {
	self.RLock()
	var cid string = uidstr(msg)
	self.in.Lock()

	if _, ok := self.in.messages[cid]; ok == false {
		self.in.Unlock()
		self.RUnlock()

		return false
	}
	delete(self.in.messages, cid)

	self.in.Unlock()
	self.RUnlock()

	return true
}

// DeleteOut disassociates a client from a outgoing packet.
func (self *MessageBox) DeleteOut(msg protobase.EDProtocol) bool {
	self.RLock()
	var cid string = uidstr(msg)
	self.out.Lock()

	if _, ok := self.out.messages[cid]; ok == false {
		self.out.Unlock()
		self.RUnlock()

		return false
	}
	delete(self.out.messages, cid)

	self.out.Unlock()
	self.RUnlock()

	return true
}

// GetAllOut returns all of available outgoing packets of a given client.
func (self *MessageBox) GetAllOut() (msgs []protobase.EDProtocol) {
	self.RLock()
	self.out.Lock()

	for _, msg := range self.out.messages {
		msgs = append(msgs, msg)
	}
	order := self.out.order
	sort.Slice(msgs, func(i, j int) bool {
		a, b := msgs[i], msgs[j]
		astr, bstr := uidstr(a), uidstr(b)
		return order[astr] < order[bstr]
	})

	self.out.Unlock()
	self.RUnlock()

	return msgs
}

func (self *MessageBox) GetAllOutStr() (msgs []string) {
	self.RLock()
	self.out.Lock()

	for _, msg := range self.out.messages {
		var sid string = uidstr(msg)
		msgs = append(msgs, sid)
	}

	self.out.Unlock()
	self.RUnlock()

	return msgs
}

func (self *MessageBox) GetInbound(uid uuid.UUID) (protobase.EDProtocol, bool) {
	self.RLock()
	self.in.Lock()

	p, ok := self.in.messages[uid.String()]

	self.in.Unlock()
	self.RUnlock()

	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (self *MessageBox) GetOutbound(uid uuid.UUID) (protobase.EDProtocol, bool) {
	self.RLock()
	self.out.Lock()

	p, ok := self.out.messages[uid.String()]

	self.out.Unlock()
	self.RUnlock()

	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (self *MessageBox) GetIDStoreO() (idstore protobase.MSGIDInterface) {
	self.RLock()
	self.out.Lock()

	idstore = self.in.ids

	self.out.Unlock()
	self.RUnlock()
	return idstore
}

func (self *MessageBox) GetIDStoreI() (idstore protobase.MSGIDInterface) {
	self.RLock()
	self.in.Lock()

	idstore = self.in.ids

	self.in.Unlock()
	self.RUnlock()
	return idstore
}
