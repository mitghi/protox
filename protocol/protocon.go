package protocol

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"fmt"
)

// protocon is the struct for handling protox compatible
// connections.
type protocon struct {
	sync.RWMutex

	Conn            net.Conn // connection section
	Reader          *bufio.Reader
	Writer          *bufio.Writer  // end
	corous          sync.WaitGroup // main section
	SendLock        sync.Mutex
	cendch          chan struct{}
	ShouldTerminate chan struct{}
	ErrChan         chan struct{}
	SendChan        chan *Packet
	PrioSendChan    chan *Packet // Priority send channel
	RecvChan        chan *Packet
	addr            string
	Status          uint32
	// State           ConnectionState // end
	// ErrorHandler func(client *protobase.ClientInterface)
}

// AllocateChannels initializes internal send/receive channels. Their creation
// are deferred to reduce unneccessary memory allocations up until all initial
// stages including critical checks are passed.
func (self *protocon) AllocateChannels() {
	self.SendChan = make(chan *Packet, 1024)
	self.PrioSendChan = make(chan *Packet, 1024)
	self.RecvChan = make(chan *Packet, 1024)
	self.ErrChan = make(chan struct{}, 1)
	// self.cendch = make(chan struct{})
	// self.ShouldTerminate = make(chan struct{})
}

// Receive is a helper function which creates a new packet from incomming data.
// Result should be checked later to create/cast to a particular packet.
func (self *protocon) Receive() (packet *Packet, err error) {
	pck, cmd, length, err := self.receive()
	if err != nil {
		return nil, err
	}
	resultPacket := NewPacket(pck, cmd, length)

	return resultPacket, nil
}

// ReceiveWithTimeout waits `timeout` seconds before returning an erorr. It polls a coroutine
// for `timeout` seconds.
func (self *protocon) ReceiveWithTimeout(timeout time.Duration, inbox chan *Packet) (pcket *Packet, err error) {
	period := time.After(timeout)
	self.corous.Add(1)

	go func(pchan chan<- *Packet, wg *sync.WaitGroup) {
		packet, err := self.Receive()
		// NOTE: should not  close the channel here
		if err != nil {
		} else {
			pchan <- packet
		}
		// Signal workgroups that we are done
		wg.Done()
	}(inbox, &self.corous)
	select {
	case packet, ok := <-inbox:
		// pchan = nil
		// period = nil
		if !ok {
			return nil, BadMsgTypeError
		}
		return packet, nil
	case <-period:
		// pchan = nil
		logger.Debug("[PeriodError] Timeout occured")
		return nil, CriticalTimeout
	}
}

// Send routine is responsible to write a packet into send channel. Actuall sending
// is done by a coroutine.
func (self *protocon) Send(packet *Packet) {
	self.SendLock.Lock()
	if self.SendChan != nil {
		// NOTE: IMPORTANT: NEW
		self.SendChan <- packet
		// select {
		// case self.SendChan <- packet:
		// default:
		// }

	} else {
		logger.Debug("[NOTICE]: SendChannel is nil")
	}
	self.SendLock.Unlock()
}

// SendPrio is the s end hadnler for packets with higher priority.
func (self *protocon) SendPrio(packet *Packet) {
	self.SendLock.Lock()
	if self.PrioSendChan != nil {
		// NOTE: IMPORTANT: NEW:
		// select {
		// case self.PrioSendChan <- packet:
		// default:
		// 	logger.Debug("* [NOTICE]: PrioSend failed")
		// }
		self.PrioSendChan <- packet
	} else {
		logger.Debug("* [NOTICE]: PrioSendChannel is nil")
	}
	self.SendLock.Unlock()
}

// SendDirect directly writes a packet to a socket.
func (self *protocon) SendDirect(packet *Packet) {
	self.send(packet)
}

// GetStatus atomically returns the current status.
func (self *protocon) GetStatus() (stat uint32) {
	stat = atomic.LoadUint32(&self.Status)
	return stat
}

// GetErrChan returns a chan for errors.
func (self *protocon) GetErrChan() chan struct{} {
	return self.ErrChan
}

// SetStatus atomically sets the current status. It will be evaluated in the
// main loop.
func (self *protocon) SetStatus(status uint32) {
	atomic.StoreUint32(&self.Status, status)
}

