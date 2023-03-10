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

	"github.com/google/uuid"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
)

// Ensure protocol (interface) conformance.
var _ protobase.ProtoConnection = (*Connection)(nil)

// ErrorHandler is `ClientInterface` error handler signature.
type ErrorHandler func(client *protobase.ClientInterface)

// Section: structs

// Connection is high level manager acting as
// hub, connecting subsystems and utilizing them
// to form meaningful procedures and performing
// meaningful operations neccessary to have
// a block composing basic steps in bootstrapping
// protox.
type Connection struct {
	// TODO
	// . refactor permissionDelegate via PermissionInterface
	// . check alignment
	protocon

	ErrorHandler       func(client *protobase.ClientInterface)
	clientDelegate     func(string, string, string) protobase.ClientInterface // client delegate method
	permissionDelegate func(protobase.AuthInterface, ...string) bool          // permission delegate
	server             protobase.ServerInterface                              // server delegate
	auth               protobase.AuthInterface                                // auth subsystem
	storage            protobase.MessageStorage                               // storage subsystem
	client             protobase.ClientInterface                              // client subsystem
	State              protobase.ConnectionState                              // connection state
	deadline           *time.Ticker                                           // ping timeout
	connTimeout        int                                                    // connection timeout (initial)
	heartbeat          int                                                    // maximum idle time
	unclean            uint32                                                 // refurbished flag
	justStarted        bool                                                   // status flag
}

// Section: initializers [ constructors ]

// NewConnection allocates and initializes a new
// `Connection` struct from supplied `net.Conn`
// without performing health checks and returns
// its pointer. It sets the state to `Genseis`
// procedure responsible for handling new
// network connections.
func NewConnection(conn net.Conn) *Connection {
	var (
		c *Connection = &Connection{
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
			justStarted:        true, // deadline is connTimeout when true
			connTimeout:        CConnectionDefaultTimeout,
			heartbeat:          CConnectionDefaultHeartbeat,
			ErrorHandler:       nil,
			State:              nil,
			client:             nil,
			server:             nil,
			auth:               nil,
			clientDelegate:     nil,
			permissionDelegate: nil,
		}
	)
	// set connection driver
	c.State = NewGenesis(c)
	return c
}

// Section: Connection methods

// MarkClean atomically flips `unclean` flag.
// Indicates struct is reusable.
func (c *Connection) MarkClean() {
	atomic.StoreUint32(&c.unclean, 0)
}

// MarkUnClean atomically flips `unclean` flag.
// Indicates struct is not reusable.
func (c *Connection) MarkUnClean() {
	atomic.StoreUint32(&c.unclean, 1)
}

// IsClean atomically checks if connection is clean.
func (c *Connection) IsClean() bool {
	return atomic.LoadUint32(&c.unclean) == 0
}

// Section: Connection Setter-Getter methods

// GetClient returns client field.
func (c *Connection) GetClient() protobase.ClientInterface {
	return c.client
}

// SetInitiateTimeout sets initial timeout.
func (c *Connection) SetInitiateTimeout(timeout int) {
	c.connTimeout = timeout
}

// SetHeartBeat sets grace threshold between
// since last 'Ping'.
func (c *Connection) SetHeartBeat(heartbeat int) {
	c.heartbeat = heartbeat
}

// SetClient sets client struct.
func (c *Connection) SetClient(cl protobase.ClientInterface) {
	c.client = cl
	c.State.SetClient(cl)
}

// SetMessageStorage sets message storage.
func (c *Connection) SetMessageStorage(storage protobase.MessageStorage) {
	c.storage = storage
}

// SetServer sets internal pointer to given server instance.
func (c *Connection) SetServer(sn protobase.ServerInterface) {
	c.server = sn
	c.State.SetServer(sn)
}

// SetAuthenticator sets authentication delegate.
func (c *Connection) SetAuthenticator(auth protobase.AuthInterface) {
	c.auth = auth
}

// SetClientDelegate sets client delegate.
func (c *Connection) SetClientDelegate(cl func(string, string, string) protobase.ClientInterface) {
	c.clientDelegate = cl
}

