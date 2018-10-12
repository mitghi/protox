package networking

import (
	"sync"

	"github.com/mitghi/protox/protobase"
)

type CLBCallback func(protobase.OptionInterface, protobase.MsgInterface)

type CLBPacketInterface interface {
}

type CLBPacket struct {
	*sync.RWMutex

	clbpub map[uint16]CLBCallback
	clbsub map[uint16]CLBCallback
}

func NewCLBPacket() *CLBPacket {
	return &CLBPacket{
		RWMutex: &sync.RWMutex{},
		clbpub:  make(map[uint16]CLBCallback),
		clbsub:  make(map[uint16]CLBCallback),
	}
}

// func (clbp *CLBPacket) insert(id uint16, callback CLBCallback, table ) (ok bool) {
//   _, ok := clbp.
// }

// func (clbp *CLBPacket) insert(id uint16, callback CLBCallback) (ok bool) {

// }

// func (clbp *CLBPacket)
