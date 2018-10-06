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
	"github.com/mitghi/timer"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
	"github.com/mitghi/protox/protocol/packet"
)

// ensure interface (protocol) conformance.
var _ protobase.ProtoConnection = (*Connection)(nil)

// ErrorHandler is `ClientInterface` error handler signature.
type ErrorHandler func(client *protobase.ClientInterface)

// Connection is high level manager acting as 
// hub connecting subsystems and utilizing them
// to form meaningful procedures and performing
// meaningful operations manifesting itself as 
// basic protox system.
type Connection struct {
	// TODO
	// . refactor permissionDelegate via PermissionInterface
	// . check alignment  
	protocon

	ErrorHandler       func(client *protobase.ClientInterface)
	server             protobase.ServerInterface                              // server delegate
	auth               protobase.AuthInterface                                // auth subsystem
	storage            protobase.MessageStorage                               // storage subsystem
	client             protobase.ClientInterface                              // client subsystem
	clientDelegate     func(string, string, string) protobase.ClientInterface // client delegate method
	permissionDelegate func(protobase.AuthInterface, ...string) bool          // permission delegate
	deadline           *timer.Timer                                           // ping timeout
	State              protobase.ConnectionState                              // connection state
	connTimeout        int                                                    // connection timeout (initial)
	heartbeat          int                                                    // maximum idle time
	unclean            uint32                                                 // status flag
	justStarted        bool                                                   // status flag    
}

// NewConnection allocates and initializes a new 
// `Connection` struct from supplied `net.Conn`
// without performing health checks and returns 
// its pointer. It sets the state to `Genseis`
// procedure responsible for handling new
// network connections.
func NewConnection(conn net.Conn) *Connection {
	// TODO
	// . remove hardcoded values
  // . add methods to set `connTimeout` 
	//   and `heartbeat`.
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
  )
  // set the connection state
	c.State = NewGenesis(c)
	return c
}

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

// SetInitiateTimeout is a receiver method that sets internal
// timeout value to `timeout` argument. This value is threshold
// to terminate connections that their initialization process takes
// longer than `timeout`.
func (c *Connection) SetInitiateTimeout(timeout int) {
	c.connTimeout = timeout
}

// SetHeartBeat sets the internal timer to maximum idle time which not 
// receiving from a client will not results in connection termination. 
// It is the responsibility of the client to keep up and send `Ping`
// packets to restart this timer.
func (c *Connection) SetHeartBeat(heartbeat int) {
	c.heartbeat = heartbeat
}

// SetClient sets the internal active client to `cl` argument. It is 
// used for callbacks and notifications.
func (c *Connection) SetClient(cl protobase.ClientInterface) {
	c.client = cl
	c.State.SetClient(cl)
}

// SetMessageStorage sets internal message store to `storage`.
func (c *Connection) SetMessageStorage(storage protobase.MessageStorage) {
	c.storage = storage
}

// SetServer saves the server memory address for callbacks.
func (c *Connection) SetServer(sn protobase.ServerInterface) {
	c.server = sn
	c.State.SetServer(sn)
}

// GetClient returns the responsible struct implementing `ClientInterface`.
func (c *Connection) GetClient() protobase.ClientInterface {
	return c.client
}

// SetAuthenticator sets the delegate for authenication system. Argument `auth`
// will be used to authorize each client before passing into new stages.
func (c *Connection) SetAuthenticator(auth protobase.AuthInterface) {
	c.auth = auth
}

// SetClientDelegate is client delegate.
func (c *Connection) SetClientDelegate(cl func(string, string, string) protobase.ClientInterface) {
	c.clientDelegate = cl
}

// GetConnection returns underlaying `net.Conn` struct.
func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

// SetAuthTimeout is the initial deadline for authorization. If a client is not able
// to send its information in this period, it results in termination.
func (c *Connection) SetAuthTimeout(timeout int) {
	c.connTimeout = timeout
}

// SetPermissionDelegate sets permission callback which is used to ensure if certain
// operation is valid given restriction and rights of a user.
func (c *Connection) SetPermissionDelegate(pd func(protobase.AuthInterface, ...string) bool) {
	c.permissionDelegate = pd
}

