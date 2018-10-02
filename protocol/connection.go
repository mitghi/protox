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
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/timer"
)

// ensure interface (protocol) conformance.
var _ protobase.ProtoConnection = (*Connection)(nil)

// ProtocolConnection is the main interface for connections. Any struct that
// implement this interface can be passed to server.
type ProtocolConnection interface {
	receive() (result *[]byte, code byte, length int, err error)
	Receive() (packet *Packet, err error)
	ReceiveWithTimeout() (pcket *Packet, err error)
	Send(packet *Packet)
	SendDirect(packet *Packet)
	send(packet *Packet) (err error)
	HandleError(err error)
	Shutdown()
	Handle()
	AllocateChannels()
}

type Connection struct {
	protocon

	ErrorHandler       func(client *protobase.ClientInterface)
	server             protobase.ServerInterface                              // server delegate
	auth               protobase.AuthInterface                                // auth subsystem
	storage            protobase.MessageStorage                               // storage subsystem
	client             protobase.ClientInterface                              // client subsystem
	clientDelegate     func(string, string, string) protobase.ClientInterface // client delegate method
	permissionDelegate func(protobase.AuthInterface, ...string) bool          // permission delegate
	connTimeout        int                                                    // connection timeout (initial)
	heartbeat          int                                                    // maximum idle time
	unclean            uint32                                                 // status flag
	deadline           *timer.Timer                                           // ping timeout
	justStarted        bool                                                   // status flag
	State              ConnectionState                                        // connection state
	// TODO
	// . add PermissionInterface
	// . check alignment
}

// NewConnection returns a pointer to a new `Connection` struct and starts `Genseis`
// logic.
func NewConnection(conn net.Conn) *Connection {
	// TODO
	// . remove hardcoded values and add methods
	//   to set `connTimeout` and `heartbeat`.
	var result *Connection = &Connection{
		protocon: protocon{
			Conn:            conn,
			Reader:          bufio.NewReader(conn),
			Writer:          bufio.NewWriter(conn),
			corous:          sync.WaitGroup{},
			ErrChan:         nil,
			ShouldTerminate: nil,
			SendChan:        nil,
			Status:          STATDISCONNECT,
		},
		//initials
		justStarted:        true, // initial connection timeout is on
		connTimeout:        1,    // 1 second
		heartbeat:          1,    // 1 second ping interval
		ErrorHandler:       nil,
		State:              nil,
		client:             nil,
		server:             nil,
		auth:               nil,
		clientDelegate:     nil,
		permissionDelegate: nil,
	}
	result.State = NewGenesis(result)

	return result
}

// ErrorHandler is the signature for setting a error handler for `ClientInterface`.
type ErrorHandler func(client *protobase.ClientInterface)

// MarkClean atomically flips `unclean` flag.
func (self *Connection) MarkClean() {
	atomic.StoreUint32(&self.unclean, 0)
}

// MarkUnClean atomically flips `unclean` flag.
func (self *Connection) MarkUnClean() {
	atomic.StoreUint32(&self.unclean, 1)
}

// IsClean atomically checks if connection is clean.
func (self *Connection) IsClean() bool {
	return atomic.LoadUint32(&self.unclean) == 0
}

// SetInitiateTimeout is a receiver method that sets internal
// timeout value to `timeout` argument. This timeout is used to
// terminate connections that their initialization process takes
// longer than `timeout`.
func (self *Connection) SetInitiateTimeout(timeout int) {
	self.connTimeout = timeout
}

// SetHeartBeat sets the internal timer to maximum idle time which not receiving from a
// client will not results in connection termination. It is the responsibility of the client
// to keep up and send `Ping` packets to restart this timer.
func (self *Connection) SetHeartBeat(heartbeat int) {
	self.heartbeat = heartbeat
}

// SetClient sets the internal active client to `cl` argument. It is used for callbacks and
// notifications.
func (self *Connection) SetClient(cl protobase.ClientInterface) {
	self.client = cl
	self.State.SetClient(cl)
}

// SetMessageStorage sets internal message store to `storage`.
func (self *Connection) SetMessageStorage(storage protobase.MessageStorage) {
	self.storage = storage
}

// SetServer saves the server memory address for callbacks.
func (self *Connection) SetServer(sn protobase.ServerInterface) {
	self.server = sn
	self.State.SetServer(sn)
}

