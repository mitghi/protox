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
	"sort"
	"sync"

	"github.com/google/uuid"

	"github.com/mitghi/protox/protobase"
)

// Constant values for minimum and maximum number of packets. In reality, this number
// will not be reached easily as the enteries gets reused when they become free.
const (
	MSGMINLEN = 1
	MSGMAXLEN = 65535
)

// MsgEntry maps a client id to a protocol packet.
type MsgEntry struct {
	sync.Mutex
	messages map[string]protobase.EDProtocol
	order    map[string]int
	ids      *MessageId
	counter  int
}

// MessageStore is a struct which acts as an entry into the persistent storage.
type MessageStore struct {
	sync.RWMutex
	in  map[string]*MsgEntry
	out map[string]*MsgEntry
}

// NewMessageStore returns a pointer to a new `MessageStore`.
func NewMessageStore() *MessageStore {
	// TODO
	// . review and refactor to remove this function.
	var result *MessageStore = &MessageStore{}
	return result
}

// NewInitedMessageStore returns a pointer to a new `MessageStore` and allocates memory.
// This is equivalent of calling `Init()` after `NewMessageStore`
func NewInitedMessageStore() *MessageStore {
	var result *MessageStore = &MessageStore{
		in:  make(map[string]*MsgEntry),
		out: make(map[string]*MsgEntry),
	}
	return result
}

// Init initializes all of internal struct members such as inbound and
// outbound mappings.
func (self *MessageStore) Init() {
	self.in = make(map[string]*MsgEntry)
	self.out = make(map[string]*MsgEntry)
}

// AddClient initializes and adds a entry for a given client.
func (self *MessageStore) AddClient(client string) {
	self.Lock()
	self.in[client] = &MsgEntry{
		messages: make(map[string]protobase.EDProtocol),
		order:    make(map[string]int),
		ids:      NewMessageId(),
	}
	self.out[client] = &MsgEntry{
		messages: make(map[string]protobase.EDProtocol),
		order:    make(map[string]int),
		ids:      NewMessageId(),
	}
	self.Unlock()
}

// Close removes the entry of the client from internal mappings. It returns `false`
// if client does not exist.
func (self *MessageStore) Close(client string) bool {
	self.Lock()
	if ok := self.nomexist(client); !ok {
		self.Unlock()
		return false
	}
	delete(self.in, client)
	delete(self.out, client)
	self.Unlock()

	return true
}

// AddInbound associates a client to a incoming packet. It checks
// if the given client is in the storage, otherwise returns `false`.
func (self *MessageStore) AddInbound(client string, msg protobase.EDProtocol) bool {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return false
	}
	var cid string = msg.UUID().String()
	self.in[client].Lock()
	if _, ok := self.in[client].messages[cid]; ok == true {
		self.in[client].Unlock()
		self.RUnlock()

		return false
	}
	self.in[client].messages[cid] = msg
	self.in[client].order[cid] = self.in[client].GenSeqID()
	self.in[client].Unlock()
	self.RUnlock()

	return true
}

// AddOutbound associates a client to a ougoing packet.
func (self *MessageStore) AddOutbound(client string, msg protobase.EDProtocol) bool {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return false
	}
	var cid string = msg.UUID().String()
	self.out[client].Lock()
	if _, ok := self.out[client].messages[cid]; ok == true {
		self.out[client].Unlock()
		self.RUnlock()

		return false
	}
	self.out[client].messages[cid] = msg
	self.out[client].order[cid] = self.out[client].GenSeqID()
	self.out[client].Unlock()
	self.RUnlock()

	return true
}

// DeleteIn disassociates a client from a incoming packet.
func (self *MessageStore) DeleteIn(client string, msg protobase.EDProtocol) bool {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return false
	}
	var cid string = msg.UUID().String()
	self.in[client].Lock()
	if _, ok := self.in[client].messages[cid]; ok == false {
		self.in[client].Unlock()
		self.RUnlock()

		return false
	}
	delete(self.in[client].messages, cid)
	delete(self.in[client].order, cid)
	self.in[client].Unlock()
	self.RUnlock()

	return true
}

// DeleteOut disassociates a client from a outgoing packet.
func (self *MessageStore) DeleteOut(client string, msg protobase.EDProtocol) bool {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return false
	}
	var cid string = msg.UUID().String()
	self.out[client].Lock()
	if _, ok := self.out[client].messages[cid]; ok == false {
		self.out[client].Unlock()
		self.RUnlock()

		return false
	}
	delete(self.out[client].messages, cid)
	delete(self.out[client].order, cid)
	self.out[client].Unlock()
	self.RUnlock()

	return true
}

// Exists returns a `bool` indicating whether a client is already registered or not.
func (self *MessageStore) Exists(client string) (ok bool) {
	self.RLock()
	defer self.RUnlock()
	ok = self.nomexist(client)
	return ok
}

// nomexist returns a `bool` indicating whether a client is already registered
// or not. It is used instead of `Exists` internally because its locking is done
// manually by a caller.
func (self *MessageStore) nomexist(client string) bool {
	var (
		okIn  bool
		okOut bool
	)
	_, okIn = self.in[client]
	_, okOut = self.out[client]

	return okIn && okOut
}