// SendMessage creates a new packet and sends it to send channel.
func (c *Connection) SendMessage(pb protobase.MsgInterface, isOwner bool) {
	const fn string = "SendMessage"
	var (
		envelope protobase.MsgEnvelopeInterface = pb.Envelope()
		topic    string                         = envelope.Route()
		message  []byte                         = envelope.Payload()
		cl       protobase.ClientInterface      = c.GetClient()
		clid     string                         = cl.GetIdentifier()
		msg      *Publish                       = protocol.NewPublish()
		qos      byte                           = pb.QoS()
		puid     uuid.UUID
	)
	msg.Message = message
	msg.Topic = topic
	/* d e b u g */
	// if isOwner {
	// 	logger.FDebug(fn, "* [SendMessage/IsOwner] redirecting to c(%s), MessageId(%d).", clid, pb.MessageId())
	// 	msg.Meta.MessageId = pb.MessageId()
	// 	msg.Meta.Qos = qos
	// 	msg.Encode()
	// 	data := msg.Encoded.Bytes()
	// 	packet := NewPacket(&data, msg.Command, msg.Encoded.Len())
	// 	c.Send(packet)
	// 	return
	// }
	/* d e b u g */
	if qos > 0 {
		logger.FDebug(fn, "* [QoS] QoS>0 in [SendMessage].", "qos", qos)
		puid = (*msg.Id)
		logger.FDebug("SendMessage", "Publish QoS.", qos, "msgdir", pb.Dir())
		if !c.storage.AddOutbound(clid, msg) {
			logger.Warn("- [MessageStore] unable to add outbound message in [SendMessage].")
		}
		idstore := c.storage.GetIDStoreO(clid)
		msg.Meta.MessageId = idstore.GetNewID(puid)
		msg.Meta.Qos = qos
		logger.FDebugf("SendMessage", "* [MessageId] id(%d). ", msg.Meta.MessageId)
	}
	msg.Encode()
	data := msg.Encoded.Bytes()
	packet := packet.NewPacket(&data, msg.Command, msg.Encoded.Len())
	c.Send(packet)
}