// GetClient returns the responsible struct implementing `ClientInterface`.
func (self *Connection) GetClient() protobase.ClientInterface {
	return self.client
}

// SetAuthenticator sets the delegate for authenication system. Argument `auth`
// will be used to authorize each client before passing into new stages.
func (self *Connection) SetAuthenticator(auth protobase.AuthInterface) {
	self.auth = auth
}

// SetClientDelegate is client delegate.
func (self *Connection) SetClientDelegate(cl func(string, string, string) protobase.ClientInterface) {
	self.clientDelegate = cl
}

// GetConnection returns underlaying `net.Conn` struct.
func (self *Connection) GetConnection() net.Conn {
	return self.Conn
}

// SetAuthTimeout is the initial deadline for authorization. If a client is not able
// to send its information in this period, it results in termination.
func (self *Connection) SetAuthTimeout(timeout int) {
	self.connTimeout = timeout
}

// SetPermissionDelegate sets permission callback which is used to ensure if certain
// operation is valid given restriction and rights of a user.
func (self *Connection) SetPermissionDelegate(pd func(protobase.AuthInterface, ...string) bool) {
	self.permissionDelegate = pd
}

// SendMessage creates a new packet and sends it to send channel.
func (self *Connection) SendMessage(pb protobase.MsgInterface, isOwner bool) {
	const fn = "SendMessage"
	var (
		envelope protobase.MsgEnvelopeInterface = pb.Envelope()
		topic    string                         = envelope.Route()
		message  []byte                         = envelope.Payload()
		cl       protobase.ClientInterface      = self.GetClient()
		clid     string                         = cl.GetIdentifier()
		msg      *Publish                       = NewPublish()
		qos      byte                           = pb.QoS()
		puid     uuid.UUID
	)
	msg.Message = message
	msg.Topic = topic
	/* d e b u g */
	// if isOwner {
	// 	logger.FDebug(fn, "* [SendMessage/IsOwner] redirecting to self(%s), MessageId(%d).", clid, pb.MessageId())
	// 	msg.Meta.MessageId = pb.MessageId()
	// 	msg.Meta.Qos = qos
	// 	msg.Encode()
	// 	data := msg.Encoded.Bytes()
	// 	packet := NewPacket(&data, msg.Command, msg.Encoded.Len())
	// 	self.Send(packet)
	// 	return
	// }
	/* d e b u g */
	if qos > 0 {
		logger.FDebug(fn, "* [QoS] QoS>0 in [SendMessage].", "qos", qos)
		puid = (*msg.Id)
		logger.FDebug("SendMessage", "Publish QoS.", qos, "msgdir", pb.Dir())
		if !self.storage.AddOutbound(clid, msg) {
			logger.Warn("- [MessageStore] unable to add outbound message in [SendMessage].")
		}
		idstore := self.storage.GetIDStoreO(clid)
		msg.Meta.MessageId = idstore.GetNewID(puid)
		msg.Meta.Qos = qos
		logger.FDebugf("SendMessage", "* [MessageId] id(%d). ", msg.Meta.MessageId)
	}
	msg.Encode()
	data := msg.Encoded.Bytes()
	packet := NewPacket(&data, msg.Command, msg.Encoded.Len())
	self.Send(packet)
}

func (self *Connection) SendRedelivery(pb protobase.EDProtocol) {
	var (
		p      *Packet  = pb.GetPacket().(*Packet)
		msg    *Publish = NewPublish()
		packet *Packet
	)
	/* d e b u g */
	// 	switch pb.(type) {
	// case *Publish:
	// 	msg = pb.(*Publish)
	// 	msg.Meta.Dup = true
	// 	msg.Encoded = nil
	// 	msg.Encode()
	// }
	// logger.FDebug("SendRedelivery", "* [Redelivery] sending stored packages to client.")
	// packet = pb.GetPacket().(*Packet)
	// self.Send(packet)
	/* d e b u g */
	if err := msg.DecodeFrom(p.Data); err != nil {
		logger.FDebug("sendRedelivery", "- [Redelivery] cannot decode a publish packet.")
	}
	msg.Meta.Dup = true
	logger.FDebugf("SendRedelivery", "* [Redelivery] sending stored packages to client  QoS(%d) Duplicate(%t).", msg.Meta.Qos, msg.Meta.Dup)
	msg.Encode()
	packet = msg.GetPacket().(*Packet)
	self.Send(packet)
}

