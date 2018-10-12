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
	"github.com/mitghi/protox/protobase"
)

// Ensure protocol (interface) conformance.
var (
	_ protobase.ConStateInterface = (*constatebase)(nil)
	_ protobase.ConnectionState   = (*constate)(nil)
)

// constatebase is base struct version of constate.
// It covers proto packet handler and other methods
// required by 'protobase.BaseControlInterface' that
// other structs can use to embed and extend its
// methods.
type constatebase struct {
	Conn protobase.BaseControlInterface
}

// constate is a base struct containing valid inherited
// methods fulfilling requirements for post 'Genesis'
// state.
type constate struct {
	constatebase

	client protobase.ClientInterface
	server protobase.ServerInterface
}

// Section: constructors

// newconstatebase returns a pointer to a new instance
// of 'constatebase'. It embeds conn into its public
// field.
func newconstatebase(conn protobase.BaseControlInterface) *constatebase {
	return &constatebase{
		Conn: conn,
	}
}

// newconstate returns a pointer to a new instance
// of 'constate'. It embeds conn into its public
// field and extends its methods.
func newconstate(conn protobase.ConnectionState) *constate {
	return &constate{
		constatebase: constatebase{
			Conn: conn,
		},
		client: nil,
		server: nil,
	}
}

// Section: constatebase methods.

// SetNextState pushes into next state.
// NOTE: empty method
func (csb *constatebase) SetNextState() {
	// NOP
}

// SetServer sets server instance.
// NOTE: empty method
func (csb *constatebase) SetServer(server protobase.ServerInterface) {
}

// SetClient sets client instance.
// NOTE: empty method
func (csb *constatebase) SetClient(client protobase.ClientInterface) {
}

// Run is main method containing business logic
// NOTE: empty method
func (csb *constatebase) Run() {
	// NOP
}

// Handle is blocking method and main loop.
// NOTE: empty method
func (csb *constatebase) Handle(packet protobase.PacketInterface) {
	// NOP
}

// HandleDefault is non-blocking entry method; just
// valid in 'Genesis' stage.
// NOTE: empty method
func (csb *constatebase) HandleDefault(packet protobase.PacketInterface) (status bool) {
	return true
}

// Shutdown is shutdown method.
func (csb *constatebase) Shutdown() {
}

// Section: proto event handler 'protobase.ProtoEventInterface' [ event-handler proto-event ]

// OnCONNECT handles 'Connection' packet.
// NOTE: empty method
func (csb *constatebase) OnCONNECT(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Connect.")
	csb.Conn.Shutdown()
}

// OnCONNACK handles 'Connack' packet.
// NOTE: empty method
func (csb *constatebase) OnCONNACK(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Connack.")
	csb.Conn.Shutdown()
}

// OnPUBLISH handles 'Publish' packet.
// NOTE: empty method
func (csb *constatebase) OnPUBLISH(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Publish.")
	csb.Conn.Shutdown()
}

// OnPUBACK handles 'Puback' packet.
// NOTE: empty method
func (csb *constatebase) OnPUBACK(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Puback.")
	csb.Conn.Shutdown()
}

// OnSUBSCRIBE handles 'Subscribe' packet.
// NOTE: empty method
func (csb *constatebase) OnSUBSCRIBE(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Subscribe.")
	csb.Conn.Shutdown()
}

// OnSUBACK handles 'Suback' packet.
// NOTE: empty method
func (csb *constatebase) OnSUBACK(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Suback.")
	csb.Conn.Shutdown()
}

// OnPING handles 'Ping' packet.
// NOTE: empty method
func (csb *constatebase) OnPING(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Ping.")
	csb.Conn.Shutdown()
}

// OnPong handles 'Pong' packet.
// NOTE: empty method
func (csb *constatebase) OnPONG(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Pong.")
}

// OnDisconnect handles 'Disconnect' packet.
// NOTE: empty method
func (csb *constatebase) OnDISCONNECT(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Disconnect.")
	csb.Conn.Shutdown()
}

// OnQueue handles 'Queue' control packets.
// NOTE: empty method
func (csb *constatebase) OnQUEUE(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] Queue.")
	csb.Conn.Shutdown()
}

// OnQueueAck handles 'QueueAck' packet.
// NOTE: empty method
func (csb *constatebase) OnQueueAck(packet protobase.PacketInterface) {
	logger.Debug("+ [constatebase] QueueAck.")
	csb.Conn.Shutdown()
}

/* Section: constate methods */

// SetServer sets server instance.
func (cs *constate) SetServer(server protobase.ServerInterface) {
	cs.server = server
}

// SetClient sets client instance.
func (cs *constate) SetClient(client protobase.ClientInterface) {
	cs.client = client
}

// Shutdown is shutdown method.
func (cs *constate) Shutdown() {
	cs.client.Disconnected(protobase.PUNone)
}
