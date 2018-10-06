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

package networking

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mitghi/protox/protocol"
	"github.com/mitghi/protox/protocol/packet"
)

// protocon is the data flow layer embedding
// network fd, low level struct implementing
// and providign I/O, compatible to protox.
type protocon struct {
	// TODO
	// . check padding
	sync.RWMutex

	Conn   net.Conn // connection section
	Reader *bufio.Reader
	Writer *bufio.Writer // end

	corous          sync.WaitGroup // spawned coroutines
	SendLock        sync.Mutex
	cendch          chan struct{}
	ShouldTerminate chan struct{}
	ErrChan         chan struct{}
	SendChan        chan *packet.Packet
	PrioSendChan    chan *packet.Packet // Priority send channel
	RecvChan        chan *packet.Packet
	addr            string // ip address
	Status          uint32 // status flag
	// TODO
	// . implement connection state ( reuse this struct. Prevent new allocations. )
	// . set the external error handler ( non-critical errors )
	// State           ConnectionState // end
	// ErrorHandler func(client *protobase.ClientInterface)
}

// AllocateChannels initializes internal
// send/receive channels. Their creation
// are deferred to reduce unneccessary
// memory allocations up until all initial
// stages including critical checks are passed.
func (pc *protocon) AllocateChannels() {
	pc.SendChan = make(chan *Packet, 1024)
	pc.PrioSendChan = make(chan *Packet, 1024)
	pc.RecvChan = make(chan *Packet, 1024)
	pc.ErrChan = make(chan struct{}, 1)
	/* d e b u g */
	// pc.cendch = make(chan struct{})
	// pc.ShouldTerminate = make(chan struct{})
	/* d e b u g */
}

// Receive is a helper function which creates
// a new packet from incomming data. Result
// should be later checked to create/cast
// into a packet. When succesfull, it returns
// the formed packet; nil with error when
// unsuccessfull. NOTE: Blocking routine.
func (pc *protocon) Receive() (p *Packet, err error) {
	var (
		pck    *[]byte // incoming data; don't pass slice header
		cmd    byte    // packet type flag
		length int     // length of remaining data in buffer
	)
	pck, cmd, length, err = pc.receive()
	if err != nil {
		return nil, err
	}
	p = packet.NewPacket(pck, cmd, length)
	return p, nil
}

// ReceiveWithTimeout waits `timeout` seconds before
// returning an erorr. It polls responsible coroutine
// `timeout` seconds and issues timeout with no data
// and error set to 'CriticalTimeout'.
func (pc *protocon) ReceiveWithTimeout(timeout time.Duration, inbox chan *Packet) (pcket *Packet, err error) {
	// TODO
	// . investigate timer channel
	//   for deadlock.
	var (
		period <-chan time.Time = time.After(timeout) // timeout channel
	)
	// increment wait group
	// and spawn the job
	pc.corous.Add(1)
	go func(pchan chan<- *Packet, wg *sync.WaitGroup) {
		packet, err := pc.Receive()
		// NOTE: should not close the channel here
		if err != nil {
		} else {
			pchan <- packet
		}
		// Signal workgroups that we are done
		wg.Done()
	}(inbox, &pc.corous)
	select {
	case packet, ok := <-inbox:
		/* d e b u g */
		// accelerate garbage collection
		// pchan = nil
		// period = nil
		/* d e b u g */
		if !ok {
			return nil, protocol.BadMsgTypeError
		}
		return packet, nil
	case <-period:
		/* d e b u g */
		// pchan = nil
		/* d e b u g */
		logger.Debug("[PeriodError] Timeout occured")
		return nil, protocol.CriticalTimeout
	}
}

// Send routine is responsible to write a packet into send channel. Actuall sending
// is done by a coroutine.
func (pc *protocon) Send(packet *Packet) {
	pc.SendLock.Lock()
	if pc.SendChan != nil {
		// NOTE: IMPORTANT: NEW
		pc.SendChan <- packet
		// select {
		// case pc.SendChan <- packet:
		// default:
		// }

	} else {
		logger.Debug("[NOTICE]: SendChannel is nil")
	}
	pc.SendLock.Unlock()
}

// SendPrio is the s end hadnler for packets with higher priority.
func (pc *protocon) SendPrio(packet *Packet) {
	pc.SendLock.Lock()
	if pc.PrioSendChan != nil {
		// NOTE: IMPORTANT: NEW:
		// select {
		// case pc.PrioSendChan <- packet:
		// default:
		// 	logger.Debug("* [NOTICE]: PrioSend failed")
		// }
		pc.PrioSendChan <- packet
	} else {
		logger.Debug("* [NOTICE]: PrioSendChannel is nil")
	}
	pc.SendLock.Unlock()
}

// SendDirect directly writes a packet to a socket.
func (pc *protocon) SendDirect(packet *Packet) {
	pc.send(packet)
}

// GetStatus atomically returns the current status.
func (pc *protocon) GetStatus() (stat uint32) {
	stat = atomic.LoadUint32(&pc.Status)
	return stat
}

// GetErrChan returns a chan for errors.
func (pc *protocon) GetErrChan() chan struct{} {
	return pc.ErrChan
}

