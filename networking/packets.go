package networking

import (
	"github.com/mitghi/protox/protocol"
)

// Packet alias
type (
  Connect     = protocol.Connect
  Disconnect  = protocol.Disconnect
  Connack     = protocol.Connack
  ConnackOpts = protocol.ConnackOpts
  Subscribe   = protocol.Subscribe
  Suback      = protocol.Suback
  Publish     = protocol.Publish
  Puback      = protocol.Puback
)

// Packet constructors
var (
   NewConnect     func() *Connect     = protocol.NewConnect
   NewDisconnect  func() *Disconnect  = protocol.NewDisconnect
   NewConnack     func() *Connack     = protocol.NewConnack
   NewConnackOpts func() *ConnackOpts = protocol.NewConnackOpts
   NewSubscribe   func() *Subscribe   = protocol.NewSubscribe
   NewSuback      func() *Suback      = protocol.NewSuback
   NewPublish     func() *Publish     = protocol.NewPublish
   NewPuback      func() *Puback      = protocol.NewPuback
)
