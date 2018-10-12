package networking

import (
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
	_packet "github.com/mitghi/protox/protocol/packet"
)

// Aliases
type (
	Packet = _packet.Packet
	PI     = protobase.PacketInterface
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
	Ping        = protocol.Ping
	Pong        = protocol.Pong
)

// Packet constructors
var (
	NewPacket     func([]byte, byte, int) *Packet = _packet.NewPacket
	NewConnect    func(PI) *Connect               = protocol.NewConnect
	NewDisconnect func(PI) *Disconnect            = protocol.NewDisconnect
	NewConnack    func(PI) *Connack               = protocol.NewConnack
	NewSubscribe  func(PI) *Subscribe             = protocol.NewSubscribe
	NewSuback     func(PI) *Suback                = protocol.NewSuback
	NewPublish    func(PI) *Publish               = protocol.NewPublish
	NewPuback     func(PI) *Puback                = protocol.NewPuback
	NewPing       func(PI) *Ping                  = protocol.NewPing
	NewPong       func(PI) *Pong                  = protocol.NewPong

	NewConnackOpts func() *ConnackOpts = protocol.NewConnackOpts

	NewRawConnect    func() *Connect    = protocol.NewRawConnect
	NewRawDisconnect func() *Disconnect = protocol.NewRawDisconnect
	NewRawConnack    func() *Connack    = protocol.NewRawConnack
	NewRawSubscribe  func() *Subscribe  = protocol.NewRawSubscribe
	NewRawSuback     func() *Suback     = protocol.NewRawSuback
	NewRawPublish    func() *Publish    = protocol.NewRawPublish
	NewRawPuback     func() *Puback     = protocol.NewRawPuback
	NewRawPing       func() *Ping       = protocol.NewRawPing
	NewRawPong       func() *Pong       = protocol.NewRawPong
)

var (
	IsValidCommand func(byte) bool = _packet.IsValidCommand
)