// SetStatus atomically sets the current status. It will be evaluated in the
// main loop.
func (pc *protocon) SetStatus(status uint32) {
	atomic.StoreUint32(&pc.Status, status)
}

// receive reads the protocol command flag
// and strips fixed header. It fills the
// the buffer according to packet header.
// Error indicates broken I/O pipe or
// protocol violation ( malformed data ,
// incorrect header , ... ).
func (pc *protocon) receive() (result *[]byte, code byte, length int, err error) {
	const fn string = "receive"
	var (
		pack []byte
		rl   uint32 // remaining length
		cmd  byte   // command flag (first higher byte of new packet)
	)
	msg, err := pc.Reader.ReadByte()
	if err != nil {
		logger.FDebug(fn, "readpacket(readbyte)error, msg, err", msg, err)
		return nil, 0, 0, err
	}
	// NOTE: 0xF0 is mask for command byte
	cmd = (msg & 0xF0) >> 4
	/* d e b u g */
	// check if command flag is valid,
	// early return when invalid
	// if pc.precheck(&cmd) == false {
	// 	return nil, 0, 0, InvalidCmdForState
	// }
	/* d e b u g */
	pack = append(pack, msg)
	// Read remaining bytes after the fixed header
	err = protocol.ReadPacket(pc.Reader, &pack, &rl)
	if err != nil {
		logger.FDebug(fn, "readpacket(receive)error, msg, err", msg, err, rl)
		logger.FDebug(fn, "readpacket pc.reader:", pc.Reader, pc.Reader.Size())
		return nil, 0, 0, err
	}
	// NOTE: rl not included
	return &pack, cmd, 0, nil
}

// send writes packet data to underlying connection.
func (pc *protocon) send(packet *Packet) (err error) {
	const fn string = "send"
	pc.Lock()
	defer pc.Unlock()
	if pc.Writer == nil {
		logger.FWarn(fn, "- [send] no buffer(writer) for writing data.")
		logger.FDebug(fn, "- [send] protocon.Writer==nil.")
		return protocol.EINVLWRTBFR
	}
	// TODO
	// . write remaining data
	_, err = pc.Writer.Write(*packet.Data)
	if err != nil {
		return err
	}
	err = pc.Writer.Flush()
	if err != nil {
		logger.FDebug(fn, "unable to flush data for transportation.", err)
		return err
	}
	return nil
}

// uniSendHandler is the transport coroutine
// responsible for sending data to its remote
// destination.
func (pc *protocon) uniSendHandler() {
	// TODO
	// . handle broken channels
	// . write test case
	const fname string = "uniSendHandler"
	defer func() {
		logger.FDebug(fname, "* [UniSendHandler] is stopped. Signaling done to work group.")
		pc.corous.Done()
	}()
	for {
		select {
		case packet, ok := <-pc.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] PrioSendChan is closed.", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendChan] error while sending priority packets via [PrioSendChan].", "error:", err)
				return
			}
		case packet, ok := <-pc.SendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] SendChan is closed.", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [UniSendHandler] error while sending packets.", "error:", err)
				return
			}
		case <-pc.cendch:
			logger.FDebug(fname, "* [UniSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

// SendHandler is the transport coroutine
// responsible for writing outgoing data to
// send channel.
func (pc *protocon) sendHandler() {
	// TODO
	// . propagate error and call
	//   responsible error handler.
	const fname string = "sendHandler"
	defer func() {
		logger.FInfo(fname, "is stopped. Signaling done to work group.")
		pc.corous.Done()
	}()
	for {
		select {
		case packet, ok := <-pc.SendChan:
			if !ok {
				logger.FDebug(fname, "- [SendChan] is closed.", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [SendHandler] error while sending packets.", "error:", err)
				return
			}
		case <-pc.cendch:
			logger.FDebug(fname, "* [SendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

// prioSendHandler is the transport coroutine
// responsible for writing high priority
// outgoing data to priority send channel.
func (pc *protocon) prioSendHandler() {
	const fname string = "prioSendHandler"
	defer func() {
		logger.FInfo(fname, "is stopped. Signaling done to work group.")
		pc.corous.Done()
	}()
	for {
		select {
		case packet, ok := <-pc.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "PrioSendChan is closed", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets", "error:", err)
				return
			}
		case <-pc.cendch:
			logger.FDebug(fname, "* [PrioSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

// HandleSendError is a error handler. It is used for errors
// caused by sending packets. Currently it terminates the
// connection.
func (pc *protocon) handleSendError(err error) {
	pc.Conn.Close()
}

// - - - - - - STASH - - - - - - -

/* d e b u g */
// // recvHandler is the main receive handler.
// func (pc *protocon) recvHandler() {
// 	const fname string = "recvHandler"
// 	defer func() {
// 		logger.FInfo(fname, "is stopped. Signaling done to work group.")
// 		pc.corous.Done()
// 	}()
// 	for {
// 		packet, err := pc.Receive()
// 		if err != nil {
// 			logger.FError(fname, "- [RecvHandler] error while receiving packets.")
// 			// TODO
// 			//  handle errors
// 			return
// 		}
// 		pc.RecvChan <- packet
// 	}
// 	// pc.corous.Done()
// }
/* d e b u g */
