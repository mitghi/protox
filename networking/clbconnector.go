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
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
	"github.com/mitghi/timer"
)

// ensure interface (protocol) conformance
var _ protobase.ProtoClientConnection = (*CLBConnection)(nil)

// Client Connector error messages
var (
	ECLBSendFailure = errors.New("protocol(clientConnector): unable to send packet.")
)

// Constants
var (
	cDefaultSleepTime time.Duration = time.Millisecond * 2
)

type CLBConnection struct {
	// TODO
	// . check alignment
	// . add packet processor delegate
	// . add callback manager
	*protocon

	clock          *sync.RWMutex
	clblock        *sync.RWMutex
	State          protobase.ConnectionState
	storage        protobase.MessageBox
	client         protobase.ClientInterface
	pinger         *timer.Timer
	heartbeat      int
	justStarted    bool
	shouldContinue bool
	stateOpts      map[byte]protobase.OptionInterface
	clbpub         map[uint16]func(protobase.OptionInterface, protobase.MsgInterface)
	clbsub         map[uint16]func(protobase.OptionInterface, protobase.MsgInterface)
}

func dialRemoteAddr(addr string, tlsopts *tls.Config, isTLS bool) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)
	if isTLS {
		conn, err = tls.Dial("tcp", addr, tlsopts)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewClientConnection(addr string) (clbc *CLBConnection) {
	// TODO
	// . add options
	clbc = &CLBConnection{
		protocon: &protocon{
			// NOTE:
			// . connection is purposefully not initialized
			Conn:            nil,
			Reader:          nil,
			Writer:          nil,
			corous:          sync.WaitGroup{},
			ErrChan:         nil,
			ShouldTerminate: nil,
			SendChan:        nil,
			Status:          STATDISCONNECT,
			addr:            addr,
		},
		clock:          &sync.RWMutex{},
		clblock:        &sync.RWMutex{},
		heartbeat:      CCLBConnectionDefaultHeartbeat,
		pinger:         nil,
		client:         nil,
		storage:        nil,
		shouldContinue: true,
		stateOpts:      make(map[byte]protobase.OptionInterface),
		clbpub:         make(map[uint16]func(protobase.OptionInterface, protobase.MsgInterface)),
		clbsub:         make(map[uint16]func(protobase.OptionInterface, protobase.MsgInterface)),
	}
	// set the state to client genesis
	clbc.State = NewCGenesis(clbc)
	return clbc
}

func (clbc *CLBConnection) GetTermChan() chan struct{} {
	var ch chan struct{}
	clbc.clock.Lock()
	ch = clbc.ShouldTerminate
	clbc.clock.Unlock()
	return ch
}

func (clbc *CLBConnection) SetHeartBeat(heartbeat int) {
	clbc.heartbeat = heartbeat
}

func (clbc *CLBConnection) GetConnection() net.Conn {
	var conn net.Conn
	clbc.protocon.RLock()
	conn = clbc.Conn
	clbc.protocon.RUnlock()
	return conn
}

// SendMessage creates a new packet and sends it to send channel.
func (clbc *CLBConnection) SendMessage(pb protobase.MsgInterface) (err error) {
	const fn = "SendMessage"
	var (
		envelope protobase.MsgEnvelopeInterface = pb.Envelope()
		topic    string                         = envelope.Route()
		message  []byte                         = envelope.Payload()

		msg  *Publish = protocol.NewRawPublish()
		qos  byte     = pb.QoS()
		puid uuid.UUID

		data   []byte
		packet *Packet
	)
	msg.Message = ([]byte)(message)
	msg.Topic = topic
	if qos > 0 {
		logger.FDebug(fn, "* [QoS] QoS>0 in [SendMessage].", "qos", qos)
		puid = (msg.Id)
		logger.FDebugf(fn, "* [SendMessage] publish QoS(%b), message direction(%b).", pb.QoS(), pb.Dir())
		if !clbc.storage.AddOutbound(msg) {
			logger.Warn("- [SendMessage ]unable to add outbound message.")
		}
		idstore := clbc.storage.GetIDStoreO()
		msg.Meta.MessageId = idstore.GetNewID(puid)
		logger.FDebugf(fn, "* [MessageId] messageid(%d) in [SendMessage].", msg.Meta.MessageId)
	}
	err = msg.Encode()
	if err != nil {
		return err
	}
	// data = msg.Encoded.Bytes()
	data = msg.GetBytes()
	packet = NewPacket(data, msg.Command, msg.Encoded.Len())
	clbc.Send(packet)
	return nil
}

func (clbc *CLBConnection) Disconnect() (err error) {
	const fn string = "Disconnect"
	if clbc.GetStatus() == protobase.STATONLINE {
		err = clbc.sendDisconnect()
		if err != nil {
			logger.FDebug(fn, "- [CLBConnection] unable to disconnect. err:", err)
		}
		clbc.SetStatus(protobase.STATGODOWN)
		return err
	}
	logger.FWarn(fn, "- [CLBConnection] invalid connection status to disconnect.")
	return ECLBCONNINVALDISCONN
}

func (clbc *CLBConnection) sendDisconnect() (err error) {
	const fn string = "sendDisconnect"
	var (
		disconn *Disconnect = protocol.NewRawDisconnect()
		packet  *Packet
	)
	err = disconn.Encode()
	if err != nil {
		logger.FWarnf(fn, "- [CLBConnection] unable to encode disconnect packet. error:", err)
		return err
	}
	packet = disconn.GetPacket().(*Packet)
	clbc.protocon.Send(packet)
	return nil
}

func (clbc *CLBConnection) sendPing() (err error) {
	const fn string = "sendPing"
	var (
		ping   *Ping = protocol.NewRawPing()
		packet *Packet
	)
	err = ping.Encode()
	if err != nil {
		logger.FWarnf(fn, "- [CLBConnection] unable to encode ping packet. error:", err)
		return err
	}
	packet = ping.GetPacket().(*Packet)
	clbc.Send(packet)
	return nil
}

func (clbc *CLBConnection) Send(packet *Packet) {
	clbc.SendLock.Lock()
	if clbc.SendChan != nil {
		if clbc.pinger != nil {
			clbc.pinger.Reset(time.Duration(time.Second * time.Duration(clbc.heartbeat)))
		}
		clbc.SendChan <- packet
	} else {
		logger.Info("* [NOTICE] [SendChannel] is (nil).")
	}
	clbc.SendLock.Unlock()
}

// SendHandler reads from send channel and writes to a socket.
func (clbc *CLBConnection) sendHandler() {
	const fname = "sendHandler"
	var (
		dur time.Duration = time.Second * time.Duration(clbc.heartbeat)
	)
	defer func() {
		logger.FInfo(fname, "+ [WorkGroup] decrementing workgroup.")
		clbc.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-clbc.SendChan:
			if !ok {
				logger.FDebugf(fname, "* [SendHandler] is (closed) with status(%t).", ok)
				return
			}
			clbc.send(packet)
		case <-clbc.cendch:
			logger.FDebug(fname, "* [SendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		case _, ok := <-clbc.pinger.C:
			if !ok {
				logger.FDebug(fname, "- [HeartBeat] cannot get a new ping from pinger.")
				return
			}
			_ = clbc.sendPing()
			logger.Info("* [CLBConnection] sending ping.")
		}
		clbc.pinger.Reset(dur)
	}
}

func (clbc *CLBConnection) uniSendHandler() {
	const (
		fname = "uniSendHandler"
	)
	var (
		dur time.Duration = time.Second * time.Duration(clbc.heartbeat)
	)
	defer func() {
		logger.FDebug(fname, "+ [WorkGroup] decrementing workgroup")
		// stop heartbeat ticker
		clbc.pinger.Stop()
		clbc.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-clbc.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] PrioSendChan is closed.", "ok status:", ok)
				return
			}
			if err := clbc.send(packet); err != nil {
				logger.FError(fname, "- [UniSendHandler] error while sending packets to [PrioSendChan].", "error:", err)
				return
			}
		case packet, ok := <-clbc.SendChan:
			if !ok {
				logger.FDebug(fname, "- [UniSendHandler] SendChan is closed.", "ok status:", ok)
				return
			}
			if err := clbc.send(packet); err != nil {
				logger.FError(fname, "- [PrioSendHandler] error while sending priority packets.", "error:", err)
				return
			}
		case _, ok := <-clbc.pinger.C:
			if !ok {
				logger.FDebug("Handle", "- [HeartBeat] cannot get a new ping from pinger.")
				return
			}
			clbc.sendPing()
			logger.Info("* [CLBConnection] sending ping.")
		case <-clbc.cendch:
			logger.FDebug(fname, "* [UniSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
		clbc.pinger.Reset(dur)
	}
}

