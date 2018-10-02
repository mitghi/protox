package protocol

/* d e b u g */

// func (self *Connection) SendMessage(message string, topic string) {
// 	msg := NewPublish()
// 	puid := (*msg.Id)
// 	cl := self.GetClient()
// 	clid := cl.GetIdentifier()
// 	if !self.storage.AddOutbound(clid, msg) {
// 		logger.Warn("unable to add outbound message.")
// 	}
// 	idstore := self.storage.GetIDStoreO(clid)
// 	msg.Message = message
// 	msg.Topic = topic
// 	msg.Meta.MessageId = idstore.GetNewID(puid)
// 	logger.FDebug("SendMessage", "msg.Meta.Messageid.", msg.Meta.MessageId)
// 	msg.Encode()
// 	data := msg.Encoded.Bytes()
// 	packet := NewPacket(&data, msg.Command, msg.Encoded.Len())
// 	self.Send(packet)
// }

// Connection is most important struct which implements `ProtocolConnection` interface.
// It has all neccessary informations for a connection between a client and the broker.
// Most of low level logics are handled in `Connection` methods and most of delegate
// methods are called here as well.
// type Connection struct {
// 	// protobase.ProtoConnection
// 	// WrMutex sync.Mutex
// 	Conn            net.Conn // connection section
// 	Reader          *bufio.Reader
// 	Writer          *bufio.Writer  // end
// 	corous          sync.WaitGroup // main section
// 	ShouldTerminate chan struct{}
// 	SendChan        chan *Packet
// 	PrioSendChan    chan *Packet // Priority send channel
// 	ErrChan         chan struct{}
// 	SendLock        sync.Mutex
// 	RecvChan        chan *Packet
// 	ErrorHandler    func(client *protobase.ClientInterface)
// 	State           ConnectionState           // end
// 	server          protobase.ServerInterface // delegate section
// 	auth            protobase.AuthInterface
// 	storage         protobase.MessageStorage
// 	client          protobase.ClientInterface
// 	clientDelegate  func(string, string, string) protobase.ClientInterface // end
// 	connTimeout     int                                                    // state section
// 	heartbeat       int
// 	unclean         uint32
// 	Status          uint32
// 	deadline        *timer.Timer
// 	justStarted     bool // end
// }

// func NewConnection(conn net.Conn) *Connection {
// 	var result *Connection = &Connection{
// 		//initials
// 		Conn:            conn,
// 		Reader:          bufio.NewReader(conn),
// 		Writer:          bufio.NewWriter(conn),
// 		ErrorHandler:    nil,
// 		corous:          sync.WaitGroup{},
// 		justStarted:     true,
// 		State:           nil,
// 		ShouldTerminate: nil,
// 		ErrChan:         nil,
// 		connTimeout:     1,
// 		heartbeat:       1,
// 		SendChan:        nil,
// 		// SendLock:        nil,
// 		client:         nil,
// 		server:         nil,
// 		auth:           nil,
// 		clientDelegate: nil,
// 		Status:         STATDISCONNECT,
// 		//channels
// 		// WrMutex: nil,
// 	}
// 	result.State = NewGenesis(result)

// 	return result
// }

// // AllocateChannels initializes internal send/receive channels. Their creation
// // are deferred to reduce unneccessary memory allocations up until all initial
// // stages including critical checks are passed.
// func (self *Connection) AllocateChannels() {
// 	self.SendChan = make(chan *Packet, 1024)
// 	self.PrioSendChan = make(chan *Packet, 1024)
// 	// self.SendLock = &sync.Mutex{}
// 	self.RecvChan = make(chan *Packet, 1024)
// 	self.ErrChan = make(chan struct{}, 1)
// }

// // Receive is a helper function which creates a new packet from incomming data.
// // Result should be checked later to create/cast to a particular packet.
// func (self *Connection) Receive() (packet *Packet, err error) {
// 	pck, cmd, length, err := self.receive()
// 	if err != nil {
// 		return nil, err
// 	}
// 	resultPacket := NewPacket(pck, cmd, length)

// 	return resultPacket, nil
// }

// // ReceiveWithTimeout waits `timeout` seconds before returning an erorr. It polls a coroutine
// // for `timeout` seconds.
// func (self *Connection) ReceiveWithTimeout(timeout time.Duration, inbox chan *Packet) (pcket *Packet, err error) {
// 	period := time.After(timeout)
// 	self.corous.Add(1)

// 	go func(pchan chan<- *Packet, wg *sync.WaitGroup) {
// 		packet, err := self.Receive()
// 		// NOTE: should not  close the channel here
// 		if err != nil {
// 		} else {
// 			pchan <- packet
// 		}
// 		// Signal workgroups that we are done
// 		wg.Done()
// 	}(inbox, &self.corous)
// 	select {
// 	case packet, ok := <-inbox:
// 		// pchan = nil
// 		// period = nil
// 		if !ok {
// 			return nil, BadMsgTypeError
// 		}
// 		return packet, nil
// 	case <-period:
// 		// pchan = nil
// 		logger.Debug("[PeriodError] Timeout occured")
// 		return nil, CriticalTimeout
// 	}
// }

