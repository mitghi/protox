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

package protocol

import (
	"bytes"
	"testing"
	"github.com/mitghi/protox/protocol/packet"
)

func TestQueue(t *testing.T) {
	var (
		q   *Queue = NewQueue()
		nq  *Queue = NewQueue()
		p   *packet.Packet
		err error
	)
	// setup queue
	q.Action = QAInitialize
	q.Address = "simple/queue/path"
	q.ReturnPath = "simple/queue/results"
	q.Mark = []byte("arbitary associated data")
	q.Meta.MessageId = 1
	// encode/serialize
	err = q.Encode()
	if err != nil {
		t.Fatal("expected err==nil")
	}
	// decode from serialized
	// data.
	p = q.GetPacket().(*packet.Packet)
	err = nq.DecodeFrom(p.Data)
	if err != nil {
		t.Fatal("expected err==nil")
	}
	// test deep equality
	if nq.Action != q.Action {
		t.Fatalf("inconsistent state, assertion failed (nq==%b, q==%b).", nq.Action, q.Action)
	}
	if bytes.Compare(nq.Message, q.Message) != 0 {
		t.Fatal("inconsistent state, assertion failed.")
	}
	if bytes.Compare(nq.Mark, q.Mark) != 0 {
		t.Fatal("inconsistent state, assertion failed.", nq.Mark, q.Mark)
	}
	if nq.Meta.MessageId != q.Meta.MessageId {
		t.Fatal("inconsistent state, assertion failed.", nq.Meta.MessageId, q.Meta.MessageId)
	}
	if nq.Address != q.Address {
		t.Fatal("inconsistent state, assertion failed.")
	}
	if nq.ReturnPath != q.ReturnPath {
		t.Fatal("inconsistent state, assertion failed.")
	}
}
