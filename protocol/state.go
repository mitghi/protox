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

import "github.com/mitghi/protox/protobase"

var _ ConStateInterface = (*constatebase)(nil)
var _ ConnectionState = (*constate)(nil)

type constatebase struct {
	Conn baseControlInterface
}

type constate struct {
	constatebase

	client protobase.ClientInterface
	server protobase.ServerInterface
}

func newconstatebase(conn baseControlInterface) *constatebase {
	result := &constatebase{
		Conn: conn,
	}
	return result
}

// SetNextState pushes into next state.
func (self *constatebase) SetNextState() {
}

// Handle is the main handler ( stub for constate ).
func (self *constatebase) Handle(packet *Packet) {
}

// Run is a receiver that is invoked after state creation
// and before before real execution.
func (self *constatebase) Run() {

}

// HandleDefault is the default handler ( stub for constate ).
func (self *constatebase) HandleDefault(packet *Packet) (status bool) {
	return true
}

// onCONNECT is not valid in this stage.
func (self *constatebase) onCONNECT(packet *Packet) {
	// TODO
	self.Conn.Shutdown()
}

// onCONNACK is not valid in this stage.
func (self *constatebase) onCONNACK(packet *Packet) {
	// TODO
	self.Conn.Shutdown()
}

// onPUBLISH is the handler for `Publish` packets.
func (self *constatebase) onPUBLISH(packet *Packet) {
	logger.Debug("+ [Publish] onPublish.")
	self.Conn.Shutdown()
}

// onPING is the heartbeat handler ( other packets reset its timer as well ).
func (self *constatebase) onPING(packet *Packet) {
	logger.Debug("+ [Heartbeat] Received.")
	self.Conn.Shutdown()
}

func (self *constatebase) onPONG(packet *Packet) {
	logger.Debug("+ [Pong] received")
}

// onSUBACK is a handler which removes the outbound subscribe
// message when QoS >0.
func (self *constatebase) onSUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onSUBACK", "suback received.")
	self.Conn.Shutdown()
}

// onPUBACK is a handler which removes the outbound publish
// message when QoS >0.
func (self *constatebase) onPUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onPUBACK", " puback received.")
	self.Conn.Shutdown()
}

//
func (self *constatebase) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", " disconnect received.")
	self.Conn.Shutdown()
	// self.Conn.Shutdown()
}

// onSUBSCRIBE is a STUB ( for Genesis ).
func (self *constatebase) onSUBSCRIBE(packet *Packet) {
	logger.FDebug("onSUBSCRIBE", " in constatebase.")
	self.Conn.Shutdown()
	// self.Conn.Shutdown()
}

// onQueue is a handler for Queue control packets.
func (self *constatebase) onQUEUE(packet *Packet) {
	logger.FDebug("onQUEUE", " in constatebase.")
	self.Conn.Shutdown()
}

// onQueueAck is a handler for Queue acknowledge packets.
func (self *constatebase) onQueueAck(packet *Packet) {
	logger.FDebug("onQueueAck", " in constatebase.")
	self.Conn.Shutdown()
}

func newconstate(conn baseControlInterface) *constate {
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
func (self *constate) SetServer(server protobase.ServerInterface) {
	self.server = server
}

// SetClient sets the internal client struct pointer.
func (self *constate) SetClient(client protobase.ClientInterface) {
	self.client = client
}

// Shutdown sets the status to error which notifies the supervisor
// and cleanly terminates the connection.
func (self *constate) Shutdown() {
	// logger.Debug("* [Genesis] Closing.")
	// TODO
	// atomic.StoreUint32(&(self.Conn).Status, STATERR)
	self.client.Disconnected(protobase.PUNone)
}
