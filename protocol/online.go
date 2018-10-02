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
func (self *Online) HandleDefault(packet *Packet) (status bool) {
	return true
}

// Shutdown sets the status to error which notifies the supervisor
// and cleanly terminates the connection.
func (self *Online) Shutdown() {
	logger.Debug("* [Genesis] Closing.")
	atomic.StoreUint32(&(self.Conn).Status, STATERR)
	self.client.Disconnected(protobase.PUForceTerminate)
}

// onCONNECT is not valid in this stage.
func (self *Online) onCONNECT(packet *Packet) {
	// TODO
}

// onCONNACK is not valid in this stage.
func (self *Online) onCONNACK(packet *Packet) {
	// TODO
}

// onPUBLISH is the handler for `Publish` packets.
func (self *Online) onPUBLISH(packet *Packet) {
	var (
		publish  *Publish = NewPublish()
		cid      string   = self.client.GetIdentifier()
		userType protobase.AuthUserType
	)
	if err := publish.DecodeFrom(packet.Data); err != nil {
		logger.Debugf("- [DecodeErr(onPublish)] Unable to decode data for Client(%s).", err, cid)
		self.Shutdown()
		return
	}
	userType, err := self.Conn.auth.GetUserType(cid)
	if err != nil {
		logger.Debugf("onPUBLISH", "- [Packet] unable to find associated User Type for Client(%s).", cid)
		return
	}
	if self.Conn.permissionDelegate != nil {
		if !self.Conn.permissionDelegate(self.Conn.auth, "can", "publish", publish.Topic) {
			logger.Debugf("onPUBLISH", "- [Packet] unable to find corresponding permission for Client(%s).", cid)
			self.Shutdown()
			return
		}
	} else {
		role := self.Conn.auth.GetACL().GetRole((string)(userType))
		if role == nil {
			logger.Debug("onPUBLISH", "- [Role] role==nil.")
			self.Shutdown()
			return
		}
    // TDOO
    // . refactor hard-coded permissions
		if !role.HasPerm("can", "publish", publish.Topic) {
			logger.Debugf("onPUBLISH", "- [Packet] unable to find corresponding permission ( direct ) for Client(%s).", cid)
			self.Shutdown()
			return
		}
	}
	if stat := self.Conn.storage.AddInbound(cid, publish); stat == false {
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
			self.Shutdown()
			return
		}
		var pckt *Packet = puback.GetPacket().(*Packet)
		self.Conn.SendPrio(pckt)
		if stat := self.Conn.storage.DeleteIn(cid, publish); stat == false {
			logger.Debug("? [NOTICE] deleteinbound returned false (online/publish).")
		}
	}
	pb := NewMsgBox(publish.Meta.Qos, publish.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(publish.Topic, publish.Message))
	// publish box clone
	pbc := pb.Clone(protobase.MDInbound)
	self.client.Publish(pbc)
	self.server.NotifyPublish(self.Conn, pb)
}

// onSUBSCRIBE is the handler for `Subscribe` packets.
func (self *Online) onSUBSCRIBE(packet *Packet) {
	var (
		subscribe *Subscribe = NewSubscribe()
		cid       string     = self.client.GetIdentifier()
		userType  protobase.AuthUserType
	)
	if err := subscribe.DecodeFrom(packet.Data); err != nil {
		logger.Debugf("- [DecodeErr(onSubscribe)] Unable to decode data for Client(%s).", err, cid)
		self.Shutdown()
		return
	}
	userType, err := self.Conn.auth.GetUserType(cid)
	if err != nil {
		logger.Debugf("onSUBSCRIBE", "- [Packet] unable to find associated User Type for Client(%s).", cid)
		self.Shutdown()
		return
	}
  if self.Conn.auth.GetMode() != protobase.AUTHModeNone {
    if self.Conn.permissionDelegate != nil {
      if !self.Conn.permissionDelegate(self.Conn.auth, "can", "subscribe", subscribe.Topic) {
        self.Shutdown()
        return
      }
    } else {    
      role := self.Conn.auth.GetACL().GetRole((string)(userType))
      if role == nil {
        logger.Debug("onSUBSCRIBE", "- [Role] role==nil.")
        self.Shutdown()
        return
      }
      if !role.HasPerm("can", "subscribe", subscribe.Topic) {
        self.Shutdown()
        return
      }
    }
  }
	if stat := self.Conn.storage.AddInbound(cid, subscribe); stat == false {
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
			self.Shutdown()
			return
		}
		var pckt *Packet = suback.GetPacket().(*Packet)
		self.Conn.SendPrio(pckt)
		if stat := self.Conn.storage.DeleteIn(cid, subscribe); stat == false {
			logger.Debug("? [NOTICE] deleteinbound returned false (online/subscribe).")
		}
	}
	pb := NewMsgBox(subscribe.Meta.Qos, subscribe.Meta.MessageId, protobase.MDInbound, NewMsgEnvelope(subscribe.Topic, nil))
	// subscribe box clone
	pbc := pb.Clone(protobase.MDInbound)
	self.client.Subscribe(pbc)
	self.server.NotifySubscribe(self.Conn, pb)
}

// onPING is the heartbeat handler ( other packets reset its timer as well ).
func (self *Online) onPING(packet *Packet) {
	logger.Debug("+ [Heartbeat] Received.")
}

// onSUBACK is a handler which removes the outbound subscribe
// message when QoS >0.
func (self *Online) onSUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onSUBACK", "+ [SubAck] received.")
}

// onPUBACK is a handler which removes the outbound publish
// message when QoS >0.
func (self *Online) onPUBACK(packet *Packet) {
	// TODO
	logger.FDebug("onPUBACK", "+ [PubAck] received.")
	var (
		clid string  = self.Conn.GetClient().GetIdentifier()
		pa   *Puback = NewPuback()
		uid  uuid.UUID
	)
	if err := pa.DecodeFrom(packet.Data); err != nil {
		logger.FDebug("onPUBACK", "- [PubAck] uanble to decode .", "error", err)
		return
	}
	oidstore := self.Conn.storage.GetIDStoreO(clid)
	msgid := pa.Meta.MessageId
	uid, ok := oidstore.GetUUID(msgid)
	if !ok {
		logger.FWarn("onPUBACK", "- [PubAck] no packet coult be found in storage with msgid.", "msgid", msgid)
		return
	}
	np, ok := self.Conn.storage.GetOutbound(clid, uid)
	if !ok {
		logger.FWarn("onPUBACK", "- [PubAck] no packet found with uid.", "uid", uid)
	}
	if !self.Conn.storage.DeleteOut(clid, np) {
		logger.FWarn("onPUBACK", "- [PubAck] failed to remove packet from storage.")
	}
	logger.FDebug("onPUBACK", "+ [PubAck] successfull acknowledge.")
	oidstore.FreeId(msgid)
}

func (self *Online) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", "* [Disconnect] disconnect packet received.")
	atomic.StoreUint32(&(self.Conn).Status, STATGODOWN)
	// NOTE
	// . no need to call (*ClientInterface).Disconnected explicitely,
	//   it will be handled by the main loop. Setting state code is
	//   sufficient.
}

func (self *Online) onPONG(packet *Packet) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}

func (self *Online) onQUEUE(packet *Packet) {
	logger.FDebug("onQUEUE", "* [QUEUE] packet received.")
	self.server.NotifyQueue(self.Conn, nil)
}
