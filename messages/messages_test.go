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
	"sort"
	"testing"

	"github.com/google/uuid"

	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
)

var (
	ECLNEX = "client %s does not exists"
	ECLIN  = "error inserting client %s into msgstore"
	EINVS  = "invalid state"
	EADD   = "error adding packet to msgstore for client %s ( %s )"
)

const (
	DEFCLN = "test"
)

func init() {
	var _ protobase.MessageStorage = (*MessageStore)(nil)
}

func newMStore() *MessageStore {
	var result *MessageStore = NewMessageStore()
	result.Init()

	return result
}

func TestAddClient(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()

	store.AddClient(DEFCLN)
	if _, ok := store.in[DEFCLN]; !ok {
		t.Fatalf(ECLIN, DEFCLN)
	}

	store.Close(DEFCLN)
}

func TestClose(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()
	store.AddClient(DEFCLN)
	if ok := store.Close(DEFCLN); !ok {
		t.Fatal(EINVS)
	}
	if ok := store.Close("TEST"); ok {
		t.Fatal(EINVS)
	}
}

func TestInsertInbound(t *testing.T) {
	var store *MessageStore = newMStore()
	var pckt *protocol.Ping

	store.AddClient(DEFCLN)
	pckt = protocol.NewRawPing()
	if ok := store.AddInbound(DEFCLN, pckt); !ok {
		t.Fatalf(EADD, DEFCLN, " INBOUND ")
	}
	if ok := store.AddInbound("nonexisting-user", pckt); ok {
		t.Fatal(EINVS)
	}
	// reinsert an existing packet - should return false
	if ok := store.AddInbound(DEFCLN, pckt); ok {
		t.Fatalf(EADD, DEFCLN, " INBOUND ")
	}

	store.Close(DEFCLN)
}

func TestInsertOutbound(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()
	var pckt *protocol.Pong

	store.AddClient(DEFCLN)
	pckt = protocol.NewRawPong()
	if ok := store.AddOutbound(DEFCLN, pckt); !ok {
		t.Fatalf(EADD, DEFCLN, " OUTBOUND ")
	}
	if ok := store.AddOutbound("TEST", pckt); ok {
		t.Fatal(EINVS)
	}
	if ok := store.AddOutbound(DEFCLN, pckt); ok {
		t.Fatal(EINVS)
	}
	if ok := store.AddOutbound(DEFCLN, pckt); ok {
		t.Fatalf(EADD, DEFCLN, " OUTBOUND ")
	}

	store.Close(DEFCLN)
}

func TestGetAllOut(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()
	var clnames map[string][]string = make(map[string][]string)
	var fmtclname string = "test_%d"

	// add dummies
	for i := 0; i < 5; i++ {
		var clname string = fmt.Sprintf(fmtclname, i)
		store.AddClient(clname)
		clnames[clname] = make([]string, 0, 5)
		for j := 0; j < 5; j++ {
			var pckt *protocol.Pong = protocol.NewRawPong()
			var uid uuid.UUID = pckt.UUID()
			var sid string = uid.String()
			clnames[clname] = append(clnames[clname], sid)
			store.AddOutbound(clname, pckt)
		}
		sort.Strings(clnames[clname])
	}
	for k, _ := range clnames {
		if ok := store.Exists(k); !ok {
			t.Fatalf(EINVS)
		}
		var v []string = clnames[k]
		var values []string = store.GetAllOutStr(k)
		var vlen int = len(values)
		var clen int = len(v)
		sort.Strings(values)
		if vlen != clen {
			t.Fatalf(EINVS)
		}
		for i := 0; i < vlen; i++ {
			if values[i] != v[i] {
				t.Fatalf(EINVS)
			}
		}
	}
	// GetAllOut length
	if msgs := store.GetAllOut("test_0"); len(msgs) != len(clnames["test_0"]) {
		t.Fatal(EINVS)
	}
	// non existing user
	if msgs := store.GetAllOutStr("NON_EXISTING_USER"); len(msgs) != 0 {
		t.Fatal(EINVS)
	}
	// non existing user - GetAllOut
	if msgs := store.GetAllOut("NON_EXISTING_USER"); len(msgs) != 0 {
		t.Fatal(EINVS)
	}
}

