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
	"testing"
)

func TestMQAck(t *testing.T) {
	var (
		mq  *MessageQueue = NewMessageQueue()
		ack *sync.Cond
	)
	if !mq.Ack.CreateInAck(0) {
		t.Fatal("assertion failed: expected true, got false.")
	}
	if !mq.Ack.CreateInAck(2) {
		t.Fatal("assertion failed: expected true, got false.")
	}
	if ack = mq.Ack.GetInAck(0); ack == nil {
		t.Fatalf("assertion failed: expected ack!=nil, got %v.", ack)
	}
	if ack = mq.Ack.GetInAck(1); ack != nil {
		t.Fatalf("inconsistent state: expected ack==nil, got %v.", ack)
	}
	if !mq.Ack.RemoveInAck(0) {
		t.Fatal("inconsistent state: expected true, got 'false'.")
	}
	if mq.Ack.RemoveInAck(0) {
		t.Fatal("inconsistent state: expected false, got 'true'. Attempt to remove unexisting value.")
	}
	if !mq.ReleaseIn(2) {
		t.Fatal("inconsistent state: expected true, got 'false'. Unable to release existing item.")
	}
	if mq.ReleaseIn(2) {
		t.Fatal("inconsistent state: expected false, got 'true'. Attempt to release unexisting item.")
	}
}
