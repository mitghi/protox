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
	"fmt"
	"testing"

	"github.com/mitghi/protox/protobase"
)

type dummyPacket struct {
	protobase.EDProtocol

	name string
}

const (
	eFATALfmt string = "inconsistent state, expected %s, got %v ( %s %v )."
)

func TestRetain(t *testing.T) {
	var (
		topics [][]byte = [][]byte{
			[]byte("a/test/topic"),
			[]byte("one/level"),
			[]byte("a/restricted/multi/level/topic"),
			[]byte("self/general"),
		}
		retain *Retain      = NewRetain()
		packet *dummyPacket = &dummyPacket{name: "test"}
		err    error
	)
	// insertion
	err = retain.rinsert(topics[0], packet)
	if err != nil {
		t.Fatalf("assertion failed, expected err==nil, got %v.", err)
	}
	// existing case
	np, err := retain.rfind(topics[0])
	if err != nil {
		t.Fatalf(eFATALfmt, "err==nil", err, "np", np)
	}
	// check if its null pointer
	if np == nil {
		t.Fatalf(eFATALfmt, "np!=nil", np, "", "")
	}
	// typecast, because `dummyPacket` doesn't
	// conform to `protobase.EDProtocol` in a
	// valid way.
	dp, ok := np.(*dummyPacket)
	if !ok {
		t.Fatalf(eFATALfmt, "ok==true", ok, "dp", dp)
	}
	if dp.name != packet.name {
		t.Fatalf(eFATALfmt, "dp.name==packet.name", dp, "packet.name:", packet.name)
	}
	// non-existing case
	_, err = retain.rfind(topics[1])
	if err == nil {
		t.Fatalf(eFATALfmt, "err!=nil", err, "", "")
	}
	// insert all topics
	for i, v := range topics {
		err = retain.rinsert(v, &dummyPacket{
			name: fmt.Sprintf("%s%d", "test", i),
		})
		if err != nil {
			t.Fatalf("assertion failed, expected err==nil, got %v.", err)
		}
	}
	np, err = retain.rfind(topics[0])
	if err != nil {
		t.Fatalf(eFATALfmt, "err==nil", err, "np", np)
	}
	if np == nil {
		t.Fatalf(eFATALfmt, "np!=nil", np, "", "")
	}
	dp, ok = np.(*dummyPacket)
	if !ok {
		t.Fatalf(eFATALfmt, "ok==true", ok, "dp", dp)
	}
	if dp.name != "test0" {
		t.Fatalf(eFATALfmt, "dp.name==packet.name", dp, "packet.name:", packet.name)
	}
	// remove a node
	err = retain.rremove(topics[0])
	if err != nil {
		t.Fatalf(eFATALfmt, "err==nil", err, "", "")
	}
	// ensure it's removal
	_, err = retain.rfind(topics[0])
	if err == nil {
		t.Fatalf(eFATALfmt, "err!=nil", err, "", "")
	}
	// check another topic
	_, err = retain.rfind(topics[1])
	if err != nil {
		t.Fatalf(eFATALfmt, "err==nil", err, "", "")
	}
}