//
func (clbc *CLBConnection) prioSendHandler() {
	const fname = "prioSendHandler"
	defer func() {
		logger.FInfo(fname, "+ [WorkGroup] decrementing workgroup.")
		clbc.corous.Done()
	}()

	for {
		select {
		case packet, ok := <-clbc.PrioSendChan:
			if !ok {
				logger.FDebug(fname, "[PrioSendChan] is (closed) with status(%t).", ok)
				return
			}
			clbc.SendPrio(packet)
		case <-clbc.cendch:
			logger.FDebug(fname, "* [PrioSendHandler] received end signal from [cendch] channel, terminating coroutine.")
			return
		}
	}
}

func (clbc *CLBConnection) ContinueFlag(f bool) {
	/* critical section */
	clbc.clock.Lock()
	clbc.shouldContinue = f
	clbc.clock.Unlock()
	/* critical section - end */
}

// Handle is the entry routine into `Connection`. It is the main loop
// for handling initial logics/allocating and passing data to different stages.
func (clbc *CLBConnection) Handle(_ protobase.PacketInterface) {
	// initalize variables regardless of
	// possibility of early return
	// variable definition instructions
	// gets rearranged by the compiler** ( TODO : dig into golang source code )
	var (
		dur  time.Duration = time.Second * time.Duration(clbc.heartbeat) // heartbeat interval
		ch   chan *Packet                                                // timeout channel ( initial packet )
		stat uint32                                                      // connection status
	)
	/* critical section */
	clbc.clock.RLock()
	shcont := clbc.shouldContinue
	clbc.clock.RUnlock()
	/* critical section - end */
	if !shcont {
		logger.Debug("* [Flag] discontinue flag is set, terminating....")
		return
	}
	clbc.SetStatus(STATCONNECTING)
	// mark that connect packet is received
	/* critical section */
	clbc.clock.Lock()
	clbc.State = NewCGenesis(clbc)
	clbc.clock.Unlock()
	/* critical section */
	if !clbc.State.HandleDefault(nil) {
		clbc.SetStatus(STATERR)
		// NOTE
		// . For this level, corresponding disconnect codes are sent
		//   to the client by Genesis.
		logger.Warn("- [Genesis] deafult handler returned false.")
		return
	}
	ch = make(chan *Packet, 1)
	np, err := clbc.ReceiveWithTimeout(time.Second*2, ch)
	if err != nil || np == nil {
		logger.Debug("Handle", "- [Timeout] recv with timeout", err, np == nil)
		clbc.SetStatus(STATERR)
		clbc.client.Disconnected(protobase.PUAckDeadline)
		return
	}
	close(ch)
	clbc.dispatch(np)
	if clbc.GetStatus() != STATONLINE {
		logger.Warn("? [Handle] status!=STATONLINE.")
		clbc.SetStatus(STATERR)
		clbc.client.Disconnected(protobase.PURejected)
		return
	}
	// send first ping packet
	// set timer to routin ping delivery
	// run I/O coroutines
	// finalize setup
	_ = clbc.sendPing()
	// prevent data race by concurrent access
	/* critical section */
	clbc.clock.Lock()
	clbc.pinger = timer.NewTimer(dur)
	// increment work group
	clbc.corous.Add(2)
	go clbc.recvHandler()
	go clbc.uniSendHandler()
	if !clbc.justStarted {
		clbc.AllocateChannels()
		clbc.cendch = make(chan struct{})
		clbc.ShouldTerminate = make(chan struct{}, 1)
	}
	// set this flag to signal fresh start
	// redeliver remaining packets
	clbc.justStarted = true
	go func() {
		// TODO
		// . remove hard-coded duration
		time.Sleep(cDefaultSleepTime)
		clbc.SendRedelivery()
		_ = clbc.client.Connected(clbc.stateOpts[protobase.CCONNACK])
	}()
	clbc.clock.Unlock()
	/* critical section - end */
	// main loop
ML:
	for {
		stat = atomic.LoadUint32(&clbc.Status)
		switch stat {
		case STATERR:
			logger.Debug("- [Error] STATERR.")
			break ML
		case STATGODOWN:
			logger.Debug("* [ForceShutdown] received force shutdown, cleaning up ....")
			stat = STATDISCONNECT
			clbc.SetStatus(STATDISCONNECT)
			clbc.Shutdown()
			break ML
		}
		select {
		case <-clbc.ShouldTerminate:
			logger.Info("* [Terminate] request. Should terminate.")
			if clbc.GetStatus() == STATONLINE {
				if err := clbc.sendDisconnect(); err != nil {
					// TODO
					// . handle error case
					logger.Warn("- [Disconnect] error while sending disconnect packet, continuing ... .")
					stat = STATDISCONNECT
				}
			}
			clbc.SetStatus(STATDISCONNECT)
			clbc.Shutdown()
			break ML
		case packet, ok := <-clbc.RecvChan:
			if !ok {
				logger.Debug("- [RecvChan] cannot fetch packets from recvchan.")
				break ML
			}
			// TODO
			// . reuse session containers
			// . run this concurrently
			logger.Debug("+ [Message] Received .", "userId", clbc.client.GetIdentifier(), "data", packet.Data)
			clbc.dispatch(packet)
		case <-clbc.ErrChan:
			logger.Warn("- [Shit] went down. Panic.")
			break ML
		default:
			time.Sleep(time.Millisecond * 2)
		}
	}
	clbc.pinger.Stop()
	clbc.terminate()
	switch stat {
	case protobase.STATGODOWN, protobase.STATDISCONNECT:
		clbc.client.Disconnected(protobase.PUDisconnect)
	case protobase.STATERR:
		clbc.client.Disconnected(protobase.PUForceTerminate)
	default:
		logger.FDebug("Handle", "- [Handler/State] Unknown state (neither statgodown or staterr)")
	}
	logger.Debug("before corous")
	clbc.corous.Wait()
	logger.Debug("after corous")
}

