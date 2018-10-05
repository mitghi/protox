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

// Online is the second stage. A connection can only be upgraded to `Online` iff it passes
// `Genesis` stage which means it must be fully authorized, valid and compatible with the
// broker.
type Online struct {
	constate

	Conn *Connection
}

// MARK: Online

// NewOnline returns a pointer to a new `Online` struct. This is the where
// interactions with a connected/authorized client happens.
func NewOnline(conn *Connection) *Online {
	result := &Online{
		constate: constate{
			constatebase: constatebase{
				Conn: conn,
			},
			client: nil,
			server: nil,
		},
		Conn: conn,
	}
	return result
}

// HandleDefault is the default handler ( stub for Online ).
func (o *Online) HandleDefault(packet *Packet) (status bool) {
	return true
}

// Shutdown sets the status to error which notifies the supervisor
// and cleanly terminates the connection.
func (o *Online) Shutdown() {
	logger.Debug("* [Genesis] Closing.")
	atomic.StoreUint32(&(o.Conn).Status, STATERR)
	o.client.Disconnected(protobase.PUForceTerminate)
}

// onCONNECT is not valid in this stage.
func (o *Online) onCONNECT(packet *Packet) {
	// TODO
}

// onCONNACK is not valid in this stage.
func (o *Online) onCONNACK(packet *Packet) {
	// TODO
}

