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
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/mitghi/protox/protobase"
)

type COnline struct {
	constate

	Conn *CLBConnection
}

// MARK: Online

// NewOnline returns a pointer to a new `Online` struct. This is the where
// interactions with a connected/authorized client happens.
func NewCOnline(conn *CLBConnection) *COnline {
	result := &COnline{
		constate: constate{
			constatebase: constatebase{
				Conn: conn,
			},
			client: nil,
			server: nil,
		},
		Conn: conn,
	}
	result.client = conn.GetClient()
	return result
}

// HandleDefault is the default handler ( stub for COnline ).
func (self *COnline) HandleDefault(packet *Packet) (status bool) {
	return true
}

// Shutdown sets the status to error which notifies the supervisor
// and cleanly terminates the connection.
func (self *COnline) Shutdown() {
	logger.Debug("* [Genesis] Closing.")
	atomic.StoreUint32(&(self.Conn).Status, STATERR)
}

// onCONNECT is not valid in this stage.
func (self *COnline) onCONNECT(packet *Packet) {
	// TODO
	// NOTE
	// . this is new
	self.Shutdown()
	self.Conn.protocon.Conn.Close()
}

// onCONNACK is not valid in this stage.
func (self *COnline) onCONNACK(packet *Packet) {
	// TODO
}