func (clbc *CLBConnection) SetClient(cl protobase.ClientInterface) {
	clbc.client = cl
}

// ShutDown terminates the connection.
func (clbc *CLBConnection) Shutdown() {
	const fn string = "Shutdown"
	var (
		err error
	)
	err = clbc.Conn.Close()
	if err != nil {
		logger.FWarnf(fn, "- [CLBConnection] unable to close the connection. error:", err)
		return
	}
	logger.Debug("* [Event] Shutting down stream.")
}

// HandleSendError is a error handler. It is used for errors
// caused by sending packets. Currently it terminates the
// connection.
func (clbc *CLBConnection) handleSendError(err error) {
	clbc.Conn.Close()
}

// terminate shuts the connection down and undoes some side effects.
func (clbc *CLBConnection) terminate() {
	clbc.protocon.Lock()

	clbc.protocon.Conn.Close()
	clbc.protocon.Conn = nil
	clbc.protocon.Writer = nil
	clbc.protocon.Reader = nil

	clbc.protocon.Unlock()
}

// dispatch is responsible to call correct methods on state structures.
func (clbc *CLBConnection) dispatch(packet protobase.PacketInterface) {
	switch packet.GetCode() {
	case protobase.PCONNECT:
		clbc.State.OnCONNECT(packet)
	case protobase.PCONNACK:
		clbc.State.OnCONNACK(packet)
	case protobase.PSUBSCRIBE:
		clbc.State.OnSUBSCRIBE(packet)
	case protobase.PSUBACK:
		clbc.State.OnSUBACK(packet)
	case protobase.PPUBLISH:
		clbc.State.OnPUBLISH(packet)
	case protobase.PPUBACK:
		clbc.State.OnPUBACK(packet)
	case protobase.PPING:
		clbc.State.OnPING(packet)
	case protobase.PPONG:
		clbc.State.OnPONG(packet)
	case protobase.PDISCONNECT:
		clbc.State.OnDISCONNECT(packet)
		// NOTE: Rest of protocol data suite should be integrated in this case
	}
}

