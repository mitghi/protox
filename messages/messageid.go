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
	"sync"

	"github.com/google/uuid"
)

// MessageId is a struct which contains mapping of packet number to their unique identifier.
type MessageId struct {
	sync.RWMutex
	id     map[uint16]uuid.UUID
	cursor int
}

type QueueId struct {
	sync.RWMutex
	id     map[uint16]struct{}
	cursor int
}

// - MARK: Initializers.

// NewMessageId returns a pointer to a new `MessageId` struct.
func NewMessageId() *MessageId {
	var result *MessageId = &MessageId{
		id: make(map[uint16]uuid.UUID),
	}
	return result
}

// NewQueueId returns a pointer to a new initialized and allocated
// `QueueId` struct.
func NewQueueId() *QueueId {
	var qi *QueueId = &QueueId{
		id: make(map[uint16]struct{}),
	}
	return qi
}

// - MARK: MessageId section.

// GetNewId finds an empty slot and returns a new `uint16`.
func (m *MessageId) GetNewID(uid uuid.UUID) (id uint16) {
	m.Lock()
	defer m.Unlock()
	id, _ = m.getNewId(uid)
	return id
}

func (m *MessageId) getNewId(uid uuid.UUID) (uint16, int) {
	var i uint16
	if m.cursor > MSGMAXLEN {
		m.cursor = 0
	}
	for i = MSGMINLEN; i < MSGMAXLEN; i++ {
		if _, ok := m.id[i]; ok == false {
			m.id[i] = uid
			m.cursor++
			return i, m.cursor
		}
	}
	return 0, m.cursor
}

// GetFNewID finds an empty slot and returns a new `uint16` associated
// with that slot as well as a `cursor`. Cursor increments on each
// new association and restarts when maximum message length is reached.
func (m *MessageId) GetFNewID(uid uuid.UUID) (id uint16, cursor int) {
	m.Lock()
	defer m.Unlock()
	return m.getNewId(uid)
}

// GetUUID finds the associated `uuid.UUID` for a given`id`. It
func (m *MessageId) GetUUID(id uint16) (uuid.UUID, bool) {
	m.RLock()
	uid, ok := m.id[id]
	m.RUnlock()
	return uid, ok
}

// IsOccupied returns a `bool` indicating wether a certain id is in use or not.
func (m *MessageId) IsOccupied(id uint16) bool {
	m.Lock()
	_, status := m.id[id]
	m.Unlock()
	return status
}

// FreeId removes a `id` from internal mapping. Note that
// it has no effect on cursor ( cursor is not decremented ).
func (m *MessageId) FreeId(id uint16) {
	m.Lock()
	delete(m.id, id)
	m.Unlock()
}

// - MARK: QueueId section.

// GetNewId finds an empty slot and returns a new `uint16`.
func (qi *QueueId) GetNewID() (id uint16) {
	qi.Lock()
	defer qi.Unlock()
	id, _ = qi.getNewId()
	return id
}

func (qi *QueueId) getNewId() (uint16, int) {
	var i uint16
	if qi.cursor > MSGMAXLEN {
		qi.cursor = 0
	}
	for i = MSGMINLEN; i < MSGMAXLEN; i++ {
		if _, ok := qi.id[i]; ok == false {
			qi.id[i] = struct{}{}
			qi.cursor++
			return i, qi.cursor
		}
	}
	return 0, qi.cursor
}

// GetFNewID finds an empty slot and returns a new `uint16` associated
// with that slot as well as a `cursor`. Cursor increments on each
// new association and restarts when maximum message length is reached.
func (qi *QueueId) GetFNewID() (id uint16, cursor int) {
	qi.Lock()
	defer qi.Unlock()
	return qi.getNewId()
}

// IsOccupied returns a `bool` indicating wether a certain id is in use or not.
func (qi *QueueId) IsOccupied(id uint16) bool {
	qi.Lock()
	_, status := qi.id[id]
	qi.Unlock()
	return status
}

// FreeId removes a `id` from internal mapping. Note that
// it has no effect on cursor ( cursor is not decremented ).
func (qi *QueueId) FreeId(id uint16) {
	qi.Lock()
	delete(qi.id, id)
	qi.Unlock()
}