// onPUBLISH is the handler for `Publish` packets.
func (self *COnline) onPUBLISH(packet *Packet) {
	var publish *Publish = NewPublish()
	if err := publish.DecodeFrom(packet.Data); err != nil {
		logger.Debug("- [DecodeErr(onPublish)] Unable to decode data.", err)
		self.Shutdown()
		return
	}
	if stat := self.Conn.storage.AddInbound(publish); stat == false {
		logger.Debug("? [NOTICE] addinbound returned false (conline/publish).")
	}
	var puback *Puback = NewPuback()
	if publish.Meta.Qos > 0 {
		puback.Meta.Qos, puback.Meta.MessageId = publish.Meta.Qos, publish.Meta.MessageId
		if err := puback.Encode(); err != nil {
			logger.FError("onPUBLISH", "- [CONLINE] Error while encoding puback.")
			self.Shutdown()
			return
		}
		logger.FTracef(1, "onPUBLISH", "* [QoS] packet QoS(%b) Duplicate(%t) MessageID(%d).", publish.Meta.Qos, publish.Meta.Dup, int(publish.Meta.MessageId))
		var pckt *Packet = puback.GetPacket().(*Packet)
		logger.FTrace(1, "onPUBLISH", "* [PubAck] sending packet with content", pckt.Data)
		// NOTE
		// . this has changed
		// self.Conn.Send(pckt)
		self.Conn.SendPrio(pckt)
		if stat := self.Conn.storage.DeleteIn(publish); stat == false {
			logger.Debug("? [NOTICE] deleteinbound returned false (conline/publish).")
		}
	}
	pb := NewMsgBox(publish.Meta.Qos, publish.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
	// publish box clone
	pbc := pb.Clone(protobase.MDInbound)

	// NOTE
	// . this has changed
	// if publish.Meta.Qos > 0 {
	// 	self.Conn.clblock.Lock()
	// 	// because publish is received when a subscribtion
	// 	// for rotue exists.
	// 	callback, ok := self.Conn.clbsub[puback.Meta.MessageId]
	// 	if ok {
	// 		delete(self.Conn.clbsub, puback.Meta.MessageId)
	// 	}
	// 	self.Conn.clblock.Unlock()

	// 	if ok && callback != nil {
	// 		go func() {
	// 			callback(nil, pbc)
	// 			self.client.Publish(pbc)
	// 		}()
	// 		return
	// 	}
	// }

	// if publish.Meta.Qos > 0 {
	// 	if ok && callback != nil {
	// 		go func() {
	// 			self.client.Publish(pbc)
	// 		}()
	// 		return
	// 	}
	// }

	// TODO
	// . send this to a worker thread
	// go func() { self.client.Publish(pbc) }()
	self.client.Publish(pbc)

}

// onSUBSCRIBE is the handler for `Subscribe` packets.
func (self *COnline) onSUBSCRIBE(packet *Packet) {
	logger.Debug("* [Subscribe] packet is received.")
	// subscribe := NewSubscribe()
	// if err := subscribe.DecodeFrom(packet.Data); err != nil {
	// 	self.Shutdown()
	// 	return
	// }
	// pb := NewMsgBox(subscribe.Meta.Qos, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
	// self.client.Subscribe(pb)
	// self.server.NotifySubscribe(self.Conn, pb)
	// self.client.Subscribe(subscribe.Topic)
	// self.server.NotifySubscribe(subscribe.Topic, self.Conn)
}

// onPING is the heartbeat handler ( other packets reset its timer as well ).
func (self *COnline) onPING(packet *Packet) {
	logger.Debug("+ [Heartbeat] Received.")
}

// onSUBACK is a handler which removes the outbound subscribe
// message when QoS >0.
func (self *COnline) onSUBACK(packet *Packet) {
	// TODO
	var (
		pa  *Suback = NewSuback()
		uid uuid.UUID
	)
	logger.FDebug("onSUBACK", "* [SubAck] packet is received.")
	if err := pa.DecodeFrom(packet.Data); err != nil {
		logger.FDebug("onSUBACK", "- [Decode] uanble to decode in [SubAck].", err)
		return
	}

	oidstore := self.Conn.storage.GetIDStoreO()
	msgid := pa.Meta.MessageId
	uid, ok := oidstore.GetUUID(msgid)
	if !ok {
		logger.FWarn("onSUBACK", "- [IDStore/Suback] no packet with msgid found.", "msgid", msgid)
		return
	}
	np, ok := self.Conn.storage.GetOutbound(uid)
	if !ok {
		logger.FWarn("onSUBACK", "- [MessageBox/Suback] no packet with uid found.", uid)
	}
	if !self.Conn.storage.DeleteOut(np) {
		logger.FWarn("onSUBACK", "- [MessageBox/Suback] failed to remove message.")
	}
	oidstore.FreeId(msgid)

	npc := np.(*Subscribe)
	if npc == nil {
		// TODO
		// . handle this case
		logger.FWarn("onSUBACK", "- [MessageBox/Suback] npc==nil [FATAL].")
	}

	pb := NewMsgBox(npc.Meta.Qos, npc.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(npc.Topic, nil))
	pbc := pb.Clone(protobase.MDInbound)
	self.Conn.clblock.Lock()
	callback, ok := self.Conn.clbsub[msgid]
	if ok {
		delete(self.Conn.clbsub, msgid)
	}
	self.Conn.clblock.Unlock()

	if ok && callback != nil {
		// go func() {
		callback(nil, pbc)
		self.client.Subscribe(pbc)
		// }()
		return
	}

	// TODO
	// . send this to a worker thread
	// go func() { self.client.Subscribe(pbc) }()
	// self.client.Subscribe(pbc)

}

// onPUBACK is a handler which removes the outbound publish
// message when QoS >0.
func (self *COnline) onPUBACK(packet *Packet) {
	// TODO
	var (
		pa  *Puback = NewPuback()
		uid uuid.UUID
	)

	logger.FDebug("onPUBACK", "+ [PubAck] packet received.")
	if err := pa.DecodeFrom(packet.Data); err != nil {
		logger.FDebug("onPUBACK", "- [Decode] uanble to decode in [PubAck].", err)
		return
	}

	oidstore := self.Conn.storage.GetIDStoreO()
	msgid := pa.Meta.MessageId
	uid, ok := oidstore.GetUUID(msgid)
	if !ok {
		logger.FWarn("onPUBACK", "- [IDStore/Puback] no packet with msgid found.", "msgid", msgid)
		return
	}
	np, ok := self.Conn.storage.GetOutbound(uid)
	if !ok {
		logger.FWarn("onPUBACK", "- [MessageBox/Puback] no packet with uid found.", uid)
	}
	if !self.Conn.storage.DeleteOut(np) {
		logger.FWarn("onPUBACK", "- [MessageBox/Puback] failed to remove message.")
	}
	oidstore.FreeId(msgid)

	npc := np.(*Publish)
	if npc == nil {
		// TODO
		// . handle this case
		logger.FWarn("onPUBACK", "- [MessageBox/Puback] npc==nil [FATAL].")
	}

	pb := NewMsgBox(npc.Meta.Qos, npc.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(npc.Topic, npc.Message))
	self.Conn.clblock.Lock()
	callback, ok := self.Conn.clbpub[msgid]
	if ok {
		delete(self.Conn.clbpub, msgid)
	}
	self.Conn.clblock.Unlock()

	if ok && callback != nil {
		// go func() {
		callback(nil, pb)
		// NOTE
		// . this has changed
		// self.client.Publish(pb)
		// }()
		return
	}
	// NOTE
	// . this has changed
	// self.client.Publish(pb)
	// go func() { self.client.Publish(pb) }()
}

func (self *COnline) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", "+ [Disconnect] packet is received.")
	self.Conn.protocon.Conn.Close()
}

func (self *COnline) onPONG(packet *Packet) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}

func (self *COnline) onQueueAck(packet *Packet) {
	logger.FDebug("onQueueAck", "+ [Queue(Ack) packet is received.]")
}