// GetConnection returns underlaying `net.Conn` struct.
func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

// SetAuthTimeout sets initial authorization deadline.
// Exceeding specified timeout result in quick
// termination.
func (c *Connection) SetAuthTimeout(timeout int) {
	c.connTimeout = timeout
}

// SetPermissionDelegate sets the callback used
// to ensure validity of access given restriction
// and rights of a user.
func (c *Connection) SetPermissionDelegate(pd func(protobase.AuthInterface, ...string) bool) {
	c.permissionDelegate = pd
}

// SetErrorHandler sets the delegate for `protobase.ClientInterface` handler.
func (c *Connection) SetErrorHandler(fn func(client *protobase.ClientInterface)) {
	c.ErrorHandler = fn
}

// GetAuthenticator returns underlying auth subsystem.
func (c *Connection) GetAuthenticator() protobase.AuthInterface {
	return c.auth
}

// Section: Connection main-methods

// HandleDefault is non-blocking entry point; only
// for 'Genesis' stage.
// NOTE: empty method
func (c *Connection) HandleDefault(apckt protobase.PacketInterface) {
	// NOP
}

// SendMessage is a method for sending
// message to remote destination.
func (c *Connection) SendMessage(pb protobase.MsgInterface, isOwner bool) (err error) {
	const fn string = "SendMessage"
	var (
		envelope protobase.MsgEnvelopeInterface = pb.Envelope()
		topic    string                         = envelope.Route()
		message  []byte                         = envelope.Payload()
		cl       protobase.ClientInterface      = c.GetClient()
		clid     string                         = cl.GetIdentifier()
		msg      *Publish                       = NewRawPublish()
		qos      byte                           = pb.QoS()
		puid     uuid.UUID

		data []byte
		p    *Packet
	)
	msg.Message = message
	msg.Topic = topic
	if qos > 0 {
		logger.FDebug(fn, "* [QoS] QoS>0 in [SendMessage].", "qos", qos)
		puid = (msg.Id)
		logger.FDebug("SendMessage", "Publish QoS.", qos, "msgdir", pb.Dir())
		if !c.storage.AddOutbound(clid, msg) {
			logger.Warn("- [MessageStore] unable to add outbound message in [SendMessage].")
		}
		idstore := c.storage.GetIDStoreO(clid)
		msg.Meta.MessageId = idstore.GetNewID(puid)
		msg.Meta.Qos = qos
		logger.FDebugf("SendMessage", "* [MessageId] id(%d). ", msg.Meta.MessageId)
	}
	err = msg.Encode()
	if err != nil {
		logger.FWarnf(fn, "- [Connection] unable to encode publish packet. error:", err)
		// TODO
		// . handle errors by changing the execution flow
		return err
	}
	data = msg.Encoded.Bytes()
	p = NewPacket(data, msg.Command, msg.Encoded.Len())
	c.Send(p)
	return nil
}

func (c *Connection) SendRedelivery(pb protobase.EDProtocol) (err error) {
	const fn string = "SendRedelivery"
	var (
		p      *Packet  = pb.GetPacket().(*Packet)
		msg    *Publish = NewPublish(p)
		packet *Packet
	)
	if msg == nil {
		logger.FDebug("sendRedelivery", "- [Redelivery] cannot decode a publish packet.", pb)
	}
	// if err := msg.DecodeFrom(p.Data); err != nil {
	// 	logger.FDebug("sendRedelivery", "- [Redelivery] cannot decode a publish packet.")
	// }
	msg.Meta.Dup = true
	logger.FDebugf("SendRedelivery", "* [Redelivery] sending stored packages to client  QoS(%d) Duplicate(%t).", msg.Meta.Qos, msg.Meta.Dup)
	err = msg.Encode()
	if err != nil {
		logger.FWarnf(fn, "- [Connection] unable to encode publish packet. error:", err)
		return err
	}
	packet = msg.GetPacket().(*Packet)
	c.Send(packet)
	return nil
}