// // Send routine is responsible to write a packet into send channel. Actuall sending
// // is done by a coroutine.
// func (self *Connection) Send(packet *Packet) {
// 	self.SendLock.Lock()
// 	if self.SendChan != nil {
// 		self.SendChan <- packet
// 	} else {
// 		logger.Debug("[NOTICE]: SendChannel is nil")
// 	}
// 	self.SendLock.Unlock()
// }

// // SendPrio is the s end hadnler for packets with higher priority.
// func (self *Connection) SendPrio(packet *Packet) {
// 	self.SendLock.Lock()
// 	if self.PrioSendChan != nil {
// 		self.PrioSendChan <- packet
// 	} else {
// 		logger.Debug("[NOTICE]: PrioSendChannel is nil")
// 	}
// 	self.SendLock.Unlock()
// }

// // SendDirect directly writes a packet to a socket.
// func (self *Connection) SendDirect(packet *Packet) {
// 	self.send(packet)
// }

// // GetStatus atomically returns the current status.
// func (self *Connection) GetStatus() (stat uint32) {
// 	stat = atomic.LoadUint32(&self.Status)
// 	return stat
// }

// // GetErrChan returns a chan for errors.
// func (self *Connection) GetErrChan() chan struct{} {
// 	return self.ErrChan
// }

// // SetStatus atomically sets the current status. It will be evaluated in the
// // main loop.
// func (self *Connection) SetStatus(status uint32) {
// 	atomic.StoreUint32(&self.Status, status)
// }

// // Send function writes a packet to a socket.
// func (self *Connection) send(packet *Packet) (err error) {
// 	_, err = self.Writer.Write(*packet.Data)
// 	if err != nil {
// 		return err
// 	}
// 	// NOTE: NEW: UNTESTED:
// 	err = self.Writer.Flush()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// import (
// 	"sync/atomic"

// 	"github.com/mitghi/protox/protobase"
// )

// // Online is the second stage. A connection can only be upgraded to `Online` iff it passes
// // `Genesis` stage which means it must be fully authorized, valid and compatible with the
// // broker.
// type Online struct {
// 	ConnectionState

// 	Conn   *Connection
// 	client protobase.ClientInterface
// 	server protobase.ServerInterface
// }

// // NewOnline returns a pointer to a new `Online` struct. This is the where
// // interactions with a connected/authorized client happens.
// func NewOnline(conn *Connection) *Online {
// 	result := &Online{
// 		Conn:   conn,
// 		client: nil,
// 		server: nil,
// 	}

// 	return result
// }

// // SetServer sets the caller ( server).
// func (self *Online) SetServer(server protobase.ServerInterface) {
// 	self.server = server
// }

// // SetClient sets the internal client struct pointer.
// func (self *Online) SetClient(client protobase.ClientInterface) {
// 	self.client = client
// }

// // SetNextState pushes into next state.
// func (self *Online) SetNextState() {
// }

// // Handle is the main handler ( stub for Online ).
// func (self *Online) Handle(packet *Packet) {
// }

// // HandleDefault is the default handler ( stub for Online ).
// func (self *Online) HandleDefault(packet *Packet) (status bool) {
// 	return true
// }

// // Shutdown sets the status to error which notifies the supervisor
// // and cleanly terminates the connection.
// func (self *Online) Shutdown() {
// 	logger.Debug("* [Genesis] Closing.")
// 	atomic.StoreUint32(&(self.Conn).Status, STATERR)
// 	self.client.Disconnected()
// }

// // onCONNECT is not valid in this stage.
// func (self *Online) onCONNECT(packet *Packet) {
// 	// TODO
// }

// // onCONNACK is not valid in this stage.
// func (self *Online) onCONNACK(packet *Packet) {
// 	// TODO
// }

// // onPUBLISH is the handler for `Publish` packets.
// func (self *Online) onPUBLISH(packet *Packet) {
// 	var publish *Publish = NewPublish()
// 	if err := publish.DecodeFrom(packet.Data); err != nil {
// 		logger.Debug("- [DecodeErr(onPublish)] Unable to decode data.", err)
// 		self.Shutdown()
// 		return
// 	}
// 	if stat := self.Conn.storage.AddInbound(self.client.GetIdentifier(), publish); stat == false {
// 		logger.Debug("?? [NOTICE] addinbound returned false (online/publish).")
// 	}
// 	var puback *Puback = NewPuback()
// 	if publish.Meta.Qos > 0 {
// 		puback.Meta.Qos, puback.Meta.MessageId = publish.Meta.Qos, publish.Meta.MessageId
// 		if err := puback.Encode(); err != nil {
// 			logger.FError("onPUBLISH", "- [ONLINE] Error while encoding puback.")
// 			self.Shutdown()
// 			return
// 		}
// 		var pckt *Packet = puback.GetPacket().(*Packet)
// 		self.Conn.SendPrio(pckt)
// 		if stat := self.Conn.storage.DeleteIn(self.client.GetIdentifier(), publish); stat == false {
// 			logger.Debug("?? [NOTICE] deleteinbound returned false (online/publish).")
// 		}
// 	}
// 	pb := NewMsgBox(publish.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
// 	// publish box clone
// 	pbc := pb.Clone(protobase.MDInbound)
// 	// self.client.Publish(publish.Topic, publish.Message, protobase.MDInbound)
// 	self.client.Publish(pbc)
// 	// self.server.NotifyPublish(publish.Topic, publish.Message, self.Conn, protobase.MDInbound)