// NOTE: this is not thread safe, lock must be acquired by the
// caller.
func (clbc *CLBConnection) SetNetConnection(conn net.Conn) {
	if clbc.protocon.Conn != nil {
		clbc.protocon.Conn.Close()
	}
	clbc.protocon.Conn = conn
	clbc.protocon.Writer = bufio.NewWriter(clbc.protocon.Conn)
	clbc.protocon.Reader = bufio.NewReader(clbc.protocon.Conn)
}

func (clbc *CLBConnection) SetMessageStorage(storage protobase.MessageBox) {
	clbc.storage = storage
}

// GetClient returns the responsible struct implementing `ClientInterface`.
func (clbc *CLBConnection) GetClient() protobase.ClientInterface {
	// TODO
	// . add lock
	return clbc.client
}

// recvHandler is the main receive handler.
func (clbc *CLBConnection) recvHandler() {
	const fname = "recvHandler"
	var (
		packet *Packet
		err    error
	)
	defer func() {
		logger.FInfo(fname, "+ [WorkGroup] decrementing workgroup.")
		clbc.corous.Done()
		if err != nil {
			clbc.SetStatus(STATERR)
		}
		// NOTE
		// . this has changed from (send handler, prio handler) to unihandler
		// clbc.cendch <- struct{}{} // send handler
		// clbc.cendch <- struct{}{} // prio handler
		clbc.cendch <- struct{}{} // uni handler
		logger.FInfo(fname, "+ [WorkGroup] terminated other coroutines, returning.")
	}()
	for {
		packet, err = clbc.Receive()
		if err != nil {
			logger.FError(fname, "- [RecvHandler] error while receiving packets.", "error:", err)
			// TODO
			//  handle errors
			return
		}
		clbc.RecvChan <- packet
	}
}

