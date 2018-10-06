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
	"sync"
	"sync/atomic"
	"time"

	"github.com/mitghi/protox/protobase"
)

type cstflg byte

const (
	CLConnected cstflg = iota
	CLDisconnected
	CLSent
	CLRecv
	CLReject
	CLFault
)

type connstats struct {
	sent       uint64
	recv       uint64
	connect    uint64
	disconnect uint64
	reject     uint64
	fault      uint64
}

type serverState struct {
	sync.RWMutex

	clients map[string]*connection
	mode    byte
	// TODO
	// conns   map[net.Conn]*connection
}

type connection struct {
	sync.RWMutex
	connstats
	conninfo

	typ        SConnType
	conn       *net.Conn
	uid        string
	ip         string
	hasSession bool
	persist    bool
	// TODO
	//  add callbacks
}

type conninfo struct {
	start    *time.Time
	end      *time.Time
	proto    protobase.ProtoConnection
	client   protobase.ClientInterface
	msgstore protobase.MessageStorage
	authsys  protobase.AuthInterface
}

func newConnection(typ SConnType, uid string, conn net.Conn, session bool, persist bool) *connection {
	return &connection{
		conn:       &conn,
		uid:        uid,
		hasSession: session,
		persist:    persist,
		// ip:         conn.RemoteAddr().String(),
	}
}

func newServerState(mode byte) *serverState {
	ret := &serverState{
		clients: make(map[string]*connection),
		mode:    mode,
		// conns:   make(map[net.Conn]*connection),
	}
	return ret
}

func (c *connection) Started() {
	t := time.Now()
	c.ip = (*c.conn).RemoteAddr().String()
	c.conninfo.start = &t
	c.conninfo.end = nil
}

func (c *connection) Ended() {
	t := time.Now()
	c.conninfo.end = &t
}

func (c *connection) setInfo(conn net.Conn, proto protobase.ProtoConnection, client protobase.ClientInterface, msgstore protobase.MessageStorage, authsys protobase.AuthInterface) {
	c.conn = &conn
	c.conninfo.proto = proto
	c.conninfo.client = client
	c.conninfo.msgstore = msgstore
	c.conninfo.authsys = authsys
}

func (c *connection) Inc(statics cstflg) {
	switch statics {
	case CLConnected:
		atomic.AddUint64(&c.connect, 1)
	case CLDisconnected:
		atomic.AddUint64(&c.disconnect, 1)
	case CLSent:
		atomic.AddUint64(&c.sent, 1)
	case CLRecv:
		atomic.AddUint64(&c.recv, 1)
	case CLReject:
		atomic.AddUint64(&c.reject, 1)
	case CLFault:
		atomic.AddUint64(&c.fault, 1)
	default:
		// invalid
	}
}

func (c *connection) update() {
	t := time.Now()
	c.start = &t
	c.ip = (*c.conn).RemoteAddr().String()
}

func (c *connection) Status() uint32 {
	return c.proto.GetStatus()
}

func (c *connection) Statics() (sent, recv, connect, disconnect, reject, fault uint64) {
	sent = atomic.LoadUint64(&c.sent)
	recv = atomic.LoadUint64(&c.recv)
	connect = atomic.LoadUint64(&c.connect)
	disconnect = atomic.LoadUint64(&c.disconnect)
	reject = atomic.LoadUint64(&c.reject)
	fault = atomic.LoadUint64(&c.fault)
	return sent, recv, connect, disconnect, reject, fault
}

func (c *connection) ResetStatics() {
	c.Lock()
	c.sent = 0
	c.recv = 0
	c.connect = 0
	c.disconnect = 0
	c.reject = 0
	c.fault = 0
	c.Unlock()
}

func (s *serverState) get(cid string) (val *connection) {
	s.RLock()
	defer s.RUnlock()
	if val, ok := s.clients[cid]; ok {
		return val
	}
	return nil
}

func (s *serverState) set(cid string, info *connection) {
	s.Lock()
	s.clients[cid] = info
	s.Unlock()
}

func (s *serverState) pruneByCid(cid string) (val *connection) {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.clients[cid]; ok {
		delete(s.clients, cid)
		return val
	}
	return nil
}