// 	self.server.NotifyPublish(self.Conn, pb)
// 	// self.server.NotifyPublish(publish.Topic, publish.Message, self.Conn, protobase.MDInbound)
// }

// // onSUBSCRIBE is the handler for `Subscribe` packets.
// func (self *Online) onSUBSCRIBE(packet *Packet) {
// 	subscribe := NewSubscribe()
// 	if err := subscribe.DecodeFrom(packet.Data); err != nil {
// 		self.Shutdown()
// 		return
// 	}
// 	pb := NewMsgBox(subscribe.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
// 	self.client.Subscribe(pb)
// 	self.server.NotifySubscribe(self.Conn, pb)
// 	// self.client.Subscribe(subscribe.Topic)
// 	// self.server.NotifySubscribe(subscribe.Topic, self.Conn)

// }

// // onPING is the heartbeat handler ( other packets reset its timer as well ).
// func (self *Online) onPING(packet *Packet) {
// 	logger.Debug("+ [Heartbeat] Received.")
// }

// // onSUBACK is a handler which removes the outbound subscribe
// // message when QoS >0.
// func (self *Online) onSUBACK(packet *Packet) {
// 	// TODO
// 	logger.FDebug("onSUBACK", "suback received.")
// }

// // onPUBACK is a handler which removes the outbound publish
// // message when QoS >0.
// func (self *Online) onPUBACK(packet *Packet) {
// 	// TODO
// 	logger.FDebug("onPUBACK", "puback received.")
// }

// // SetServer sets the caller ( server).
// func (self *Online) SetServer(server protobase.ServerInterface) {
// 	self.server = server
// }

// // SetClient sets the internal client struct pointer.
// func (self *Online) SetClient(client protobase.ClientInterface) {
// 	self.client = client
// }

// // SetNextState pushes into next state.
// func (self *Online) SetNextState() {
// }

// // Handle is the main handler ( stub for Online ).
// func (self *Online) Handle(packet *Packet) {
// }

// // onCONNECT is not valid in this stage.
// func (self *Online) onCONNECT(packet *Packet) {
// 	// TODO
// }

// // onCONNACK is not valid in this stage.
// func (self *Online) onCONNACK(packet *Packet) {
// 	// TODO
// }

// import (
// 	"sync/atomic"

// 	"github.com/mitghi/protox/protobase"
// )

// // Online is the second stage. A connection can only be upgraded to `Online` iff it passes
// // `Genesis` stage which means it must be fully authorized, valid and compatible with the
// // broker.
// type Online struct {
// 	constate

// 	Conn *Connection
// }

// // NewOnline returns a pointer to a new `Online` struct. This is the where
// // interactions with a connected/authorized client happens.
// func NewOnline(conn *Connection) *Online {
// 	result := &Online{
// 		constate: constate{
// 			constatebase: constatebase{
// 				Conn: conn,
// 			},
// 			client: nil,
// 			server: nil,
// 		},
// 		Conn: conn,
// 	}
// 	return result
// }

// // // SetServer sets the caller ( server).
// // func (self *Online) SetServer(server protobase.ServerInterface) {
// // 	self.server = server
// // }

// // // SetClient sets the internal client struct pointer.
// // func (self *Online) SetClient(client protobase.ClientInterface) {
// // 	self.client = client
// // }

// // HandleDefault is the default handler ( stub for Online ).
// func (self *Online) HandleDefault(packet *Packet) (status bool) {
// 	return true
// }

// // Shutdown sets the status to error which notifies the supervisor
// // and cleanly terminates the connection.
// func (self *Online) Shutdown() {
// 	logger.Debug("* [Genesis] Closing.")
// 	atomic.StoreUint32(&(self.Conn).Status, STATERR)
// 	self.client.Disconnected()
// }

// // onCONNECT is not valid in this stage.
// func (self *Online) onCONNECT(packet *Packet) {
// 	// TODO
// }

// // onCONNACK is not valid in this stage.
// func (self *Online) onCONNACK(packet *Packet) {
// 	// TODO
// }

// // onPUBLISH is the handler for `Publish` packets.
// func (self *Online) onPUBLISH(packet *Packet) {
// 	var publish *Publish = NewPublish()
// 	if err := publish.DecodeFrom(packet.Data); err != nil {
// 		logger.Debug("- [DecodeErr(onPublish)] Unable to decode data.", err)
// 		self.Shutdown()
// 		return
// 	}
// 	if stat := self.Conn.storage.AddInbound(self.client.GetIdentifier(), publish); stat == false {
// 		logger.Debug("?? [NOTICE] addinbound returned false (online/publish).")
// 	}
// 	var puback *Puback = NewPuback()
// 	if publish.Meta.Qos > 0 {
// 		puback.Meta.Qos, puback.Meta.MessageId = publish.Meta.Qos, publish.Meta.MessageId
// 		if err := puback.Encode(); err != nil {
// 			logger.FError("onPUBLISH", "- [ONLINE] Error while encoding puback.")
// 			self.Shutdown()
// 			return
// 		}
// 		var pckt *Packet = puback.GetPacket().(*Packet)
// 		self.Conn.SendPrio(pckt)
// 		if stat := self.Conn.storage.DeleteIn(self.client.GetIdentifier(), publish); stat == false {
// 			logger.Debug("?? [NOTICE] deleteinbound returned false (online/publish).")
// 		}
// 	}
// 	pb := NewMsgBox(publish.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
// 	// publish box clone
// 	pbc := pb.Clone(protobase.MDInbound)
// 	// self.client.Publish(publish.Topic, publish.Message, protobase.MDInbound)
// 	self.client.Publish(pbc)
// 	// self.server.NotifyPublish(publish.Topic, publish.Message, self.Conn, protobase.MDInbound)