func (clbc *CLBConnection) SendRedelivery() {
	const fn = "SendRedelivery"
	var (
		// clid     string = clbc.GetClient().GetIdentifier()
		outbound []protobase.EDProtocol
		packet   *Packet
		err      error
	)
	logger.FDebug(fn, "+ [Redeliver] starting redelivery.")
	outbound = clbc.storage.GetAllOut()
	logger.FDebug(fn, "[OUTBOUNDS]", outbound)
	for _, p := range outbound {
		if clbc.GetStatus() == protobase.STATONLINE {
			switch p.(type) {
			case *Publish:
				tmp := p.(*Publish)
				tmp.Meta.Dup = true
				/* d e b u g */
				// tmp.Meta.Qos = p.QoS()
				/* d e b u g */
				tmp.Encoded = nil
				err = tmp.Encode()
				if err != nil {
					logger.FWarnf(fn, "- [CLBConnection] unable to encode publish packet. error:", err)
				}
			case *Subscribe:
				logger.Warn("- [sendRedelivery] PACKET TYPE IS [Subscribe]")
			default:
				logger.Warn("- [sendRedelivery] UNKNOWN PACKET TYPE, TODO:")
			}
			packet = p.GetPacket().(*Packet)
			logger.FDebug(fn, "+ [Redliver] undelivered packages are in their path to broker.")
			clbc.Send(packet)
		}
	}
}

func (clbc *CLBConnection) MakeEnvelope(route string, payload []byte, qos byte, messageId uint16, dir protobase.MsgDir) protobase.MsgInterface {
	var (
		box protobase.MsgInterface = protocol.NewMsgBox(qos, messageId, dir, protocol.NewMsgEnvelope(route, payload))
	)
	return box
}

// - MARK: Protocol communication routines section.

func (clbc *CLBConnection) Publish(topic string, message []byte, qos byte, fn func(protobase.OptionInterface, protobase.MsgInterface)) (err error) {
	const _fn string = "Publish"
	var (
		clid    string            = clbc.GetClient().GetIdentifier()
		pb      *protocol.Publish = protocol.NewRawPublish()
		puid    uuid.UUID
		idstore protobase.MSGIDInterface
	)
	logger.FDebugf(_fn, "* [Publish/QoS][CLBConnection] qos is (%b) [Topic] is (%s) [Message] is (%s).", qos, topic, message)
	puid = (pb.Id)
	// set topic and message
	pb.Topic = topic
	pb.Message = message
	// handle quality of service > 0
	if qos > 0 {
		pb.Meta.Qos = qos
		idstore = clbc.storage.GetIDStoreO()
		pb.Meta.MessageId = idstore.GetNewID(puid)
		logger.FDebugf(_fn, "* [Publish<-]CLBConnection] with id (%d).", pb.Meta.MessageId)
		if !clbc.storage.AddOutbound(pb) {
			logger.FWarn(_fn, "- [NOTICE][CLBConnection] unable to add outbound packet to [MessageBox] for userId(%s).", clid)
			return ECLBSendFailure
		}
		clbc.clblock.Lock()
		clbc.clbpub[pb.Meta.MessageId] = fn
		clbc.clblock.Unlock()
	}
	err = pb.Encode()
	if err != nil {
		logger.FWarnf(_fn, "- [CLBConnection] unable to encode publish packet. error:", err)
		return ECLBSendFailure
	}
	packet := pb.GetPacket().(*Packet)
	if clbc.GetStatus() == STATONLINE {
		clbc.Send(packet)
		logger.Infof("* [Publish<-] Publishing [Message](%s) -> [Topic](%s) with [QoS](%b), [MessageId](%d).",
			message, topic, qos, pb.Meta.MessageId)
		if qos == 0 {
			fn(nil, clbc.MakeEnvelope(topic, message, qos, pb.Meta.MessageId, protobase.MDInbound))
		}
	} else {
		// drop packets (QoS == 0)
		logger.FWarn("Publish", "- [Publish] dropping packet due to status.")
		return ECLBSendFailure
	}
	return nil
}