// Handle is the entry routine into `Connection`. It is the main loop
// for handling initial logics/allocating and passing data to different stages.
func (c *Connection) Handle() {
	var (
		dur   time.Duration = time.Second * time.Duration(c.heartbeat)
		pchan chan *Packet
		clid  string
		stat  uint32
	)
	if c.justStarted {
		atomic.StoreUint32(&c.Status, STATCONNECTING)
		// wait for packet starting 'Genesis' stage
		pchan = make(chan *Packet, 1)
		packet, err := c.ReceiveWithTimeout(time.Second*time.Duration(c.connTimeout), pchan)
		close(pchan)
		if err != nil {
			c.Conn.Close()
			c.SetStatus(protobase.STATDISCONNECT)
			c.MarkUnClean()
			// deadline exceeded
			c.server.NotifyReject(c)
			return
		}
		ok := c.State.HandleDefault(packet)
		if !ok {
			c.Conn.Close()
			c.SetStatus(protobase.STATDISCONNECT)
			c.MarkUnClean()
			// failed 'Genesis' stage
			c.server.NotifyReject(c)
			return
		} else {
			c.AllocateChannels()
		}
	}
	clid = c.client.GetIdentifier()
	if !c.storage.Exists(clid) {
		c.storage.AddClient(clid) // add client to message store
	}
	c.justStarted = false // swap fresh start flag
	// run I/O coroutines
	c.corous.Add(3)
	go c.prioSendHandler()
	go c.sendHandler()
	go c.recvHandler()
	c.client.Connected(nil)          // TODO: pass execution to background thread; in case of blocking call.
	c.SetStatus(STATONLINE)          // connection is established
	c.deadline = time.NewTicker(dur) // ping interval
	c.server.NotifyConnected(c)      // perform blocking call ( ensure serial execution )
	// main loop
ML:
	for {
		stat = c.GetStatus()
		switch stat {
		case STATERR:
			logger.Infof("- [Error][Connection] status is error flag, exiting main loop.")
			break ML
		case STATGODOWN:
			logger.Infof("* [ForceShutdown][Connection] received force shutdown, terminating main loop.")
			break ML
		}
		select {
		case <-c.deadline.C:
			logger.Debug("- [Connection] connection exceeded the deadline. Termination in progress.",
				"userId", clid)
			c.SetStatus(protobase.STATERR)
			break
		case packet := <-c.RecvChan:
			c.deadline.Reset(dur)
			logger.Debug("+ [Message][Connection] received message on receive channel.",
				"userId", clid, "data", packet.Data)
			c.dispatch(packet)
		case <-c.ErrChan:
			logger.Warn("- [Error][Connection] error channel contains error. Termination in progress.")
			break ML
		}
	}
	c.terminate()
	logger.Debugf("* [Connection][Termination]**(%s) waiting for jobs to finish.", clid)
	c.SetStatus(STATDISCONNECTED)  // connection terminated
	c.server.NotifyDisconnected(c) // TODO: this should be done either in client or state
	c.corous.Wait()                // Wait for all coroutines to finish before cleaning up
	c.SetStatus(STATDISCONNECT)    // connection remains inactive
	logger.Debugf("+ [Connection][Termination]++(%s) all jobs are done. Exiting connection handler.", clid)
}

// ShutDown terminates the connection.
func (c *Connection) Shutdown() {
	// TODO
	// . handle error
	c.Conn.Close()
	logger.Debug("* [Event][Termination] shutting down connection.")
}

// Receive is a routine to read incomming packets from internal connection. It terminates
// a connection immediately if any abnormalities are detected or packet is malformed. This
// function blocks on purpose.
func (c *Connection) receive() (result []byte, code byte, length int, err error) {
	var (
		pack []byte
		rl   uint32
		cmd  byte //command byte from first byte of new packet
	)
	msg, err := c.Reader.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}
	// NOTE:
	// . 0xF0 is mask for command byte
	cmd = (msg & 0xF0) >> 4
	if c.precheck(&cmd) == false {
		return nil, 0, 0, protocol.InvalidCmdForState
	}
	pack = append(pack, msg)
	// Read remaining bytes after the fixed header
	err = protocol.ReadPacket(c.Reader, &pack, &rl)
	if err != nil {
		return nil, 0, 0, err
	}
	// NOTE: rl not included
	return pack, cmd, 0, nil
}