func (c *Connection) SendRedelivery(pb protobase.EDProtocol) {
	var (
		p      *Packet  = pb.GetPacket().(*Packet)
		msg    *Publish = protocol.NewPublish()
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
	// c.Send(packet)
	/* d e b u g */
	if err := msg.DecodeFrom(p.Data); err != nil {
		logger.FDebug("sendRedelivery", "- [Redelivery] cannot decode a publish packet.")
	}
	msg.Meta.Dup = true
	logger.FDebugf("SendRedelivery", "* [Redelivery] sending stored packages to client  QoS(%d) Duplicate(%t).", msg.Meta.Qos, msg.Meta.Dup)
	msg.Encode()
	packet = msg.GetPacket().(*Packet)
	c.Send(packet)
}

// SetErrorHandler sets the delegate for `ClientInterface` error handler.
func (c *Connection) SetErrorHandler(fn func(client *protobase.ClientInterface)) {
	c.ErrorHandler = fn
}

// GetAuthenticator returns the actual auth subsystem used by the `Connection`.
func (c *Connection) GetAuthenticator() protobase.AuthInterface {
	return c.auth
}

// HandleDefault is to satisfy interface requirements. It is 
// stub routine.
func (c *Connection) HandleDefault(apckt protobase.PacketInterface) {
}

// Handle is the entry routine into `Connection`. It is the main loop
// for handling initial logics/allocating and passing data to different stages.
func (c *Connection) Handle() {
	var (
		dur   time.Duration = time.Second * time.Duration(c.heartbeat)
		pchan chan *Packet
	)
	if c.justStarted {
		atomic.StoreUint32(&c.Status, STATCONNECTING)
		// - MARK: Wait for Genesis packet.
		pchan = make(chan *Packet, 1)
		packet, err := c.ReceiveWithTimeout(time.Second*time.Duration(c.connTimeout), pchan)
		close(pchan)
		// - MARK: End
		if err != nil {
			// NOTE: new
			// c.handleSendError(err)
			c.Conn.Close()
			c.SetStatus(protobase.STATDISCONNECT)
			c.MarkUnClean()
			c.server.NotifyReject(c)
			return
		}
		ok := c.State.HandleDefault(packet)
		if !ok {
			// TODO
			// c.handleSendError(err)
			// NOTE: new
			c.Conn.Close()
			c.SetStatus(protobase.STATDISCONNECT)
			c.MarkUnClean()
			c.server.NotifyReject(c)
			return
		} else {
			// TODO
			c.AllocateChannels()
		}
	}
	// add client to message store
	clid := c.client.GetIdentifier()
	if !c.storage.Exists(clid) {
		// TODO
		c.storage.AddClient(clid)
	}
	// mark that connect packet is received
	c.justStarted = false
	c.corous.Add(3)
	go c.prioSendHandler()
	go c.sendHandler()
	go c.recvHandler()
	c.client.Connected(nil)
	atomic.StoreUint32(&c.Status, STATONLINE)
	c.deadline = timer.NewTimer(dur)
	// TODO
	// . maybe let notifyconnected run in a seperate coroutine ????
	c.server.NotifyConnected(c)
	// main loop
ML:
	for {
		stat := atomic.LoadUint32(&c.Status)
		switch stat {
		case STATERR:
			logger.Debug("- [Error] STATERR.")
			break ML
		case STATGODOWN:
			logger.Debug("* [ForceShutdown] received force shutdown, cleaning up ....")
			break ML
		}
		select {
		case <-c.deadline.C:
			logger.Debug("- [Connection] DEADLINE, terminating ....", "userId", c.client.GetIdentifier())
			atomic.StoreUint32(&c.Status, protobase.STATERR)
			// TODO
			//  clean up
			break
		case packet := <-c.RecvChan:
			// TODO
			// . reuse session containers
			// . run this concurrently
			c.deadline.Reset(dur)
			logger.Debug("+ [Message] Received .", "userId", c.client.GetIdentifier(), "data", *packet.Data)
			c.dispatch(packet)
		case <-c.ErrChan:
			logger.Warn("- [Shit] went down. Panic.")
			break ML
		}
	}
	c.terminate()
	logger.Debug("* [MLHEAD]**(%s) beginning to wait for coroutines to finish**", c.client.GetIdentifier())
	atomic.StoreUint32(&c.Status, STATDISCONNECTED)
	c.server.NotifyDisconnected(c) // TODO: this should be done either in client or state
	c.corous.Wait()                   // Wait for all coroutines to finish before cleaning up
	atomic.StoreUint32(&c.Status, STATDISCONNECT)
	logger.Debugf("+ [MLEND]++(%s) all coroutines are finished, exiting conn++", c.client.GetIdentifier())
}

// ShutDown terminates the connection.
func (c *Connection) Shutdown() {
	c.Conn.Close()
	logger.Debug("* [Event] Shutting down stream.")
}

// Receive is a routine to read incomming packets from internal connection. It terminates
// a connection immediately if any abnormalities are detected or packet is malformed. This
// function blocks on purpose.
func (c *Connection) receive() (result *[]byte, code byte, length int, err error) {
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
	return &pack, cmd, 0, nil
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
	return packet.IsValidCommand(msg)
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

// HandleSendError is a error handler. It is used for errors
// caused by sending packets. Currently it terminates the
// connection.
func (c *Connection) handleSendError(err error) {
	c.Conn.Close()
}

// terminate shuts down the connection and undoes all side effects.
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
	case protocol.PCONNECT:
		c.State.OnCONNECT(packet)
	case protocol.PCONNACK:
		c.State.OnCONNACK(packet)
	case protocol.PQUEUE:
		c.State.OnQUEUE(packet)
	case protocol.PSUBSCRIBE:
		c.State.OnSUBSCRIBE(packet)
	case protocol.PSUBACK:
		c.State.OnSUBACK(packet)
	case protocol.PPUBLISH:
		c.State.OnPUBLISH(packet)
	case protocol.PPUBACK:
		c.State.OnPUBACK(packet)
	case protocol.PPING:
		c.State.OnPING(packet)
	case protocol.PDISCONNECT:
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
