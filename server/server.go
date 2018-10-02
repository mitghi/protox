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

package server

import (
	"net"
	"sync/atomic"
	"time"

	buffpool "github.com/mitghi/lfpool"
	"github.com/mitghi/protox/protobase"
)

// Ensure interface (protocol) conformance.
var (
	_ protobase.ServerInterface = (*Server)(nil)
)

// NewServer creates a new server instance. It is the entry point
// for using underlaying subsystems. Most of the subsystems can be
// customized by providing a handler function or delegating.
func NewServer() *Server {
	var (
		server *Server = &Server{}
	)
	server.Clients = make(map[net.Conn]protobase.ProtoConnection)
	server.Router = make(map[string]map[string]protobase.ProtoConnection)
	server.heartbeat = 1
	server.Status = protobase.ServerNone
  server.StatusChan = make(chan uint32, 1)
	server.critical = make(chan struct{}, 1)
	server.buffer = buffpool.NewBuffPool()
	server.rt = NewRouterWithBuffer(server.buffer)
	server.State = newServerState(0)

	return server
}

// NewServerWithConfigs creates a new server instance using options
// given as `opts` argument and returns a pointer to it. It returns
// an error in case of invalid options or unsucces.
func NewServerWithConfigs(opts ServerConfigs) (*Server, error) {
	var (
		s *Server = NewServer()
	)
	s.opts = &opts
	err := precheckOpts(s.opts)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetStatus atomically loads current state of the `Server`.
func (s *Server) GetStatus() uint32 {
	return atomic.LoadUint32(&s.Status)
}

// GetStatusChan returns a channel containing the server status.
// It is used for reterieving initial status.
func (s *Server) GetStatusChan() <-chan uint32 {
  return s.StatusChan
}

// GetErrChan returns a channel to the caller. It sends a packet when server
// gets into a fatal error.
func (s *Server) GetErrChan() <-chan struct{} {
	return s.critical
}

// SetClientHandler sets the delegate to a function with a signature
// of `func(string, string, string) protocol.ClientInterface`. It gets
// called only after the connection is fully authorized and
// passed `Genesis` stage.
func (s *Server) SetClientHandler(fn func(string, string, string) protobase.ClientInterface) {
	s.onNewClient = fn
}

// SetAuthenticator sets the delegate to use the provided handler. It should
// be called prior to running the server and listening for new connections.
func (s *Server) SetAuthenticator(authenticator protobase.AuthInterface) {
	s.Authenticator = authenticator
}

// SetConnectionHandler creates a connection of type `protocol.ProtoConnection`
// which handles all low-level protocol logics.
func (s *Server) SetConnectionHandler(fn ConnectionDelegate) {
	s.onNewConnection = fn
}

// SetPermissionDelegate sets the delegate for a subsystem that checks
// right accesses and permissions.
func (s *Server) SetPermissionDelegate(pd func(protobase.AuthInterface, ...string) bool) {
	s.permissionDelegate = pd
}

// SetHeartBeat sets the maximum tolerable time ( heartbeat ) in which not
// receiving packets from a client does not cause connection termination.
func (s *Server) SetHeartBeat(heartbeat int) {
	s.heartbeat = heartbeat
}

// SetMessageStore sets the internal storage delegate to a `protobase.MessageStorage`
// compatible struct.
func (s *Server) SetMessageStore(store protobase.MessageStorage) {
	s.Store = store
}
// SetLogger is a method that implements `prtobase.ILoggable`.
func (s *Server) SetLogger(l protobase.LoggingInterface) {
	logger = l
}

func (s *Server) Redeliver(prc protobase.ProtoConnection) {
	const fn = "Redeliver"
	var (
		cl       protobase.ClientInterface = prc.GetClient()
		clid     string                    = cl.GetIdentifier()
		outbound []protobase.EDProtocol
	)
	if c := s.State.get(clid); c != nil {
		logger.FDebugf(fn, "+ [Redeliver] Starting packet redelivery for client(%s).", clid)
		outbound = s.Store.GetAllOut(clid)
		for _, p := range outbound {
			logger.FDebugf(fn, "+ [Redeliver] client(%s) has (%+v) packet.", clid, p)
			if prc.GetStatus() == protobase.STATONLINE {
				prc.SendRedelivery(p)
			}
		}
	}
}

// SetMessageHandler is a callback to server whenever a new package is received
// and parsed in the protocol layer.
func (s *Server) SetMessageHandler(fn ServerHandlerFunc) {
	s.onNewMessage = fn
}
// SetStatus sets the internal status to the new argument status and
// returns its old value, atomically.
func (s *Server) SetStatus(status uint32) uint32 {
	return atomic.SwapUint32(&s.Status, status)
}

// Shutdown is a graceful shutdown method which returns a channel.
// It can be polled if server fails to shuts down.
func (s *Server) Shutdown() (<-chan struct{}, error) {
	const fn = "Shutdown"
	stat := s.GetStatus()
	switch stat {
	case protobase.ServerNone, protobase.ServerStopped, protobase.ForceShutdown:
		if stat == protobase.ForceShutdown {
			logger.FDebug(fn, "- [Server] already in force-shutdown state, dismissing request....")
		}
		return nil, SRVShutdownError
	}
	var (
		ch   chan struct{} = make(chan struct{}, 1)
		tick *time.Ticker  = time.NewTicker(time.Millisecond * 500)
	)
	s.SetStatus(protobase.ForceShutdown)
	go func() {
		for _ = range tick.C {
			if stat := s.GetStatus(); stat == protobase.ServerStopped {
				ch <- struct{}{}
				break
			}
		}
		tick.Stop()
	}()

	return ch, nil
}

// NotifyConnected is a delegate routine. It gets called whenever a new client
// successfully passes `Genesis` stage and authorization levels. This function
// is directly called by a compatible `protocol.ProtoConnection` structure and
// only after that the corresponsing function on a compatible `client.ClientInterface`
// will be called. Subsequently, a `client.ClientInterface` struct can drop the
// connection after.
func (s *Server) NotifyConnected(prc protobase.ProtoConnection) {
	const fn = "NotifyDisconnected"
	var (
		cl   protobase.ClientInterface = prc.GetClient()
		clid string                    = cl.GetIdentifier()
		conn net.Conn                  = prc.GetConnection()
		c    *connection
	)

	logger.Infof("+ [Server] Client(%s) passed [Genesis] state and is now [Online].", clid)

	if c = s.State.get(cl.GetIdentifier()); c != nil {
		logger.FDebugf(fn, "* [Client] client (%s) already exists.", clid)
		c.Lock()
		c.setInfo(conn, prc, cl, nil, s.Authenticator)
		c.update()
		c.Inc(CLConnected)
		c.Unlock()

		s.Redeliver(prc)

	} else {
		c = newConnection(STCLIENT, clid, conn, true, true)
		c.setInfo(conn, prc, cl, nil, s.Authenticator)
		c.Inc(CLConnected)
		s.State.set(clid, c)
	}
}

// NotifySubscribe is a delegate routine that registers client subscriptions. It
// is only accessible above `Genesis` stage. Calling this method is theresponsibility
// of `client.ClientInterface` and consequently - implementor may decide to drop
// new subscriptions by not calling this function.
//
// Note: this is the new implementation and is under development.
//
func (s *Server) NotifySubscribe(prc protobase.ProtoConnection, msg protobase.MsgInterface) {
	// TODO
	const _fn = "NotifySubscribe"
	var (
		clid  string = prc.GetClient().GetIdentifier()
		topic string = msg.Envelope().Route()
		qos   byte   = msg.QoS()
	)
	logger.FDebugf(_fn, "+ [Client][Layer] client(%s) attached to stream of (%s) with QoS(%d).", clid, topic, int(qos))
	logger.Infof("+ [Subscription] Client(%s) subscribed to stream (%s) with QoS(%d).", clid, topic, int(qos))
	s.rt.AddSub(clid, topic, qos)
}

// NotifyPublish sends messages from publishers to subscribers. A compatible
// `client.ClientInterface` structure is responsible to call this function and
// may decide not to if messages must be dropped.
// Note: this uses the new implementation and is under development.
func (s *Server) NotifyPublish(prc protobase.ProtoConnection, msg protobase.MsgInterface) {
	const fn = "NotifyPublish"
	var (
		topic   string = msg.Envelope().Route()
		message []byte = msg.Envelope().Payload()
		// dir     protobase.MsgDir = msg.Dir()
	)
	m, _ := s.rt.FindSub(topic)
	for k, wqos := range m {
		cl := s.State.get(k)
		logger.FDebug(fn, "* [Publish] client found.", cl)
		if cl != nil {
			var (
				user protobase.ClientInterface = cl.proto.GetClient()
				clid string                    = cl.uid
			)
			if cl.proto == prc {
				logger.FDebug(fn, "? [Publish] cl.proto==prc ? ", "userId", clid)
				// NOTE: IMPORTANT:
				// . this has changed
				// . remove after test
				// begin
				// continue
				// end
			}
			prclid := prc.GetClient().GetIdentifier()
			// if stat := cl.proto.GetStatus(); stat == protobase.STATONLINE {
			if wqos == 0 && cl.proto.GetStatus() != protobase.STATONLINE {
				continue
			}
			logger.FDebugf(fn, "+ [Publish] prc(%s) - client(%s) is [Online], sending message(%s) from route(%s), QoS(%d).", prclid, clid, message, clid, msg.QoS())

			npb := msg.Clone(protobase.MDOutbound)
			// currQoS := npb.QoS()
			npb.SetWishQoS(wqos)
			// logger.Debugf("+ [Publish     ] Routing Topic(%s)-> Message(%s) for Client(%s) with QoS(%d) [ WishQoS(%d), wqos(%d) ].", topic, message, clid, currQoS, npb.QoS(), wqos)
			logger.Infof("+ [Publish     ] Routing Topic(%s)-> Message(%s) for Client(%s) with QoS(%d).", topic, message, clid, npb.QoS())
			// NOTE
			// . this has changed
			// cl.proto.SendMessage(npb)
			cl.proto.SendMessage(npb, cl.proto == prc)
			user.Publish(npb)
			// } else {
			// 	// TODO
			// 	// . add to outbound messages
			// }
		}
	}
	return
}

func (s *Server) NotifyQueue(prc protobase.ProtoConnection, msg protobase.MsgInterface) {
	const _fname = "NotifyQueue"
	logger.FInfo(_fname, "+ [Queue    ] message received.")
}

// Setup is for prechecks before running the server. It must crash
// ( or recover ) to indicate fatal problems early on.
func (s *Server) Setup() {
	// TODO
	// . improve this
	// . check validity of permissionDelegate
	if s.Authenticator == nil {
		panic("AUTHENTICATOR IS NIL")
	}
	if s.Store == nil {
		panic("STORAGE IS NIL")
	}
	if s.onNewClient == nil {
		panic("ONNEWCLIENT IS NIL")
	}
	if s.onNewConnection == nil {
		panic("ONNEWCONNECTION IS NIL")
	}
	// TODO
	// . fallback to default logger ??
}

// disconnectAll atomically sets the shutdown flag for all `connected` clients.
func (s *Server) disconnectAll() {
	const fn = "disconnectAll"
	switch s.GetStatus() {
	case protobase.ServerNone, protobase.ServerStopped:
		return
	default:
		s.State.RLock()
		for _, v := range s.State.clients {
			if stat := v.Status(); stat == protobase.STATONLINE {
				logger.FDebug(fn, "- [STATCONNECTED] setting godown flag for [CIENT].", "userId", v.uid)
				sent, recv, connect, disconnect, reject, fault := v.Statics()
				logger.FDebug(fn, "- [STATCONNECTED] stats for [CLIENT].", "stats", "userId", v.uid, "stats", sent, recv, connect, disconnect, reject, fault)
				v.proto.SetStatus(protobase.STATGODOWN)
				errch := v.proto.GetErrChan()
				// send notification to client's err chan,
				// drop quitely if chan is closed
				select {
				case errch <- struct{}{}:
				default:
				}
			} else {
				logger.FDebug(fn, "- [Status=%d] client is not connected.", int(stat))
			}
		}
		s.State.RUnlock()
	}
}

// NotifyRejected notifies the server that the connection is
// rejected ( invalid creds, ban , ..... ).
func (s *Server) NotifyReject(prc protobase.ProtoConnection) {
	const fn = "NotifyRejected"
	logger.FDebug(fn, "- [Server] received [Rejection], connection to broker is rejected.")
	logger.Info("- [Server  ] A Client has been rejected.")
	// tell the workgroup that conn is finished
	var (
		cl   protobase.ClientInterface
		clid string
	)
	cl = prc.GetClient()
	if cl == nil {
		logger.FDebug(fn, "? [Server] cl==nil ?")
		s.corous.Done()
		return
	}
	if conn := s.State.get(clid); conn != nil {
		conn.Lock()
		// conn.conn = nil
		conn.Ended()
		conn.Inc(CLReject)
		conn.Unlock()
	}
	s.corous.Done()
	// TODO
}

// NotifyDisconnected notifies the server that the connection to
// a client is either lost or broken.
func (s *Server) NotifyDisconnected(prc protobase.ProtoConnection) {
	const fn = "NotifyDisconnected"
	var (
		cl       protobase.ClientInterface = prc.GetClient()
		clid     string                    = cl.GetIdentifier()
		isGodown bool                      = false
	)
	if stat := s.GetStatus(); stat != protobase.ServerRunning {
		logger.FWarn(fn, "- [SERVER] is not in running state and received a [Disconnect] request from client(%s).", clid)
		s.corous.Done()
		return
	}
	var (
		prcstat uint32 = prc.GetStatus()
	)
	switch prcstat {
	case protobase.STATFATAL:
		// TODO
		//  clean up
		//  error handling
		logger.Warn(fn, "- [Fatal] Internal error.")
		s.corous.Done()
		s.Shutdown()
		return
	case protobase.STATGODOWN:
		logger.FDebugf("? [Server] request to [Disconnect] client(%s).", clid)
		isGodown = true
	}
	logger.FDebugf(fn, "- [Event][Death] Detached [Client](%s).", clid)
	conn := s.State.get(clid)
	if conn == nil {
		// NOTE TODO
		// . this is fatal, recheck
		logger.Fatal("? [Server] conn==nil ?")
	} else {
		conn.Lock()
		logger.FDebugf(fn, "- [Server][Event][ConnEnded] for [Client](%s).", clid)
		conn.conn = nil
		conn.Ended()
		conn.Inc(CLDisconnected)
		conn.Unlock()
		if isGodown {
			cl.Disconnected(protobase.PUDisconnect)
		} else {
			cl.Disconnected(protobase.PUForceTerminate)
		}
	}
	logger.Infof("- [Server  ] Client(%s) disconnected.", clid)
	s.corous.Done()
}

// HandleIncomingConnection uses delegate functions to build and run new connection.
// It also passes all the necessary informations such as certain delegate function to
// a compatible `protocol.ProtoConnection` struct ( made by using delegates ).
func (s *Server) handleIncomingConnection(conn net.Conn) {
	if s.onNewConnection == nil {
		panic("No handler is specified for 'ClientHandler'")
	}
	var newConnection protobase.ProtoConnection = s.onNewConnection(conn)
	newConnection.SetAuthenticator(s.Authenticator)
	newConnection.SetServer(s)
	newConnection.SetClientDelegate(s.onNewClient)
	newConnection.SetMessageStorage(s.Store)
	newConnection.SetHeartBeat(s.heartbeat)
	newConnection.SetPermissionDelegate(s.permissionDelegate)
	// this is for `newConnection.Handle()`
	s.corous.Add(1)
	go newConnection.Handle()
	// Signal that handleIncomingConnection is finished.
	// Receiver is listener of choice which adds this
	// coroutine to the workgroup.
	s.corous.Done()
}

// RegisterClient adds the client to internal data structures is called
// when the client is fully authorized and passes `Genesis` stage.
func (s *Server) RegisterClient(prc protobase.ProtoConnection) {
	const fn = "RegisterClient"
	s.Lock()
	logger.FDebug(fn, "+ [Client] Passed [Genesis] state and is now [Online].")
	s.Clients[prc.GetConnection()] = prc
	// s.State.set
	s.Unlock()
}

func (s *Server) Serve() error {
	// TODO
	s.Setup()
	if s.opts == nil {
		// fallback to tcp server
		// NOTE
		// . this is for test, remove hard coded addr.
		return s.ServeTCP(":52909")
	} else {
		// TODO
		// . run server based on options
	}
	return nil
}

func precheckOpts(opts *ServerConfigs) error {
	// TODO
	// . precheck server options before proceeding
	_, _, err := net.SplitHostPort(opts.Addr)
	if err != nil {
		return SRVInvalidAddr
	}
	switch opts.Mode {
	case ProtoTCP:
		// NOTE:
		// . temporarily its ok
		return nil
	case ProtoTLS:
		if _, ok := opts.Config.(TCPOptions); !ok {
			return SRVMissingOptions
		}
		return nil
	case ProtoSSL:
		// TODO
		return SRVInvalidMode
	case ProtoUNIXSO:
		// TODO
		return SRVInvalidMode
	default:
		return SRVInvalidMode
	}
}