// 	self.server.NotifyPublish(self.Conn, pb)
// 	// self.server.NotifyPublish(publish.Topic, publish.Message, self.Conn, protobase.MDInbound)
// }

// // onSUBSCRIBE is the handler for `Subscribe` packets.
// func (self *Online) onSUBSCRIBE(packet *Packet) {
// 	subscribe := NewSubscribe()
// 	if err := subscribe.DecodeFrom(packet.Data); err != nil {
// 		self.Shutdown()
// 		return
// 	}
// 	pb := NewMsgBox(subscribe.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
// 	self.client.Subscribe(pb)
// 	self.server.NotifySubscribe(self.Conn, pb)
// 	// self.client.Subscribe(subscribe.Topic)
// 	// self.server.NotifySubscribe(subscribe.Topic, self.Conn)

// }

// // onPING is the heartbeat handler ( other packets reset its timer as well ).
// func (self *Online) onPING(packet *Packet) {
// 	logger.Debug("+ [Heartbeat] Received.")
// }

// // onSUBACK is a handler which removes the outbound subscribe
// // message when QoS >0.
// func (self *Online) onSUBACK(packet *Packet) {
// 	// TODO
// 	logger.FDebug("onSUBACK", "suback received.")
// }

// // onPUBACK is a handler which removes the outbound publish
// // message when QoS >0.
// func (self *Online) onPUBACK(packet *Packet) {
// 	// TODO
// 	logger.FDebug("onPUBACK", "puback received.")
// }

// import (
// 	"fmt"

// 	"github.com/mitghi/protox/protobase"
// )

// // Genesis is the initial and most important stage. All new connections can connect
// // to broker iff they pass this stage. This stage only accepts `Connect` packets.
// // Any other control packet results in immediate termination ( it can be adjusted using
// // policies.
// type Genesis struct {
// 	constate

// 	Conn           *Connection
// 	gotFirstPacket bool
// }

// // NewGenesis creates a pointer to a new `Gensis` packet.
// func NewGenesis(conn *Connection) *Genesis {
// 	result := &Genesis{
// 		constate: constate{
// 			constatebase: constatebase{
// 				Conn: conn,
// 			},
// 			client: nil,
// 			server: nil,
// 		},
// 		Conn:           nil,
// 		gotFirstPacket: false,
// 	}
// 	result.Conn = result.constate.constatebase.Conn.(*Connection)
// 	if result.Conn == nil {
// 		// TODO
// 		// . this is critical.
// 	}
// 	return result
// }

// // SetNextState pushes the state machine into its next stage.
// // Initially it is from Genesis to Online ( Genesis -> Online -> .... ).
// func (self *Genesis) SetNextState() {
// 	newState := NewOnline(self.Conn)
// 	newState.SetClient(self.client)
// 	newState.SetServer(self.server)
// 	self.Conn.State = newState

// 	logger.Debug("+ [Genesis] Genesis for client [Status] ready.")
// }

// // cleanUp is a routine which removes pointers from the struct.
// func (self *Genesis) cleanUp() {
// 	self.Conn = nil
// 	self.client = nil
// 	self.server = nil
// }

// // Shutdown terminates the state and calls the handlers to terminate
// // and undo all side effects.
// func (self *Genesis) Shutdown() {
// 	self.client.Disconnected()
// }

// // Handle is only a stub to satisfy interface requirements ( for Genesis stage ).
// func (self *Genesis) Handle(packet *Packet) {
// }

// // HandleDefault is the first function invoked in `Genesis` when a new state struct is created.
// // It passes credentials from `Connect` packet to a `AuthInterface` implementor and upgrades
// // from `Genesis` to `Online` stage. It sends a `Connack` with appropirate status code, regardless.
// func (self *Genesis) HandleDefault(packet *Packet) (status bool) {
// 	// TODO
// 	//  add defer to cleanUp and check its performance impact
// 	var (
// 		// by default, assume packet is invalid
// 		valid   bool                    = false
// 		p       *Connect                = NewConnect()
// 		cack    *Connack                = NewConnack()
// 		authsys protobase.AuthInterface = self.Conn.GetAuthenticator()
// 		creds   protobase.CredentialsInterface
// 		rpacket *Packet
// 		newcl   protobase.ClientInterface
// 	)