func (clbc *CLBConnection) Subscribe(topic string, qos byte, fn func(protobase.OptionInterface, protobase.MsgInterface)) (err error) {
	const _fn string = "Subscribe"
	var (
		clid    string                   = clbc.GetClient().GetIdentifier() // client identifier
		sb      *Subscribe               = protocol.NewRawSubscribe()       // subscribe packet
		puid    uuid.UUID                                                   // packet UID
		idstore protobase.MSGIDInterface                                    // identifier provider
		packet  *Packet                                                     // empty packet for manipulating
	)
	logger.FDebugf(_fn, "* [Subscribe/QoS] qos is (%b) [Topic] is (%s).", qos, topic)
	puid = (sb.Id)
	// set the topic
	sb.Topic = topic
	if qos > 0 {
		sb.Meta.Qos = qos
		idstore = clbc.storage.GetIDStoreO()
		sb.Meta.MessageId = idstore.GetNewID(puid)
		logger.FDebugf(_fn, "* [Subscribe->][CLBConnection] with id (%d).", sb.Meta.MessageId)
		if !clbc.storage.AddOutbound(sb) {
			logger.FWarn(_fn, "- [NOTICE][CLBConnection] unable to add outbound packet to [MessageBox] for userId(%s).", clid)
			return ECLBSendFailure
		}
		// write to callback map
		clbc.clblock.Lock()
		clbc.clbsub[sb.Meta.MessageId] = fn
		clbc.clblock.Unlock()
	}
	err = sb.Encode()
	if err != nil {
		logger.FWarnf(_fn, "- [CLBConnection] unable to encode subscribe packet. error:", err)
		return ECLBSendFailure
	}
	packet = sb.GetPacket().(*Packet)
	if clbc.GetStatus() == STATONLINE {
		clbc.Send(packet)
		logger.Infof(_fn, "* [Subscribe->] Subscribing to [Topic](%s) with [QoS](%b), [MessageId](%d).",
			topic, qos, sb.Meta.MessageId)
		if qos == 0 {
			fn(nil, clbc.MakeEnvelope(topic, nil, qos, sb.Meta.MessageId, protobase.MDInbound))
		}
	} else {
		// drop packets (QoS == 0)
		logger.FWarn(_fn, "- [Subscribe] dropping packet due to status.")
		return ECLBSendFailure
	}
	return nil
}

func (clbc *CLBConnection) Queue(action protobase.QAction, address string, returnPath string, mark []byte, message []byte) (err error) {
	// TODO
	const _fn string = "Queue"
	logger.FDebugf(_fn, "* [Queue][CLBConnection] invoking with Address(%s), ReturnPath(%s), Mark(%s), Message(%s).",
		address, returnPath, string(mark), string(message))
	var (
		q *protocol.Queue = protocol.NewQueue()
		p *Packet
	)
	q.Action = action
	q.Address = address
	q.ReturnPath = returnPath
	q.Mark = mark
	q.Message = message
	q.Meta.MessageId = 1
	if err = q.Encode(); err != nil {
		return err
	}
	p = q.GetPacket().(*Packet)
	if clbc.GetStatus() == STATONLINE {
		clbc.Send(p)
		logger.Infof("* [Queue->] Sending Queue request to Address(%s) with ReturnPath(%s).", address, returnPath)
	} else {
		return ECLBSendFailure
	}
	return nil
}

func (clbc *CLBConnection) HandleDefault(packet protobase.PacketInterface) (status bool) {
	const fn string = "HandleDefault"
	logger.Infof(fn, "* [CLBConnection] is not implemented.")
	return false
}

func (clbc *CLBConnection) Run() {
	const fn string = "Run"
	logger.Infof(fn, "* [CLBConnection] is not implemented.")
}

func (clbc *CLBConnection) SetNextState() {
	const fn string = "SetNextState"
	logger.Infof(fn, "* [CLBConnection] is not implemented.")
}