// GetAllOut returns all of available outgoing packets of a given client.
func (self *MessageStore) GetAllOut(client string) (msgs []protobase.EDProtocol) {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return msgs
	}
	self.out[client].Lock()
	for _, msg := range self.out[client].messages {
		msgs = append(msgs, msg)
	}
	/* d e b u g */
	// orig := make([]protobase.EDProtocol, len(msgs), len(msgs))
	// copy(orig[:], msgs)
	/* d e b u g */
	order := self.out[client].order
	sort.Slice(msgs, func(i, j int) bool {
		a, b := msgs[i], msgs[j]
		astr, bstr := a.UUID().String(), b.UUID().String()
		return order[astr] < order[bstr]
	})

	/* d e b u g */
	// fmt.Println("sorted")
	// for i := 0; i < len(msgs); i++ {
	// 	fmt.Println(msgs[i].UUID().String())
	// }
	// fmt.Println("--")
	// fmt.Println("unsorted")
	// for i := 0; i < len(orig); i++ {
	// 	fmt.Println(orig[i].UUID().String())
	// }
	/* d e b u g */

	self.out[client].Unlock()
	self.RUnlock()

	return msgs
}

// GetAllOutStr returns string UUID repr of all available outgoing packets
// of a given client.
func (self *MessageStore) GetAllOutStr(client string) (msgs []string) {
	self.RLock()
	if ok := self.nomexist(client); !ok {
		self.RUnlock()
		return msgs
	}
	self.out[client].Lock()
	for _, msg := range self.out[client].messages {
		var uid uuid.UUID = msg.UUID()
		var sid string = uid.String()
		msgs = append(msgs, sid)
	}
	self.out[client].Unlock()
	self.RUnlock()

	return msgs
}

func (self *MessageStore) GetInbound(client string, uid uuid.UUID) (protobase.EDProtocol, bool) {
	if !self.Exists(client) {
		return nil, false
	}
	self.RLock()
	val, ok := self.in[client]
	self.RUnlock()
	if !ok || val == nil {
		return nil, false
	}
	val.Lock()
	p, ok := val.messages[uid.String()]
	val.Unlock()
	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (self *MessageStore) GetOutbound(client string, uid uuid.UUID) (protobase.EDProtocol, bool) {
	if !self.Exists(client) {
		return nil, false
	}
	self.RLock()
	val, ok := self.out[client]
	self.RUnlock()
	if !ok || val == nil {
		return nil, false
	}
	val.Lock()
	p, ok := val.messages[uid.String()]
	val.Unlock()
	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

func (self *MessageStore) GetIDStoreO(client string) (idstore protobase.MSGIDInterface) {
	self.RLock()
	defer self.RUnlock()
	entry, ok := self.out[client]
	if !ok {
		return nil
	}
	idstore = entry.ids
	return idstore
}

func (self *MessageStore) GetIDStoreI(client string) (idstore protobase.MSGIDInterface) {
	self.RLock()
	defer self.RUnlock()
	entry, ok := self.out[client]
	if !ok {
		return nil
	}
	idstore = entry.ids
	return idstore
}

// GenSeqID creates and returns a sequence id used to preserve order.
func (self *MsgEntry) GenSeqID() int {
	sid := self.counter
	self.counter++
	return sid
}

// NextLevepP seperates individual components in a topic string.
func NextLevelP(topic []byte, sep byte, wldcd byte) (nlvl []byte, rem []byte, err error) {
	var (
		state byte = LCHR
	)
	for i, c := range topic {
		switch c {
		case sep:
			if i == 0 {
				err = errors.New("messages: invalid position for seperator.")
				return nlvl, rem, err
			}
			if state == LWLCD {
				// // TODO: NOTE: handle wildcard
			}
			state = LBRK
			nlvl = topic[:i]
			rem = topic[i+1:]
			err = nil
			return nlvl, rem, err
		case wldcd:
			if i != 0 {
				err = errors.New("invalid wildcard position")
				return nlvl, rem, err
			}
			state = LWLCD
		default:
			state = LCHR
		}
	}
	nlvl = topic
	err = nil

	return nlvl, rem, err
}

// TopicComponents is a wrapper func for `DNextLevelP`. It returns
// individual topic components in a slice of byte arrays. Errors
// must be explicitly checked.
func TopicComponents(topic []byte) (res [][]byte, err error) {
	if len(topic) == 0 {
		return
	}
	var (
		t   []byte
		rem []byte = topic
	)
	for {
		t, rem, err = DNextLevelP(rem)
		if err != nil {
			return
		}
		if len(t) == 0 {
			return
		}
		res = append(res, t)
	}
}

// DNextLevelP is a wrapper for `NextLevelP` and supplies
// default topic constants to it.
func DNextLevelP(topic []byte) ([]byte, []byte, error) {
	return NextLevelP(topic, TSEP, TWLDCD)
}