// 	logger.FDebug("HandleDefault", "raw packet content", fmt.Sprintf("% #x\n", (*(*packet).Data)))
// 	// terminate immediately if packet is malformed or invalid.
// 	if err := p.DecodeFrom((*packet).Data); err != nil {
// 		logger.Debug("- [Fatal] invalid connection packet.", err)
// 		// TODO
// 		//  undo side effects
// 		self.gotFirstPacket = false
// 		self.cleanUp()
// 		return false
// 	}
// 	logger.FDebug("HandleDefault", "conn packet content.", p.String())
// 	// connection is established, can push into the next state
// 	self.gotFirstPacket = true
// 	// TODO
// 	// . improve by directly pass connect packet to auth subsystem
// 	creds, err := authsys.MakeCreds(p.Username, p.Password, p.ClientId)
// 	if err != nil {
// 		logger.Fatal("- [Fatal] cannot make credentials.", err)
// 		return false
// 	}
// 	// TODO/NOTICE
// 	//  do not create a new client until credentials are valid ( reduce memory alloc. overhead )
// 	newcl = self.Conn.clientDelegate(p.Username, p.Password, p.ClientId)
// 	// NOTE: check error explicitely
// 	valid, _ = authsys.CanAuthenticate(creds)
// 	if valid, err = authsys.CanAuthenticate(creds); valid {
// 		// newcl.SetCreds(creds)
// 		// NOTE
// 		// . don't forget to add result code explicitely
// 		// cack.SetResultCode(OK)
// 		cack.Encode()
// 		rpacket = cack.GetPacket().(*Packet)
// 		self.Conn.SetClient(newcl)
// 		self.SetNextState() // Genesis -> Online
// 		self.Conn.SendDirect(rpacket)
// 		// TODO
// 		//  these lines are moves to cleanUp, remove them when
// 		//  its finalized.
// 		//  self.Conn = nil
// 		//  self.client = nil
// 		self.cleanUp()
// 		return true
// 	} else {
// 		cack.SetResultCode(TMP_RESPOK) // <- NOTE: TODO: this means NOT OK, needs refactoring
// 		cack.Encode()
// 		rpacket = cack.GetPacket().(*Packet)
// 		self.Conn.SetClient(newcl)
// 		self.Conn.SendDirect(rpacket)
// 		// TODO
// 		//  improve error handling
// 		self.client.Disconnected()
// 		self.cleanUp()
// 		return false
// 	}
// }

// --------------------------------------------------------

// SendHandler reads from send channel and writes to a socket.
// func (self *protocon) sendHandler() {
// 	const fname = "sendHandler"
// 	defer func() {
// 		logger.FInfo(fname, "before decrementing workgroup")
// 		self.corous.Done()
// 	}()

// 	for {
// 		select {
// 		case packet, ok := <-self.SendChan:
// 			if !ok {
// 				logger.FDebug(fname, "PrioSendChan is closed", "ok status:", ok)
// 				return
// 			}
// 			if err := self.send(packet); err != nil {
// 				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets", "error:", err)
// 				return
// 			}
// 		case <-self.cendch:
// 			logger.FDebug(fname, "* [PrioSendHandler] received end signal from [cendch] channel, terminating coroutine.")
// 			return
// 		}
// 	}
// 	// for packet := range self.SendChan {
// 	// 	// Check QoS
// 	// 	err := self.send(packet)
// 	// 	if err != nil {
// 	// 		// self.handleSendError(err)
// 	// 		// return
// 	// 		// TODO
// 	// 		// Handle this case
// 	// 		logger.Debug("- [ERROR IN SENDHANDLER, BEFORE BREAKING]")
// 	// 		return
// 	// 	}
// 	// }
// 	// self.corous.Done()
// }

// -----------------------------------------------------------

//-----------------------------------------------------------------------------

// type COnline struct {
// 	constate

// 	Conn *CLBConnection
// }

// // MARK: Online

// // NewOnline returns a pointer to a new `Online` struct. This is the where
// // interactions with a connected/authorized client happens.
// func NewCOnline(conn *CLBConnection) *COnline {
// 	result := &COnline{
// 		constate: constate{
// 			constatebase: constatebase{
// 				Conn: conn,
// 			},
// 			client: nil,
// 			server: nil,
// 		},
// 		Conn: conn,
// 	}
// 	result.client = conn.GetClient()
// 	return result
// }

// // HandleDefault is the default handler ( stub for COnline ).
// func (self *COnline) HandleDefault(packet *Packet) (status bool) {
// 	return true
// }

// // Shutdown sets the status to error which notifies the supervisor
// // and cleanly terminates the connection.
// func (self *COnline) Shutdown() {
// 	logger.Debug("* [Genesis] Closing.")
// 	atomic.StoreUint32(&(self.Conn).Status, STATERR)
// }

// // onCONNECT is not valid in this stage.
// func (self *COnline) onCONNECT(packet *Packet) {
// 	// TODO
// 	// NOTE
// 	// . this is new
// 	self.Shutdown()
// 	self.Conn.protocon.Conn.Close()
// }

// // onCONNACK is not valid in this stage.
// func (self *COnline) onCONNACK(packet *Packet) {
// 	// TODO
// }

