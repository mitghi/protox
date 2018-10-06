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
  _ protobase.ConnectionState = (*constate)(nil)  
)

// constatebase is base struct for constate.
type constatebase struct {
  Conn protobase.BaseControlInterface
}

// constate implements state receiver methods.
type constate struct {
	constatebase

	client protobase.ClientInterface
	server protobase.ServerInterface
}

func newconstatebase(conn protobase.BaseControlInterface) *constatebase {
	result := &constatebase{
		Conn: conn,
	}
	return result
}

// SetNextState pushes into next state.
func (csb *constatebase) SetNextState() {
}

// Run is a receiver that is invoked after state creation
// and before before real execution.
func (csb *constatebase) Run() {

}

// Handle is the main handler ( stub for constate ).
func (csb *constatebase) Handle(packet protobase.PacketInterface) {
}

// HandleDefault is the default handler ( stub for constate ).
func (csb *constatebase) HandleDefault(packet protobase.PacketInterface) (status bool) {
	return true
}

// onCONNECT is not valid in this stage.
func (csb *constatebase) OnCONNECT(packet protobase.PacketInterface) {
	// TODO
	csb.Conn.Shutdown()
}

// onCONNACK is not valid in this stage.
func (csb *constatebase) OnCONNACK(packet protobase.PacketInterface) {
	// TODO
	csb.Conn.Shutdown()
}

// onPUBLISH is the handler for `Publish` packets.
func (csb *constatebase) OnPUBLISH(packet protobase.PacketInterface) {
	logger.Debug("+ [Publish] onPublish.")
	csb.Conn.Shutdown()
}

// onPUBACK is a handler which removes the outbound publish
// message when QoS >0.
func (csb *constatebase) OnPUBACK(packet protobase.PacketInterface) {
	// TODO
	logger.FDebug("onPUBACK", " puback received.")
	csb.Conn.Shutdown()
}

// onSUBSCRIBE is a STUB ( for Genesis ).
func (csb *constatebase) OnSUBSCRIBE(packet protobase.PacketInterface) {
	logger.FDebug("onSUBSCRIBE", " in constatebase.")
	csb.Conn.Shutdown()
	// csb.Conn.Shutdown()
}

// onSUBACK is a handler which removes the outbound subscribe
// message when QoS >0.
func (csb *constatebase) OnSUBACK(packet protobase.PacketInterface) {
	// TODO
	logger.FDebug("onSUBACK", "suback received.")
	csb.Conn.Shutdown()
}

// onPING is the heartbeat handler ( other packets reset its timer as well ).
func (csb *constatebase) OnPING(packet protobase.PacketInterface) {
	logger.Debug("+ [Heartbeat] Received.")
	csb.Conn.Shutdown()
}

func (csb *constatebase) OnPONG(packet protobase.PacketInterface) {
	logger.Debug("+ [Pong] received")
}

//
func (csb *constatebase) OnDISCONNECT(packet protobase.PacketInterface) {
	// TODO
	logger.FDebug("onDISCONNECT", " disconnect received.")
	csb.Conn.Shutdown()
	// csb.Conn.Shutdown()
}

// onQueue is a handler for Queue control packets.
func (csb *constatebase) OnQUEUE(packet protobase.PacketInterface) {
  const fn string = "OnQueue"
	logger.FDebug(fn, " in constatebase.")
	csb.Conn.Shutdown()
}

// onQueueAck is a handler for Queue acknowledge packets.
func (csb *constatebase) OnQueueAck(packet protobase.PacketInterface) {
  const fn string = "OnQueueAck"
	logger.FDebug(fn, " in constatebase.")
	csb.Conn.Shutdown()
}

// SetServer sets the caller ( server).
func (csb *constatebase) SetServer(server protobase.ServerInterface) {
}

// SetClient sets the internal client struct pointer.
func (csb *constatebase) SetClient(client protobase.ClientInterface) {
}

func (csb *constatebase) Shutdown() {
}

func newconstate(conn protobase.ConnectionState) *constate {
	result := &constate{
		constatebase: constatebase{
			Conn: conn,
		},
		client: nil,
		server: nil,
	}

	return result
}

// SetServer sets the caller ( server).
func (cs *constate) SetServer(server protobase.ServerInterface) {
	cs.server = server
}

// SetClient sets the internal client struct pointer.
func (cs *constate) SetClient(client protobase.ClientInterface) {
	cs.client = client
}

// Shutdown sets the status to error which notifies the supervisor
// and cleanly terminates the connection.
func (cs *constate) Shutdown() {
	// logger.Debug("* [Genesis] Closing.")
	// TODO
	// atomic.StoreUint32(&(cs).Status, STATERR)
	cs.client.Disconnected(protobase.PUNone)
}
