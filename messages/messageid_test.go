//+build msgidfull

package messages

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/mitghi/protox/protocol"
)

func TestGetNewIDThreaded(t *testing.T) {
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
	var curruid uuid.UUID = uids[0]
	fmt.Println("First stage ....")
	for j := 60000; j < MSGMAXLEN; j++ {
		_ = msgid.GetNewID(curruid)
	}
	fmt.Println("+First stage done")
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for j := 1; j <= 60; j++ {
		wg.Add(1)
		go func(c int, msgid *MessageId, cuid *uuid.UUID, wg *sync.WaitGroup) {
			var (
				u int = 1000
				e int = c * u
				s int = e - u
			)
			fmt.Printf("coroutine %d started.\n", c)
			for i := s; i < e; i++ {
				_ = msgid.GetNewID(*cuid)
			}
			fmt.Printf("+Couroutine %d finished.\n", c)
			wg.Done()
		}(j, msgid, &curruid, wg)
	}

	fmt.Println("Waiting for coroutines to finish ....")
	wg.Wait()

	if nuid := msgid.GetNewID(curruid); nuid != 0 {
		t.Fatal(EINVS)
	}
}