// // onPUBLISH is the handler for `Publish` packets.
// func (self *COnline) onPUBLISH(packet *Packet) {
// 	var publish *Publish = NewPublish()
// 	if err := publish.DecodeFrom(packet.Data); err != nil {
// 		logger.Debug("- [DecodeErr(onPublish)] Unable to decode data.", err)
// 		self.Shutdown()
// 		return
// 	}
// 	if stat := self.Conn.storage.AddInbound(publish); stat == false {
// 		logger.Debug("? [NOTICE] addinbound returned false (conline/publish).")
// 	}
// 	var puback *Puback = NewPuback()
// 	if publish.Meta.Qos > 0 {
// 		puback.Meta.Qos, puback.Meta.MessageId = publish.Meta.Qos, publish.Meta.MessageId
// 		if err := puback.Encode(); err != nil {
// 			logger.FError("onPUBLISH", "- [CONLINE] Error while encoding puback.")
// 			self.Shutdown()
// 			return
// 		}
// 		logger.FTracef(1, "onPUBLISH", "* [QoS] packet QoS(%b) Duplicate(%t) MessageID(%d).", publish.Meta.Qos, publish.Meta.Dup, int(publish.Meta.MessageId))
// 		var pckt *Packet = puback.GetPacket().(*Packet)
// 		logger.FTrace(1, "onPUBLISH", "* [PubAck] sending packet with content", pckt.Data)
// 		// NOTE
// 		// . this has changed
// 		self.Conn.Send(pckt)
// 		// self.Conn.SendPrio(pckt)
// 		if stat := self.Conn.storage.DeleteIn(publish); stat == false {
// 			logger.Debug("? [NOTICE] deleteinbound returned false (conline/publish).")
// 		}
// 	}
// 	pb := NewMsgBox(publish.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
// 	// publish box clone
// 	pbc := pb.Clone(protobase.MDInbound)
// 	self.client.Publish(pbc)

// }

// // onSUBSCRIBE is the handler for `Subscribe` packets.
// func (self *COnline) onSUBSCRIBE(packet *Packet) {
// 	logger.Debug("* [Subscribe] packet is received.")
// 	// subscribe := NewSubscribe()
// 	// if err := subscribe.DecodeFrom(packet.Data); err != nil {
// 	// 	self.Shutdown()
// 	// 	return
// 	// }
// 	// pb := NewMsgBox(subscribe.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
// 	// self.client.Subscribe(pb)
// 	// self.server.NotifySubscribe(self.Conn, pb)
// 	// self.client.Subscribe(subscribe.Topic)
// 	// self.server.NotifySubscribe(subscribe.Topic, self.Conn)
// }

// // onPING is the heartbeat handler ( other packets reset its timer as well ).
// func (self *COnline) onPING(packet *Packet) {
// 	logger.Debug("+ [Heartbeat] Received.")
// }

// // onSUBACK is a handler which removes the outbound subscribe
// // message when QoS >0.
// func (self *COnline) onSUBACK(packet *Packet) {
// 	// TODO
// 	var (
// 		pa  *Suback = NewSuback()
// 		uid uuid.UUID
// 	)
// 	logger.FDebug("onSUBACK", "* [SubAck] packet is received.")
// 	if err := pa.DecodeFrom(packet.Data); err != nil {
// 		logger.FDebug("onSUBACK", "- [Decode] uanble to decode in [SubAck].", err)
// 		return
// 	}

// 	oidstore := self.Conn.storage.GetIDStoreO()
// 	msgid := pa.Meta.MessageId
// 	uid, ok := oidstore.GetUUID(msgid)
// 	if !ok {
// 		logger.FWarn("onSUBACK", "- [IDStore] no packet with msgid found.", "msgid", msgid)
// 		return
// 	}
// 	np, ok := self.Conn.storage.GetOutbound(uid)
// 	if !ok {
// 		logger.FWarn("onSUBACK", "- [MessageBox] no packet with uid found.", uid)
// 	}
// 	if !self.Conn.storage.DeleteOut(np) {
// 		logger.FWarn("onSUBACK", "- [MessageBox] failed to remove message.")
// 	}
// 	oidstore.FreeId(msgid)
// }

// // onPUBACK is a handler which removes the outbound publish
// // message when QoS >0.
// func (self *COnline) onPUBACK(packet *Packet) {
// 	// TODO
// 	var (
// 		pa  *Puback = NewPuback()
// 		uid uuid.UUID
// 	)

// 	logger.FDebug("onPUBACK", "+ [PubAck] packet received.")
// 	if err := pa.DecodeFrom(packet.Data); err != nil {
// 		logger.FDebug("onPUBACK", "- [Decode] uanble to decode in [PubAck].", err)
// 		return
// 	}

// 	oidstore := self.Conn.storage.GetIDStoreO()
// 	msgid := pa.Meta.MessageId
// 	uid, ok := oidstore.GetUUID(msgid)
// 	if !ok {
// 		logger.FWarn("onPUBACK", "- [IDStore] no packet with msgid found.", "msgid", msgid)
// 		return
// 	}
// 	np, ok := self.Conn.storage.GetOutbound(uid)
// 	if !ok {
// 		logger.FWarn("onPUBACK", "- [MessageBox] no packet with uid found.", uid)
// 	}
// 	if !self.Conn.storage.DeleteOut(np) {
// 		logger.FWarn("onPUBACK", "- [MessageBox] failed to remove message.")
// 	}
// 	oidstore.FreeId(msgid)
// }
//------------------------------------------------------------------------------

// // Genesis is the initial and most important stage. All new connections can connect
// // to broker iff they pass this stage. This stage only accepts `Connect` packets.
// // Any other control packet results in immediate termination ( it can be adjusted using
// // policies.
// type CGenesis struct {
// 	constate

// 	Conn           *CLBConnection
// 	gotFirstPacket bool
// }

