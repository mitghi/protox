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

	"github.com/mitghi/protox/containers"
)

/**
* TODO:
* . write documentation.
* . write test units.
* . check code coverage.
**/

type MessageQueueAck struct {
	sync.RWMutex

	in  map[int]*sync.Cond
	out map[int]*sync.Cond
}

type MessageQueue struct {
	sync.RWMutex

	Ack  *MessageQueueAck
	Q    *containers.Queue
	Last interface{}
}

func NewMessageQueue() *MessageQueue {
	var mq *MessageQueue = &MessageQueue{
		Ack: NewMessageQueueAck(),
		Q:   containers.NewQueue(),
	}

	return mq
}

func NewMessageQueueAck() *MessageQueueAck {
	var mqa *MessageQueueAck = &MessageQueueAck{
		in:  make(map[int]*sync.Cond),
		out: make(map[int]*sync.Cond),
	}

	return mqa
}

// - MARK: MessageQueue section.

func (mq *MessageQueue) InsertIn(aid int, item interface{}) {
	mq.Lock()
	defer mq.Unlock()

	mq.push(true, aid, item)
}

func (mq *MessageQueue) InsertOut(aid int, item interface{}) {
	mq.Lock()
	defer mq.Unlock()

	mq.push(false, aid, item)
}

func (mq *MessageQueue) Insert(aid int, item interface{}) {

}

func (mq *MessageQueue) Get() (item interface{}) {
	mq.Lock()
	defer mq.Unlock()

	return mq.Q.Head()
}

func (mq *MessageQueue) ReleaseIn(aid int) (ok bool) {
	mq.Lock()
	mq.Ack.Lock()
	defer mq.Ack.Unlock()
	defer mq.Unlock()

	return mq.release(true, aid)
}

func (mq *MessageQueue) ReleaseOut(aid int) (ok bool) {
	mq.Lock()
	mq.Ack.Lock()
	defer mq.Ack.Unlock()
	defer mq.Unlock()

	return mq.release(false, aid)
}

func (mq *MessageQueue) release(isInbound bool, aid int) (ok bool) {
	if mq.Ack.hasAck(isInbound, aid) {
		_ = mq.Ack.removeAck(isInbound, aid)
		_ = mq.pop()
		mq.Last = mq.Q.Head()
		return true
	}
	return false
}

func (mq *MessageQueue) push(isInbound bool, aid int, item interface{}) (ok bool) {
	if mq.Ack.hasAck(isInbound, aid) {
		return false
	}
	mq.Q.Push(item)

	return true
}

func (mq *MessageQueue) pop() (item interface{}) {
	return mq.Q.Pop()
}

// - MARK: MessageQueueAck section.

func (mqa *MessageQueueAck) HasInAck(aid int) bool {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.hasAck(true, aid)
}

func (mqa *MessageQueueAck) CreateInAck(aid int) (ok bool) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.createAck(true, aid)
}

func (mqa *MessageQueueAck) GetInAck(aid int) (ack *sync.Cond) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.getAck(true, aid)
}

func (mqa *MessageQueueAck) HasOutAck(aid int) bool {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.hasAck(false, aid)
}

func (mqa *MessageQueueAck) CreateOutAck(aid int) (ok bool) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.createAck(false, aid)
}

func (mqa *MessageQueueAck) GetOutAck(aid int) (ack *sync.Cond) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.getAck(false, aid)
}

func (mqa *MessageQueueAck) RemoveInAck(aid int) (ok bool) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.removeAck(true, aid)
}

func (mqa *MessageQueueAck) RemoveOutAck(aid int) (ok bool) {
	mqa.Lock()
	defer mqa.Unlock()

	return mqa.removeAck(false, aid)
}

func (mqa *MessageQueueAck) createAck(isInbound bool, aid int) (ok bool) {
	if isInbound {
		_, ok = mqa.in[aid]
	} else {
		_, ok = mqa.out[aid]
	}
	if !ok {
		if isInbound {
			mqa.in[aid] = sync.NewCond(&sync.Mutex{})
		} else {
			mqa.out[aid] = sync.NewCond(&sync.Mutex{})
		}
		return true
	}

	return false
}

func (mqa *MessageQueueAck) hasAck(isInbound bool, aid int) (ok bool) {
	if isInbound {
		_, ok = mqa.in[aid]
	} else {
		_, ok = mqa.out[aid]
	}

	return ok
}

func (mqa *MessageQueueAck) getAck(isInbound bool, aid int) (ack *sync.Cond) {
	var ok bool
	if isInbound {
		ack, ok = mqa.in[aid]
	} else {
		ack, ok = mqa.out[aid]
	}
	if !ok {
		return nil
	}

	return ack
}

func (mqa *MessageQueueAck) removeAck(isInbound bool, aid int) (ok bool) {
	if isInbound {
		_, ok = mqa.in[aid]
	} else {
		_, ok = mqa.out[aid]
	}
	if ok {
		if isInbound {
			delete(mqa.in, aid)
		} else {
			delete(mqa.out, aid)
		}
		return true
	}

	return false
}