// Precheck validate packets to avoid abnormalities and malformed packet.
// It returns a bool indicating wether a control packet is valid or not.
func (c *Connection) precheck(cmd *byte) bool {
	var (
		msg byte = *cmd
	)
	if c.justStarted == true {
		// dirty check for connect packet
		// reject any other packet in this stage
		if msg == 0x01 {
			return true
		}
		return false
	}
	return IsValidCommand(msg)
}

// SendHandler reads from send channel and writes to a socket.
func (c *Connection) sendHandler() {
	// TODO
	// . check QoS
	// . error handler function
	const fn string = "sendHandler"
	for packet := range c.SendChan {
		err := c.send(packet)
		if err != nil {
			logger.FDebug(fn, "- [Connection] unable to send packet. error:", err)
			break
		}
	}
	logger.FInfo(fn, "+ [Connection] signaling done to work group.")
	c.corous.Done()
}

// prioSendHandler handles packages with high priority.
func (c *Connection) prioSendHandler() {
	const fn string = "prioSendHandler"
	for packet := range c.PrioSendChan {
		err := c.send(packet)
		if err != nil {
			logger.FError(fn, "- [PrioSendHandler] unable to send priority packet. error:", err)
			break
		}
	}
	logger.FInfo(fn, "+ [Connection] signaling done to work group.")
	c.corous.Done()
}

// recvHandler is the main receive handler.
func (c *Connection) recvHandler() {
	// TODO
	// . error handler function
	const fn string = "recvHandler"
	for {
		packet, err := c.Receive()
		if err != nil {
			logger.FError(fn, "- [RecvHandler] error while receiving packets. error:", err)
			break
		}
		c.RecvChan <- packet

	}
	logger.FInfo(fn, "+ [Connection] signaling done to work group.")
	c.corous.Done()
}

// HandleSendError handles 'send(...)' error.
func (c *Connection) handleSendError(err error) {
	c.Conn.Close()
}

// terminate shuts down the connection undoing side effects.
func (c *Connection) terminate() {
	// TODO
	// . change this and check status by sync/atomic
	//   already closed.
	if c.SendChan == nil || c.PrioSendChan == nil {
		logger.FDebug("terminate", "? [Terminate] both [SendChan], [PrioSendChan] are (nil).")
		return
	}
	c.Shutdown()
	/* critical section */
	c.SendLock.Lock()
	if c.PrioSendChan != nil {
		close(c.PrioSendChan)
	}
	if c.SendChan != nil {
		close(c.SendChan)
	}
	c.PrioSendChan = nil
	c.SendChan = nil
	c.SendLock.Unlock()
	/* critical section - end */
}

// dispatch dispatches the packet to its responsible handler.
func (c *Connection) dispatch(packet protobase.PacketInterface) {
	// TODO:
	// . add handler for remaining PDUs
	// . dispatch via lookup table
	const fn string = "dispatch"
	switch packet.GetCode() {
	case protobase.PCONNECT:
		c.State.OnCONNECT(packet)
	case protobase.PCONNACK:
		c.State.OnCONNACK(packet)
	case protobase.PQUEUE:
		c.State.OnQUEUE(packet)
	case protobase.PSUBSCRIBE:
		c.State.OnSUBSCRIBE(packet)
	case protobase.PSUBACK:
		c.State.OnSUBACK(packet)
	case protobase.PPUBLISH:
		c.State.OnPUBLISH(packet)
	case protobase.PPUBACK:
		c.State.OnPUBACK(packet)
	case protobase.PPING:
		c.State.OnPING(packet)
	case protobase.PDISCONNECT:
		c.State.OnDISCONNECT(packet)
	default:
		logger.FWarn(fn, "- [Dispacher] unrecognized cmd code in packet.", packet)
	}
}

// SetNetConnection is a receiver method that sets internal network connection
// and read/write buffers. It is used during initialization when connection
// is established or during a connection restart.
func (c *Connection) SetNetConnection(conn net.Conn) {
	c.Conn = conn
	c.Writer = bufio.NewWriter(conn)
	c.Reader = bufio.NewReader(conn)
}
