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
)

// Section: structs

// protocon is the data flow layer embedding
// network fd, low level struct implementing
// and providign I/O, compatible to protox.
type protocon struct {
	// TODO
	// . check padding
	sync.RWMutex

	Conn            net.Conn
	Reader          *bufio.Reader
	Writer          *bufio.Writer
	corous          sync.WaitGroup
	SendLock        sync.Mutex
	cendch          chan struct{}
	ShouldTerminate chan struct{}
	ErrChan         chan struct{}
	SendChan        chan *Packet
	PrioSendChan    chan *Packet
	RecvChan        chan *Packet
	addr            string
	Status          uint32
	// TODO
	// . implement connection state ( reuse this struct. Prevent new allocations. )
	// . set the external error handler ( non-critical errors )
	// State           ConnectionState // end
	// ErrorHandler func(client *protobase.ClientInterface)
}

// Section: protocon methods

// AllocateChannels allocates chan
// fields.
func (pc *protocon) AllocateChannels() {
	pc.SendChan = make(chan *Packet, 1024)
	pc.PrioSendChan = make(chan *Packet, 1024)
	pc.RecvChan = make(chan *Packet, 1024)
	pc.ErrChan = make(chan struct{}, 1)
}

// Receive reads available data from the
// buffer, blocks when no data is available.
// It returns '*Packet' containing protocol
// header and its payload.
func (pc *protocon) Receive() (p *Packet, err error) {
	var (
		pck    []byte // incoming data; don't pass slice header
		cmd    byte   // packet type flag
		length int    // length of remaining data in buffer
	)
	pck, cmd, length, err = pc.receive()
	if err != nil {
		return nil, err
	}
	p = NewPacket(pck, cmd, length)
	return p, nil
}

// ReceiveWithTimeout uses `timeout` grace duration
// waiting for data from poller coroutine before
// returning 'CriticalTimeout' error without data
// or returning immediately meanwhile when data
// becomes available.
func (pc *protocon) ReceiveWithTimeout(timeout time.Duration, inbox chan *Packet) (pcket *Packet, err error) {
	// TODO: investigate timer channel inconsistencies.
	var (
		period <-chan time.Time = time.After(timeout) // timeout channel
	)
	// increment wait group
	// and spawn the job
	// caller is responsible to close
	// 'inbox' chan
	pc.corous.Add(1)
	go func(pchan chan<- *Packet, wg *sync.WaitGroup) {
		packet, err := pc.Receive()
		if err != nil {
		} else {
			pchan <- packet
		}
		wg.Done() // Signal workgroups that we are done
	}(inbox, &pc.corous)
	select {
	case packet, ok := <-inbox:
		if !ok {
			return nil, protocol.BadMsgTypeError
		}
		return packet, nil
	case <-period:
		logger.Debug("- [PeriodError][protocon] timeout occured.")
		return nil, protocol.CriticalTimeout
	}
}

// Send routine is responsible to write a packet into send channel. Actuall sending
// is done by a coroutine.
func (pc *protocon) Send(packet *Packet) {
	const fn string = "Send"
	pc.SendLock.Lock()
	if pc.SendChan != nil {
		pc.SendChan <- packet
	} else {
		logger.Debug(fn, "- [NOTICE][protocon]: send channel is nil.")
	}
	pc.SendLock.Unlock()
}

// SendPrio is the s end hadnler for packets with higher priority.
func (pc *protocon) SendPrio(packet *Packet) {
	pc.SendLock.Lock()
	defer pc.SendLock.Unlock()
	if pc.PrioSendChan != nil {
		// NOTE: IMPORTANT: NEW:
		// select {
		// case pc.PrioSendChan <- packet:
		// default:

		// }
		if _, ok := (<-pc.PrioSendChan); ok {
			pc.PrioSendChan <- packet
		} else {
			logger.Debug("- [NOTICE][protocon] PrioSend is closed.")
		}
	} else {
		logger.Debug("- [NOTICE][protocon] PrioSendChannel is nil.")
	}
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

// SetStatus atomically sets the current status. It will be evaluated in the
// main loop.
func (pc *protocon) SetStatus(status uint32) {
	atomic.StoreUint32(&pc.Status, status)
}

// GetErrChan returns a chan for errors.
func (pc *protocon) GetErrChan() chan struct{} {
	return pc.ErrChan
}

// Section: protocon private-methods

// receive reads the protocol command flag
// and strips fixed header. It fills
// the buffer according to packet header.
// Error indicates broken I/O pipe or
// protocol violation ( malformed data ,
// incorrect header , ... ).
func (pc *protocon) receive() (result []byte, code byte, length int, err error) {
	const fn string = "receive"
	var (
		pack []byte
		rl   uint32 // remaining length
		cmd  byte   // command flag
	)
	msg, err := pc.Reader.ReadByte()
	if err != nil {
		logger.FDebugf(fn, "- [protocon] readpacket(readbyte)error, msg:", msg, ", err:", err)
		return nil, 0, 0, err
	}
	// NOTE: 0xF0 is mask for command byte
	cmd = (msg & 0xF0) >> 4
	pack = append(pack, msg)
	// Read remaining bytes after the fixed header
	err = protocol.ReadPacket(pc.Reader, &pack, &rl)
	if err != nil {
		logger.FDebugf(fn, "- [protocon] readpacket(receive)error, msg, err", msg, err, rl)
		logger.FDebugf(fn, "- [protocon] readpacket pc.reader:", pc.Reader, pc.Reader.Size())
		return nil, 0, 0, err
	}
	// NOTE: rl not included
	return pack, cmd, 0, nil
}

// send writes packet data to underlying connection.
func (pc *protocon) send(packet *Packet) (err error) {
	const fn string = "send"
	pc.Lock()
	defer pc.Unlock()
	if pc.Writer == nil {
		logger.FWarn(fn, "- [protocon] no buffer(writer) for writing data.")
		logger.FDebug(fn, "- [protocon] protocon.Writer==nil.")
		return protocol.EINVLWRTBFR
	}
	// TODO
	// . write remaining data
	_, err = pc.Writer.Write(packet.Data)
	if err != nil {
		return err
	}
	err = pc.Writer.Flush()
	if err != nil {
		logger.FDebug(fn, "- [protocon] unable to flush data for transportation.", err)
		return err
	}
	return nil
}

// uniSendHandler is the transport coroutine
// responsible for sending data to its remote
// destination.
func (pc *protocon) uniSendHandler() {
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
				logger.FDebug(fname, "- [SendChan][protocon] is closed.", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [SendHandler][protocon] error while sending packets.", "error:", err)
				return
			}
		case <-pc.cendch:
			logger.FDebug(fname, "* [SendHandler][protocon] received end signal from [cendch] channel, terminating coroutine.")
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
				logger.FDebug(fname, "- [PrioSendChan][protocon] is closed", "ok status:", ok)
				return
			}
			if err := pc.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler][protocon] error while sending priority packets", "error:", err)
				return
			}
		case <-pc.cendch:
			logger.FDebug(fname, "* [PrioSendHandler][protocon] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

// HandleSendError handles 'send(...)' error.
func (pc *protocon) handleSendError(err error) {
	// MARK: remove this method?
	const fn string = "handleSendError"
	var (
		e error
	)
	e = pc.Conn.Close()
	if e != nil {
		logger.FInfo(fn, "* [protocon] unable to close the connection. error:", err)
	}
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