// // NewGenesis creates a pointer to a new `Gensis` packet.
// func NewCGenesis(conn *CLBConnection) *CGenesis {
// 	// result := &CGenesis{
// 	// 	Conn:           conn,
// 	// 	gotFirstPacket: false,
// 	// }
// 	result := &CGenesis{
// 		constate: constate{
// 			constatebase: constatebase{
// 				Conn: conn,
// 			},
// 			client: nil,
// 			server: nil,
// 		},
// 		Conn:           nil,
// 		gotFirstPacket: false,
// 	}
// 	result.Conn = result.constate.constatebase.Conn.(*CLBConnection)
// 	result.client = conn.GetClient()
// 	if result.Conn == nil {
// 		// TODO
// 		// . this is critical.
// 	}
// 	return result
// }

// // SetNextState pushes the state machine into its next stage.
// // Initially it is from CGenesis to Online ( CGenesis -> Online -> .... ).
// func (self *CGenesis) SetNextState() {
// 	newState := NewCOnline(self.Conn)
// 	self.Conn.State = newState

// 	logger.Debug("+ [CGenesis] CGenesis for client [Status] ready.")
// }

// // cleanUp is a routine which removes pointers from the struct.
// func (self *CGenesis) cleanUp() {
// 	self.Conn = nil
// 	self.client = nil
// }

// // Shutdown terminates the state and calls the handlers to terminate
// // and undo all side effects.
// func (self *CGenesis) Shutdown() {
// 	self.client.Disconnected(nil)
// }

// // Handle is only a stub to satisfy interface requirements ( for CGenesis stage ).
// func (self *CGenesis) Handle(packet *Packet) {
// }

// func (self *CGenesis) onCONNACK(packet *Packet) {
// 	logger.FTrace(1, "onCONNACK", "+ [ConnAck] packet received.")
// 	var p *Connack = NewConnack()
// 	if err := p.DecodeFrom((*packet).Data); err != nil {
// 		logger.FTrace(1, "onCONNACK", "- [Fatal] invalid connack packet.", err)
// 	}
// 	logger.FTrace(1, "onCONNACK", "* [ConnAck] connack content.", p)
// 	// NOTE
// 	// . this has changed
// 	// if p.ResultCode == TMP_RESPOK {
// 	if p.ResultCode == RESPFAIL {
// 		logger.FTrace(1, "onCONNACK", "- [Credentials] INVALID CREDENTIALS.")
// 		self.Conn.SetStatus(STATERR)
// 		return
// 	} else if p.ResultCode == RESPOK {
// 		self.Conn.SetStatus(STATONLINE)
// 		logger.FTrace(1, "onCONNACK", "+ [Credentials] are valid and [Client] is now (Online).")
// 		// set connack result to options
// 		// and pass to to next state.
// 		for k, _ := range self.Conn.stateOpts {
// 			delete(self.Conn.stateOpts, k)
// 		}
// 		caopt := NewConnackOpts()
// 		caopt.parseFrom(p)
// 		self.Conn.stateOpts[CCONNACK] = caopt
// 		self.SetNextState()
// 		self.cleanUp()
// 		return
// 	} else {
// 		// NOTE: TODO:
// 		// . THIS IS FATAL, check invalid codes
// 		logger.FTracef(1, "onCONNACK", "- [Connack/Resp] unknown response-code(%b) in packet.", p.ResultCode)
// 		return
// 	}
// }

// // HandleDefault is the first function invoked in `CGenesis` when a new state struct is created.
// // It passes credentials from `Connect` packet to a `AuthInterface` implementor and upgrades
// // from `CGenesis` to `Online` stage. It sends a `Connack` with appropirate status code, regardless.
// func (self *CGenesis) HandleDefault(packet *Packet) (ok bool) {
// 	const fn = "HandleDefault"
// 	// TODO
// 	//  add defer to cleanUp and check its performance impact
// 	var (
// 		newcl   protobase.ClientInterface = self.Conn.GetClient()
// 		p       *Connect                  = NewConnect()
// 		Conn    *CLBConnection            = self.Conn
// 		nc      net.Conn
// 		err     error
// 		rpacket *Packet
// 		addr    string
// 		// by default, assume packet is invalid
// 		// valid   bool                      = false
// 	)
// 	defer func() {
// 		if !ok {
// 			Conn.protocon.Lock()
// 			if Conn.protocon.Conn != nil {
// 				Conn.protocon.Conn = nil
// 				Conn.protocon.Writer = nil
// 				Conn.protocon.Reader = nil
// 			}
// 			Conn.protocon.Unlock()
// 		}
// 	}()

// 	Conn.protocon.RLock()
// 	addr = Conn.protocon.addr
// 	Conn.protocon.RUnlock()
// 	if Conn.protocon.Conn == nil {
// 		if c, err := dialRemoteAddr(addr, false); err != nil {
// 			logger.FDebug(fn, "- [TCPConnect] cannot connect to remote addr.", "error", err)
// 			ok = false
// 			return
// 		} else {
// 			nc = c
// 		}
// 	}
// 	if nc == nil {
// 		ok = false
// 		return
// 	}
// 	Conn.protocon.Lock()
// 	Conn.SetNetConnection(nc)
// 	Conn.protocon.Unlock()

// 	p.Username, p.Password, p.ClientId = newcl.GetCreds().GetCredentials()
// 	if err = p.Encode(); err != nil {
// 		logger.FFatal("HandleDefault", "- [Encode] cannot encode in [CGenesis].", err)
// 		ok = false
// 		return
// 	}
// 	rpacket = p.GetPacket().(*Packet)
// 	if err = Conn.protocon.send(rpacket); err != nil {
// 		logger.Debug("- [Send] send returned an error.", err)
// 		ok = false
// 		return
// 	}