// SetErrorHandler sets the delegate for `ClientInterface` error handler.
func (self *Connection) SetErrorHandler(fn func(client *protobase.ClientInterface)) {
	self.ErrorHandler = fn
}

// GetAuthenticator returns the actual auth subsystem used by the `Connection`.
func (self *Connection) GetAuthenticator() protobase.AuthInterface {
	return self.auth
}

// Handle is the entry routine into `Connection`. It is the main loop
// for handling initial logics/allocating and passing data to different stages.
func (self *Connection) Handle() {
	var (
		dur   time.Duration = time.Second * time.Duration(self.heartbeat)
		pchan chan *Packet
	)
	if self.justStarted {
		atomic.StoreUint32(&self.Status, STATCONNECTING)
		// - MARK: Wait for Genesis packet.
		pchan = make(chan *Packet, 1)
		packet, err := self.ReceiveWithTimeout(time.Second*time.Duration(self.connTimeout), pchan)
		close(pchan)
		// - MARK: End
		if err != nil {
			// NOTE: new
			// self.handleSendError(err)
			self.Conn.Close()
			self.SetStatus(protobase.STATDISCONNECT)
			self.MarkUnClean()
			self.server.NotifyReject(self)
			return
		}
		ok := self.State.HandleDefault(packet)
		if !ok {
			// TODO
			// self.handleSendError(err)
			// NOTE: new
			self.Conn.Close()
			self.SetStatus(protobase.STATDISCONNECT)
			self.MarkUnClean()
			self.server.NotifyReject(self)
			return
		} else {
			// TODO
			self.AllocateChannels()
		}
	}
	// add client to message store
	clid := self.client.GetIdentifier()
	if !self.storage.Exists(clid) {
		// TODO
		self.storage.AddClient(clid)
	}
	// mark that connect packet is received
	self.justStarted = false
	self.corous.Add(3)
	go self.prioSendHandler()
	go self.sendHandler()
	go self.recvHandler()
	self.client.Connected(nil)
	atomic.StoreUint32(&self.Status, STATONLINE)
	self.deadline = timer.NewTimer(dur)
	// TODO
	// . maybe let notifyconnected run in a seperate coroutine ????
	self.server.NotifyConnected(self)
	// main loop
ML:
	for {
		stat := atomic.LoadUint32(&self.Status)
		switch stat {
		case STATERR:
			logger.Debug("- [Error] STATERR.")
			break ML
		case STATGODOWN:
			logger.Debug("* [ForceShutdown] received force shutdown, cleaning up ....")
			break ML
		}
		select {
		case <-self.deadline.C:
			logger.Debug("- [Connection] DEADLINE, terminating ....", "userId", self.client.GetIdentifier())
			atomic.StoreUint32(&self.Status, protobase.STATERR)
			// TODO
			//  clean up
			break
		case packet := <-self.RecvChan:
			// TODO
			// . reuse session containers
			// . run this concurrently
			self.deadline.Reset(dur)
			logger.Debug("+ [Message] Received .", "userId", self.client.GetIdentifier(), "data", *packet.Data)
			self.dispatch(packet)
		case <-self.ErrChan:
			logger.Warn("- [Shit] went down. Panic.")
			break ML
		}
	}
	self.terminate()
	logger.Debug("* [MLHEAD]**(%s) beginning to wait for coroutines to finish**", self.client.GetIdentifier())
	atomic.StoreUint32(&self.Status, STATDISCONNECTED)
	self.server.NotifyDisconnected(self) // TODO: this should be done either in client or state
	self.corous.Wait()                   // Wait for all coroutines to finish before cleaning up
	atomic.StoreUint32(&self.Status, STATDISCONNECT)
	logger.Debugf("+ [MLEND]++(%s) all coroutines are finished, exiting conn++", self.client.GetIdentifier())
}

// ShutDown terminates the connection.
func (self *Connection) Shutdown() {
	self.Conn.Close()
	logger.Debug("* [Event] Shutting down stream.")
}