func (self *protocon) receive() (result *[]byte, code byte, length int, err error) {
	var (
		pack []byte
		rl   uint32
		//command byte from first byte of new packet
		cmd byte
	)
	msg, err := self.Reader.ReadByte()
	if err != nil {
    fmt.Println("readpacket(readbyte)error, msg, err", msg, err)
		return nil, 0, 0, err
	}
	// NOTE: 0xF0 is mask for command byte
	cmd = (msg & 0xF0) >> 4
	// if self.precheck(&cmd) == false {
	// 	return nil, 0, 0, InvalidCmdForState
	// }
	pack = append(pack, msg)
	// Read remaining bytes after the fixed header
	err = ReadPacket(self.Reader, &pack, &rl)
	if err != nil {
    fmt.Println("readpacket(receive)error, msg, err", msg, err, rl)
    fmt.Println("readpacket self.reader:", self.Reader, self.Reader.Size())
		return nil, 0, 0, err
	}
	// NOTE: rl not included
	return &pack, cmd, 0, nil
}

// Send function writes a packet to a socket.
func (self *protocon) send(packet *Packet) (err error) {
	const fn = "send"
	self.Lock()
	defer self.Unlock()
	if self.Writer == nil {
		logger.FWarn(fn, "protocon.Writer==nil")
		return errors.New("protocol: attempt using null writer")
	}
	_, err = self.Writer.Write(*packet.Data)
	if err != nil {
		return err
	}
	self.Writer.Flush()
	// NOTE: NEW: UNTESTED:
	// err = ...
	// self.Writer.Flush()
	// if err != nil {
	// 	return err
	// }
	return nil
}

// uniSendHandler is the unified send handler which handles all
// scenarios for sending a package.
func (self *protocon) uniSendHandler() {
	const fname = "uniSendHandler"
	defer func() {
		logger.FDebug(fname, "before decrementing workgroup")
		self.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-self.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] PrioSendChan is closed.", "ok status:", ok)
				return
			}
			if err := self.send(packet); err != nil {
				logger.FError(fname, "- [UniSendHandler] error while sending packets to [PrioSendChan].", "error:", err)
				return
			}
		case packet, ok := <-self.SendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] SendChan is closed.", "ok status:", ok)
				return
			}
			if err := self.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets.", "error:", err)
				return
			}
		case <-self.cendch:
			logger.FDebug(fname, "* [UniSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

// SendHandler reads from send channel and writes to a socket.
func (self *protocon) sendHandler() {
	const fname = "sendHandler"
	defer func() {
		logger.FInfo(fname, "before decrementing workgroup")
		self.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-self.SendChan:
			if !ok {
				logger.FDebug(fname, "PrioSendChan is closed", "ok status:", ok)
				return
			}
			if err := self.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets", "error:", err)
				return
			}
		case <-self.cendch:
			logger.FDebug(fname, "* [PrioSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
	/* d e b u g */
	// for packet := range self.SendChan {
	// 	// Check QoS
	// 	err := self.send(packet)
	// 	if err != nil {
	// 		// self.handleSendError(err)
	// 		// return
	// 		// TODO
	// 		// Handle this case
	// 		logger.Debug("- [ERROR IN SENDHANDLER, BEFORE BREAKING]")
	// 		return
	// 	}
	// }
	// self.corous.Done()
	/* d e b u g */
}

//
func (self *protocon) prioSendHandler() {
	const fname = "prioSendHandler"
	defer func() {
		logger.FInfo(fname, "before decrementing workgroup")
		self.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-self.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "PrioSendChan is closed", "ok status:", ok)
				return
			}
			if err := self.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets", "error:", err)
				return
			}
		case <-self.cendch:
			logger.FDebug(fname, "* [PrioSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
	/* d e b u g */
	// for packet := range self.PrioSendChan {
	// 	err := self.send(packet)
	// 	if err != nil {

	// 	}
	// }
	// self.corous.Done()
	/* d e b u g */
}

// recvHandler is the main receive handler.
func (self *protocon) recvHandler() {
	const fname = "recvHandler"
	defer func() {
		logger.FInfo(fname, "before decrementing workgroup")
		self.corous.Done()
	}()
	for {
		packet, err := self.Receive()
		if err != nil {
			logger.FError(fname, "- [RecvHandler] error while receiving packets.")
			// TODO
			//  handle errors
			return
		}
		self.RecvChan <- packet
	}
	// self.corous.Done()
}

// HandleSendError is a error handler. It is used for errors
// caused by sending packets. Currently it terminates the
// connection.
func (self *protocon) handleSendError(err error) {
	self.Conn.Close()
}