// 	ok = true
// 	return

// 	// TODO
// 	//  improve error handling
// 	// self.cleanUp()
// }

// func (self *CGenesis) onDISCONNECT(packet *Packet) {
// 	// TODO
// 	logger.FDebug("onDISCONNECT", "* [Disconnect] packet received.")
// 	self.Conn.protocon.Conn.Close()
// }

// ------------------------------------------------
// func (self *CLBConnection) SendRedelivery() {
// 	const fn = "SendRedelivery"
// 	var (
// 		clid     string = self.GetClient().GetIdentifier()
// 		outbound []protobase.EDProtocol
// 		packet   *Packet
// 		// msg      *Publish
// 		// tp     *Packet
// 	)
// 	logger.FDebug(fn, "+ [Redeliver] starting redelivery.")
// 	outbound = self.storage.GetAllOut()
// 	for _, p := range outbound {
// 		fmt.Printf("client %s has this packet %+v\n", clid, p)
// 		if self.GetStatus() == protobase.STATONLINE {
// 			// NOTE
// 			// . this has changed
// 			// msg = NewPublish()
// 			// tp = p.GetPacket().(*Packet)
// 			// if err := msg.DecodeFrom(tp.Data); err != nil {
// 			// 	logger.FDebug("sendRedelivery", "- [Decode] cannot decode a publish packet in [Redeliver].")
// 			// }
// 			// msg.Meta.Dup = true
// 			// logger.FDebug("SendRedelivery", "+ [Redliver] undelivered packages are in their path to broker.")
// 			// msg.Encode()
// 			// packet = msg.GetPacket().(*Packet)
// 			// self.Send(packet)

// 			switch p.(type) {
// 			case *Publish:
// 				tmp := p.(*Publish)
// 				tmp.Meta.Dup = true
// 				// tmp.Meta.Qos = p.QoS()
// 				tmp.Encoded = nil
// 				tmp.Encode()
// 			}

// 			packet = p.GetPacket().(*Packet)
// 			logger.FDebug("SendRedelivery", "+ [Redliver] undelivered packages are in their path to broker.")
// 			self.Send(packet)

// 		}
// 	}
// }

// func (self *CLBConnection) Publish(topic string, message []byte, qos byte) {
// 	logger.Infof("* [Publish/QoS] qos is (%b) [Topic] is (%s) [Message] is (%s).", qos, topic, message)
// 	clid := self.GetClient().GetIdentifier()
// 	pb := NewPublish()
// 	puid := (*pb.Id)
// 	pb.Topic = topic
// 	pb.Message = message

// 	if qos > 0 {
// 		pb.Meta.Qos = qos
// 		idstore := self.storage.GetIDStoreO()
// 		pb.Meta.MessageId = idstore.GetNewID(puid)
// 		logger.Infof("* [Publish<-] with id (%d).", pb.Meta.MessageId)
// 		if !self.storage.AddOutbound(pb) {
// 			logger.Info("? [NOTICE] unable to add outbound packet to [MessageBox].", "userId", clid)
// 		}
// 	}
// 	pb.Encode()
// 	packet := pb.GetPacket().(*Packet)
// 	if self.GetStatus() == STATONLINE {
// 		self.Send(packet)
// 	} else {
// 		// NOTE:
// 		// . drop packages when qos == 0
// 		logger.FWarn("Publish", "- [Publish] dropping packet due to status.")
// 	}
// }

// func (self *CLBConnection) Subscribe(topic string, qos byte) {
// 	logger.Infof("* [Subscribe/QoS] qos is (%b) [Topic] is (%s).", qos, topic)
// 	clid := self.GetClient().GetIdentifier()
// 	sb := NewSubscribe()
// 	puid := (*sb.Id)
// 	sb.Topic = topic

// 	if qos > 0 {
// 		sb.Meta.Qos = qos
// 		idstore := self.storage.GetIDStoreO()
// 		sb.Meta.MessageId = idstore.GetNewID(puid)
// 		logger.Infof("* [Subscribe->] with id (%d).", sb.Meta.MessageId)
// 		if !self.storage.AddOutbound(sb) {
// 			logger.Info("? [NOTICE] unable to add outbound packet to [MessageBox].", "userId", clid)
// 		}
// 	}
// 	sb.Encode()
// 	packet := sb.GetPacket().(*Packet)
// 	if self.GetStatus() == STATONLINE {
// 		self.Send(packet)
// 	} else {
// 		// NOTE:
// 		// . drop packages when qos == 0
// 		logger.FWarn("Publish", "- [Subscribe] dropping packet due to status.")
// 	}
// }

// // Precheck validate packets to avoid abnormalities and malformed packet.
// // It returns a bool indicating wether a control packet is valid or not.
// func (self *protocon) precheck(cmd *byte) bool {
// 	msg := *cmd

// 	if self.justStarted == true {
// 		// dirty check for connect packet
// 		// reject any other packet in this stage
// 		if msg == 0x01 {
// 			// self.gotFirstPacket = true <- is relocated
// 			return true
// 		}
// 		return false
// 	}
// 	return IsValidCommand(msg)
// }