// Receive is a routine to read incomming packets from internal connection. It terminates
// a connection immediately if any abnormalities are detected or packet is malformed. This
// function blocks on purpose.
func (self *Connection) receive() (result *[]byte, code byte, length int, err error) {
	var (
		pack []byte
		rl   uint32
		cmd  byte //command byte from first byte of new packet
	)
	msg, err := self.Reader.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}
	// NOTE:
	// . 0xF0 is mask for command byte
	cmd = (msg & 0xF0) >> 4
	if self.precheck(&cmd) == false {
		return nil, 0, 0, InvalidCmdForState
	}
	pack = append(pack, msg)
	// Read remaining bytes after the fixed header
	err = ReadPacket(self.Reader, &pack, &rl)
	if err != nil {
		return nil, 0, 0, err
	}
	// NOTE: rl not included
	return &pack, cmd, 0, nil
}

// Precheck validate packets to avoid abnormalities and malformed packet.
// It returns a bool indicating wether a control packet is valid or not.
func (self *Connection) precheck(cmd *byte) bool {
	msg := *cmd

	if self.justStarted == true {
		// dirty check for connect packet
		// reject any other packet in this stage
		if msg == 0x01 {
			// self.gotFirstPacket = true <- is relocated
			return true
		}
		return false
	}
	return IsValidCommand(msg)
}

// SendHandler reads from send channel and writes to a socket.
func (self *Connection) sendHandler() {
	for packet := range self.SendChan {
		// Check QoS
		err := self.send(packet)
		if err != nil {
			// self.handleSendError(err)
			// return
			// TODO
			// . handle this case
			logger.Debug("- [SendHandler] ERROR, BEFORE BREAKING.")
			break
		}
	}
	self.corous.Done()
}

// prioSendHandler handles packages with high priority.
func (self *Connection) prioSendHandler() {
	const fname = "prioSendHandler"
	for packet := range self.PrioSendChan {
		err := self.send(packet)
		if err != nil {
			logger.FError(fname, "- [PrioSendHandler] ERROR, BEFORE BREAKING.")
			break
		}
	}
	self.corous.Done()
}

// recvHandler is the main receive handler.
func (self *Connection) recvHandler() {
	const fname = "recvHandler"
	for {
		packet, err := self.Receive()
		if err != nil {
			logger.FError(fname, "- [RecvHandler] error while receiving packets.")
			// TODO
			//  handle errors
			break
		}
		self.RecvChan <- packet

	}
	self.corous.Done()
}

// HandleSendError is a error handler. It is used for errors
// caused by sending packets. Currently it terminates the
// connection.
func (self *Connection) handleSendError(err error) {
	self.Conn.Close()
}

// terminate shuts down the connection and undoes all side effects.
func (self *Connection) terminate() {
	// TODO
	// . change this and check status by sync/atomic
	// already closed
	if self.SendChan == nil || self.PrioSendChan == nil {
		logger.FDebug("terminate", "? [Terminate] both [SendChan], [PrioSendChan] are (nil).")
		return
	}
	self.Shutdown()
	self.SendLock.Lock()
	if self.PrioSendChan != nil {
		close(self.PrioSendChan)
	}
	if self.SendChan != nil {
		close(self.SendChan)
	}
	self.PrioSendChan = nil
	self.SendChan = nil
	self.SendLock.Unlock()
}

// dispatch is responsible to call the correct methods on state structures.
func (self *Connection) dispatch(packet *Packet) {
	switch (*packet).Code {
	case PCONNECT:
		self.State.onCONNECT(packet)
	case PCONNACK:
		self.State.onCONNACK(packet)
	case PQUEUE:
		self.State.onQUEUE(packet)
	case PSUBSCRIBE:
		self.State.onSUBSCRIBE(packet)
	case PSUBACK:
		self.State.onSUBACK(packet)
	case PPUBLISH:
		self.State.onPUBLISH(packet)
	case PPUBACK:
		self.State.onPUBACK(packet)
	case PPING:
		self.State.onPING(packet)
	case PDISCONNECT:
		self.State.onDISCONNECT(packet)
	// NOTE: Rest of protocol data suite should be integrated in this case
	default:
		logger.FWarn("dispatch", "- [Dispacher] unrecognized cmd code in packet.", packet)
	}
}

// SetNetConnection is a receiver method that sets internal network connection
// and read/write buffers. It is used during initialization when connection
// is established or during a connection restart.
func (self *Connection) SetNetConnection(conn net.Conn) {
	self.Conn = conn
	self.Writer = bufio.NewWriter(conn)
	self.Reader = bufio.NewReader(conn)
}
