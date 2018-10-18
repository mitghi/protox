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

func (mb *MessageBox) AddInbound(msg protobase.EDProtocol) bool {
	mb.RLock()
	defer mb.RUnlock()

	var (
		cid string = uidstr(msg)
	)

	mb.in.Lock()

	if _, ok := mb.in.messages[cid]; ok {
		mb.in.Unlock()
		return false
	}
	mb.in.messages[cid] = msg

	mb.in.Unlock()

	return true
}

func (mb *MessageBox) AddOutbound(msg protobase.EDProtocol) bool {
	mb.RLock()
	var cid string = uidstr(msg)
	mb.out.Lock()

	if _, ok := mb.out.messages[cid]; ok == true {
		mb.out.Unlock()
		mb.RUnlock()

		return false
	}
	mb.out.messages[cid] = msg

	mb.out.Unlock()
	mb.RUnlock()

	return true
}

func (mb *MessageBox) DeleteIn(msg protobase.EDProtocol) bool {
	mb.RLock()
	var cid string = uidstr(msg)
	mb.in.Lock()

	if _, ok := mb.in.messages[cid]; ok == false {
		mb.in.Unlock()
		mb.RUnlock()

		return false
	}
	delete(mb.in.messages, cid)

	mb.in.Unlock()
	mb.RUnlock()

	return true
}

// DeleteOut disassociates a client from a outgoing packet.
func (mb *MessageBox) DeleteOut(msg protobase.EDProtocol) bool {
	mb.RLock()
	var cid string = uidstr(msg)
	mb.out.Lock()

	if _, ok := mb.out.messages[cid]; ok == false {
		mb.out.Unlock()
		mb.RUnlock()

		return false
	}
	delete(mb.out.messages, cid)

	mb.out.Unlock()
	mb.RUnlock()

	return true
}

// GetAllOut returns all of available outgoing packets of a given client.
func (mb *MessageBox) GetAllOut() (msgs []protobase.EDProtocol) {
	mb.RLock()
	mb.out.Lock()

	for _, msg := range mb.out.messages {
		msgs = append(msgs, msg)
	}
	order := mb.out.order
	sort.Slice(msgs, func(i, j int) bool {
		a, b := msgs[i], msgs[j]
		astr, bstr := uidstr(a), uidstr(b)
		return order[astr] < order[bstr]
	})

	mb.out.Unlock()
	mb.RUnlock()

	return msgs
}

func (mb *MessageBox) GetAllOutStr() (msgs []string) {
	mb.RLock()
	mb.out.Lock()

	for _, msg := range mb.out.messages {
		var sid string = uidstr(msg)
		msgs = append(msgs, sid)
	}

	mb.out.Unlock()
	mb.RUnlock()

	return msgs
}

func (mb *MessageBox) GetInbound(uid uuid.UUID) (protobase.EDProtocol, bool) {
	mb.RLock()
	mb.in.Lock()

	p, ok := mb.in.messages[uid.String()]

	mb.in.Unlock()
	mb.RUnlock()

	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (mb *MessageBox) GetOutbound(uid uuid.UUID) (protobase.EDProtocol, bool) {
	mb.RLock()
	mb.out.Lock()

	p, ok := mb.out.messages[uid.String()]

	mb.out.Unlock()
	mb.RUnlock()

	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (mb *MessageBox) GetIDStoreO() (idstore protobase.MSGIDInterface) {
	mb.RLock()
	mb.out.Lock()

	idstore = mb.in.ids

	mb.out.Unlock()
	mb.RUnlock()
	return idstore
}

func (mb *MessageBox) GetIDStoreI() (idstore protobase.MSGIDInterface) {
	mb.RLock()
	mb.in.Lock()

	idstore = mb.in.ids

	mb.in.Unlock()
	mb.RUnlock()
	return idstore
}