func TestDeleteIn(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()
	var values []protobase.EDProtocol = make([]protobase.EDProtocol, 0, 5)

	store.AddClient(DEFCLN)
	for i := 0; i < 5; i++ {
		var pckt *protocol.Pong = protocol.NewRawPong()
		values = append(values, pckt)
		store.AddInbound(DEFCLN, pckt)
	}
	for i := 0; i < 5; i++ {
		var pckt *protocol.Pong = values[i].(*protocol.Pong)
		if ok := store.DeleteIn(DEFCLN, pckt); !ok {
			t.Fatal(EINVS)
		}
	}
	// non existing user
	if ok := store.DeleteIn("NON_EXISTING_CLIENT", values[0].(*protocol.Pong)); ok {
		t.Fatal(EINVS)
	}
	// existing user, non existing packet
	var pckt *protocol.Pong = protocol.NewRawPong()
	if ok := store.DeleteIn(DEFCLN, pckt); ok {
		t.Fatal(EINVS)
	}
}

func TestDeleteOut(t *testing.T) {
	var store *MessageStore = NewInitedMessageStore()
	var values []protobase.EDProtocol = make([]protobase.EDProtocol, 0, 5)

	store.AddClient(DEFCLN)
	for i := 0; i < 5; i++ {
		var pckt *protocol.Pong = protocol.NewRawPong()
		values = append(values, pckt)
		store.AddOutbound(DEFCLN, pckt)
	}
	for i := 0; i < 5; i++ {
		var pckt *protocol.Pong = values[i].(*protocol.Pong)
		if ok := store.DeleteOut(DEFCLN, pckt); !ok {
			t.Fatal(EINVS)
		}
	}
	// non existing user
	if ok := store.DeleteOut("NON_EXISTING_CLIENT", values[0].(*protocol.Pong)); ok {
		t.Fatal(EINVS)
	}
	// existing user, non existing packet
	var pckt *protocol.Pong = protocol.NewRawPong()
	if ok := store.DeleteOut(DEFCLN, pckt); ok {
		t.Fatal(EINVS)
	}
}

// TestGetNewID covers most of `MessageId` methods except a single case
// in `GetNewID(uuid.UUID)`. Full test case is excluded to a new file
// because of its long running time (`TestGetNewIDThreaded(t *testing.T)`).
func TestGetNewID(t *testing.T) {
	var msgid *MessageId = NewMessageId()
	var ids [6]uint16
	var uids [6]uuid.UUID
	var i uint16

	for i = 0; i < 6; i++ {
		var pckt *protocol.Pong = protocol.NewRawPong()
		var uid uuid.UUID = pckt.UUID()
		uids[i] = uid
		ids[i] = msgid.GetNewID(uid)
	}
	for i = 0; i < 6; i++ {
		if iocc := msgid.IsOccupied(i + 1); !iocc {
			t.Fatalf("Supposed to be occupied")
		}
	}
	for i = 0; i < 6; i++ {
		if ids[i] != i+1 {
			t.Fatalf(EINVS)
		}
	}
	for i = 0; i < 6; i++ {
		var curruid uuid.UUID = uids[i]
		if uid := msgid.GetNewID(curruid); uid == i+1 || uid == i {
			t.Fatal(EINVS)
		}
	}
	for i = 0; i < 6; i++ {
		msgid.FreeId(i + 1)
	}
	for i = 0; i < 6; i++ {
		if iocc := msgid.IsOccupied(i + 1); iocc {
			t.Fatalf("Supposed to be occupied")
		}
	}
}

func TestOrderedMessages(t *testing.T) {
	var msgstore *MessageStore = NewInitedMessageStore()
	msgstore.AddClient("test")
	for i := 0; i < 6; i++ {
		var pckt *protocol.Publish = protocol.NewRawPublish()
		pckt.Message = []byte(fmt.Sprintf("test_%d", i))
		if err := pckt.Encode(); err != nil {
			t.Fatal(EINVS)
		}

		msgstore.AddOutbound("test", pckt)
	}

	all := msgstore.GetAllOut("test")
	for i, msg := range all {
		expstr := fmt.Sprintf("test_%d", i)
		p := protocol.NewPublish(msg.GetPacket())
		if p == nil {
			t.Fatal(EINVS)
		}
		if string(p.Message) != expstr {
			t.Fatal(EINVS)
		}
	}
}