// onPUBLISH is the handler for `Publish` packets.
func (o *Online) onPUBLISH(packet *Packet) {
	var (
		publish  *Publish = NewPublish()
		cid      string   = o.client.GetIdentifier()
		userType protobase.AuthUserType
	)
	if err := publish.DecodeFrom(packet.Data); err != nil {
		logger.Debugf("- [DecodeErr(onPublish)] Unable to decode data for Client(%s).", err, cid)
		o.Shutdown()
		return
	}
	userType, err := o.Conn.auth.GetUserType(cid)
	if err != nil {
		logger.Debugf("onPUBLISH", "- [Packet] unable to find associated User Type for Client(%s).", cid)
		return
	}
	if o.Conn.permissionDelegate != nil {
		if !o.Conn.permissionDelegate(o.Conn.auth, "can", "publish", publish.Topic) {
			logger.Debugf("onPUBLISH", "- [Packet] unable to find corresponding permission for Client(%s).", cid)
			o.Shutdown()
			return
		}
	} else {
		role := o.Conn.auth.GetACL().GetRole((string)(userType))
		if role == nil {
			logger.Debug("onPUBLISH", "- [Role] role==nil.")
			o.Shutdown()
			return
		}
    // TDOO
    // . refactor hard-coded permissions
		if !role.HasPerm("can", "publish", publish.Topic) {
			logger.Debugf("onPUBLISH", "- [Packet] unable to find corresponding permission ( direct ) for Client(%s).", cid)
			o.Shutdown()
			return
		}
	}
	if stat := o.Conn.storage.AddInbound(cid, publish); stat == false {
		logger.Debug("? [NOTICE] addinbound returned false (online/publish).")
	}
	var puback *Puback = NewPuback()
	logger.FDebugf("onPUBLISH", "+ [Packet] received with [QoS] %d.", int(publish.Meta.Qos))
	if publish.Meta.Qos > 0 {
		puback.Meta.Qos, puback.Meta.MessageId = publish.Meta.Qos, publish.Meta.MessageId
		if puback.Meta.Qos > MAXQoS {
			puback.Meta.Qos = MAXQoS
		}
		if err := puback.Encode(); err != nil {
			logger.FError("onPUBLISH", "- [ONLINE] Error while encoding puback.")
			o.Shutdown()
			return
		}
		var pckt *Packet = puback.GetPacket().(*Packet)
		o.Conn.SendPrio(pckt)
		if stat := o.Conn.storage.DeleteIn(cid, publish); stat == false {
			logger.Debug("? [NOTICE] deleteinbound returned false (online/publish).")
		}
	}
	pb := NewMsgBox(publish.Meta.Qos, publish.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
	// publish box clone
	pbc := pb.Clone(protobase.MDInbound)
	o.client.Publish(pbc)
	o.server.NotifyPublish(o.Conn, pb)
}

// onSUBSCRIBE is the handler for `Subscribe` packets.
func (o *Online) onSUBSCRIBE(packet *Packet) {
	var (
		subscribe *Subscribe = NewSubscribe()
		cid       string     = o.client.GetIdentifier()
		userType  protobase.AuthUserType
	)
	if err := subscribe.DecodeFrom(packet.Data); err != nil {
		logger.Debugf("- [DecodeErr(onSubscribe)] Unable to decode data for Client(%s).", err, cid)
		o.Shutdown()
		return
	}
	userType, err := o.Conn.auth.GetUserType(cid)
	if err != nil {
		logger.Debugf("onSUBSCRIBE", "- [Packet] unable to find associated User Type for Client(%s).", cid)
		o.Shutdown()
		return
	}
  if o.Conn.auth.GetMode() != protobase.AUTHModeNone {
    if o.Conn.permissionDelegate != nil {
      if !o.Conn.permissionDelegate(o.Conn.auth, "can", "subscribe", subscribe.Topic) {
        o.Shutdown()
        return
      }
    } else {    
      role := o.Conn.auth.GetACL().GetRole((string)(userType))
      if role == nil {
        logger.Debug("onSUBSCRIBE", "- [Role] role==nil.")
        o.Shutdown()
        return
      }
      if !role.HasPerm("can", "subscribe", subscribe.Topic) {
        o.Shutdown()
        return
      }
    }
  }
	if stat := o.Conn.storage.AddInbound(cid, subscribe); stat == false {
		logger.Debug("? [NOTICE] addinbound returned false (online/subscribe).")
	}
	var suback *Suback = NewSuback()
	logger.FDebugf("onSUBSCRIBE", "+ [Packet] received with [QoS] %d.", int(subscribe.Meta.Qos))
	if subscribe.Meta.Qos > 0 {
		suback.Meta.Qos, suback.Meta.MessageId = subscribe.Meta.Qos, subscribe.Meta.MessageId
		if suback.Meta.Qos > MAXQoS {
			suback.Meta.Qos = MAXQoS
		}
		if err := suback.Encode(); err != nil {
			logger.FError("onSUBSCRIBE", "- [ONLINE] Error while encoding suback.")
			o.Shutdown()
			return
		}
		var pckt *Packet = suback.GetPacket().(*Packet)
		o.Conn.SendPrio(pckt)
		if stat := o.Conn.storage.DeleteIn(cid, subscribe); stat == false {
			logger.Debug("? [NOTICE] deleteinbound returned false (online/subscribe).")
		}
	}
	pb := NewMsgBox(subscribe.Meta.Qos, subscribe.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
	// subscribe box clone
	pbc := pb.Clone(protobase.MDInbound)
	o.client.Subscribe(pbc)
	o.server.NotifySubscribe(o.Conn, pb)
}

// onPING is the heartbeat handler ( other packets reset its timer as well ).
func (o *Online) onPING(packet *Packet) {
	logger.Debug("+ [Heartbeat] Received.")
}

// onSUBACK is a handler which removes the outbound subscribe
// message when QoS >0.
func (o *Online) onSUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onSUBACK", "+ [SubAck] received.")
}

// onPUBACK is a handler which removes the outbound publish
// message when QoS >0.
func (o *Online) onPUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onPUBACK", "+ [PubAck] received.")
	var (
		clid string  = o.Conn.GetClient().GetIdentifier()
		pa   *Puback = NewPuback()
		uid  uuid.UUID
	)
	if err := pa.DecodeFrom(packet.Data); err != nil {
		logger.FDebug("onPUBACK", "- [PubAck] uanble to decode .", "error", err)
		return
	}
	oidstore := o.Conn.storage.GetIDStoreO(clid)
	msgid := pa.Meta.MessageId
	uid, ok := oidstore.GetUUID(msgid)
	if !ok {
		logger.FWarn("onPUBACK", "- [PubAck] no packet coult be found in storage with msgid.", "msgid", msgid)
		return
	}
	np, ok := o.Conn.storage.GetOutbound(clid, uid)
	if !ok {
		logger.FWarn("onPUBACK", "- [PubAck] no packet found with uid.", "uid", uid)
	}
	if !o.Conn.storage.DeleteOut(clid, np) {
		logger.FWarn("onPUBACK", "- [PubAck] failed to remove packet from storage.")
	}
	logger.FDebug("onPUBACK", "+ [PubAck] successfull acknowledge.")
	oidstore.FreeId(msgid)
}

func (o *Online) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", "* [Disconnect] disconnect packet received.")
	atomic.StoreUint32(&(o.Conn).Status, STATGODOWN)
	// NOTE
	// . no need to call (*ClientInterface).Disconnected explicitely,
	//   it will be handled by the main loop. Setting state code is
	//   sufficient.
}

func (o *Online) onPONG(packet *Packet) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}

func (o *Online) onQUEUE(packet *Packet) {
	logger.FDebug("onQUEUE", "* [QUEUE] packet received.")
	o.server.NotifyQueue(o.Conn, nil)
}
