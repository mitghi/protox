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
	"bytes"
	"errors"

	"github.com/google/uuid"
	"github.com/mitghi/protox/protobase"
)

// Maximum supported Quality of Service
const (
	MAXQoS byte = 0x1
)

// Error messages
var (
	HeartBeatFailure   = errors.New("protox: HeartBeat not received in timewindow")
	MalformedPacket    = errors.New("protox: Packet is malformed")
	InvalidCmdForState = errors.New("protox: Command inconsistent with state")
	CriticalTimeout    = errors.New("protox: Critial timeout section missed")
	InvalidHeader      = errors.New("protox: Invalid header")
)

// Control packet codes ( shifted to left, mask : 0xF0 )
const (
	CCONNECT     byte = byte(0x1 << 4)
	CCONNACK     byte = byte(0x2 << 4)
	CQUEUE       byte = byte(0x4 << 4)
	CQUEUEACK    byte = byte(0x5 << 4)
	CPUBACK      byte = byte(0x6 << 4)
	CSUBSCRIBE   byte = byte(0x7 << 4)
	CSUBACK      byte = byte(0x8 << 4)
	CUNSUBSCRIBE byte = byte(0x9 << 4)
	CUNSUBACK    byte = byte(0xA << 4)
	CPUBLISH     byte = byte(0xB << 4)
	CPING        byte = byte(0xC << 4)
	CPONG        byte = byte(0xD << 4)
	CDISCONNECT  byte = byte(0xE << 4)
	// TODO
	//  CRESACK byte      = byte(0x6 << 4)
	//  CREQACK      byte = byte(0x4 << 4)
	//  CREQUEST     byte = byte(0x5 << 4)
	// 	CRESPONSE    byte = byte(0x3 << 4)
)

// Quality of Service codes
const (
	LQOS0 = 0x00
	LQOS1 = 0x01
	LQOS2 = 0x02
)

// Duplicate option
const (
	NDUP = 0
	YDUP = 1
)

// Retain option ( N = no, Y = yes . ex. NRET = no retain, YRET = retain)
const (
	NRET = 0
	YRET = 1
)

// Authorization status codes
const (
	RESULTFAIL = 0x01
	RESULTERR  = 0x02
	RESULTNOP  = 0x03
	RESULTOK   = 0x04
)

// Connection status codes
const (
	STATDISCONNECT   uint32 = 1
	STATCONNECTING   uint32 = 2
	STATONLINE       uint32 = 3
	STATDISCONNECTED uint32 = 4
	STATERR          uint32 = 5
	STATFATAL        uint32 = 6
	STATGODOWN       uint32 = 7
)

// Connection response codes
const (
	TMP_RESPOK = 0x10 // TODO: change this later
	RESNON     = 0x00 // no response code is set yet
	RESPFAIL   = 0x01
	RESPOK     = 0x02
	RESPNOK    = 0x03
	RESPERR    = 0x04
)

// Connection response header options
const (
	// TODO
	// . add the rest
	RHASSESSION = 0x08
	RCLEANSTART = 0x04
)

// Control packet raw codes
const (
	PNULL        = 0x00
	PCONNECT     = 0x01
	PCONNACK     = 0x02
	PQUEUE       = 0x04
	PQUEUEACK    = 0x05
	PPUBACK      = 0x06
	PSUBSCRIBE   = 0x07
	PSUBACK      = 0x08
	PUNSUBSCRIBE = 0x09
	PUNSUBACK    = 0x0A
	PPUBLISH     = 0x0B
	PPING        = 0x0C
	PPONG        = 0x0D
	PDISCONNECT  = 0x0E
	// TODO
	//  PRESACK      = 0x06
	// NOTE: new control codes should be included
	//
	// RREQUEST
	// RRESPONSE
	// RREQBCST
	// REQRBCST
	// PPROPOS
	// PRPROPS
	// PRESPONSE    = 0x03
	// PREQACK      = 0x04
	// PREQUEST     = 0x05
)

// Protocol control packets mapping
var PROTOCODES map[byte]string = map[byte]string{
	0x01: "PCONNECT",
	0x02: "PCONNACK",
	0x04: "PQUEUE",
	0x05: "PQUEUEACK",
	0x06: "PPUBACK",
	0x07: "PSUBSCRIBE",
	0x08: "PSUBACK",
	0x09: "PUNSUBSCRIBE",
	0x0A: "PUNSUBACK",
	0x0B: "PPUBLISH",
	0x0C: "PPING",
	0x0D: "PPONG",
	0x0E: "PDISCONNECT",
	// TODO
	//  0x06 : "PRESACK",
	// 0x03: "PRESPONSE",
	// 0x04: "PREQACK",
	// 0x05: "PREQUEST",
}

// NewProtoMeta returns a pointer to a new `ProtoMeta` which includes metadata
// related to a control packet such as qulity of service and message id.
func NewProtoMeta() *ProtoMeta {
	var result *ProtoMeta = &ProtoMeta{
		Qos:        0x00,
		Dup:        false,
		Ret:        false,
		CleanStart: false,
		HasSession: false,
		MessageId:  0,
	}
	return result
}

// NewProtocol returns a new `Protocol` struct. It contains
// neccessary information for header, command code, metadata and ... .
func NewProtocol(code byte) Protocol {
	var uid uuid.UUID
	uid, _ = uuid.NewUUID()
	var result Protocol = Protocol{
		Command: code,
		Header:  &bytes.Buffer{},
		Encoded: nil,
		Meta:    NewProtoMeta(),
		Id:      &uid,
	}
	return result
}

// MessageId is a receiver method which returns an `uint16`
// associated with the current Packet along with a `boolean`
// to indicate operation success.
func (p *Protocol) MessageId() (bool, uint16) {
	if p.Meta != nil {
		return true, p.Meta.MessageId
	}
	return false, 0
}

// CommandCode is a receiver method which returns a `byte`
// as protocol command identifier associated with the Packet.
func (p *Protocol) CommandCode() byte {
	return p.Command
}

// ParseHOptions is a function that parses first 0x0F
// bits into Fixed Header options.
func ParseHOptions(opts byte) (dup, retain bool, qos byte) {
	dup = (opts>>3)&0x01 > 0   // (0x1 << 3) // 0x08 bit
	retain = (opts & 0x01) > 0 // (0x1 << 0)
	qos = (opts >> 1) & 0x03   // (opts & 6) >> 1
	return dup, retain, qos
}

// ParseHCOptions is a function that parses 0x0F bits into
// initial Connect options.
func ParseHCOptions(opts byte) (hasSession bool, hasSessionId bool, cleanStart bool) {
	hasSession = (opts>>3)&0x01 > 0   // 0x8
	hasSessionId = (opts>>2)&0x01 > 0 // 0x4
	cleanStart = (opts>>1)&0x01 > 0
	// TODO
	return hasSession, hasSessionId, cleanStart
}

// Packet represents a PDU ( protocol data units )
type Packet struct {
	Data   *[]byte
	Code   byte
	Length int
}

// IsValidCommand returns wether a packet code is in the mapping or is invalid.
func IsValidCommand(cmd byte) bool {
	_, ok := PROTOCODES[cmd]
	return ok
}

// NewPacket crafts a new `Packet` and returns a pointer to it.
func NewPacket(data *[]byte, code byte, length int) *Packet {
	result := &Packet{
		Data:   data,
		Code:   code,
		Length: length,
	}
	return result
}

// SetData sets internal byte slice pointer to `data` argument.
func (self *Packet) SetData(data *[]byte) {
	self.Data = data
}

// SetCode sets internal control packet code to `code` argument.
func (self *Packet) SetCode(code byte) {
	self.Code = code
}

// SetLength sets the total length of byte slice `data`.
func (self *Packet) SetLength(length int) {
	self.Length = length
}

// IsValid returns wether a given control packet code is in the mapping or not.
func (self *Packet) IsValid() bool {
	if self.Code == 0 {
		return false
	}
	return IsValidCommand(self.Code)
}

// GetData returns a pointer to packet data.
func (self *Packet) GetData() *[]byte {
	return self.Data
}

// GetCode returns the associated Protocol Command Code.
func (self *Packet) GetCode() byte {
	return self.Code
}

// GetLength returns bytes total size.
func (self *Packet) GetLength() int {
	return self.Length
}

// MsgEnvelope is a struct that conforms to `protobse.MsgEnvelopeInterface`.
// It is used to abstract message content and passed to client functions
// for consumption.
type MsgEnvelope struct {
	route   string
	payload []byte
	// TODO
	// . add meta information
}

// NewMsgEnvelope is a function that allocates and initializes a new
// `MsgEnvelope` and returns a pointer to it.
func NewMsgEnvelope(route string, payload []byte) *MsgEnvelope {
	return &MsgEnvelope{route, payload}
}

// Route returns destination route.
func (me *MsgEnvelope) Route() string {
	return me.route
}

// Payload returns data in bytes.
func (me *MsgEnvelope) Payload() []byte {
	return me.payload
}

// MsgBox is a struct that conforms to `protobase.MessageBoxInterface`.
// It is used to abstract Protocol packets and simplify packet storage.
type MsgBox struct {
	qos       byte
	messageId uint16
	dir       protobase.MsgDir
	envelope  protobase.MsgEnvelopeInterface
	// TODO
	// meta      protobase.MetaEnvelopeInterface
}

// NewMsgBox is a function that allocates and initializes a new `MsgBox`
// and return a pointer to it.
func NewMsgBox(qos byte, messageId uint16, dir protobase.MsgDir, envelope protobase.MsgEnvelopeInterface) *MsgBox {
	return &MsgBox{qos, messageId, dir, envelope}
}

// SetWishQoS sets calculates and sets a feasible Quality of Service
// value.
func (mb *MsgBox) SetWishQoS(qos byte) {
	mb.qos = calcMinQoS(mb.qos, qos)
}

// MessageId returns the associated ID.
func (mb *MsgBox) MessageId() uint16 {
	return mb.messageId
}

// QoS returns associated Quality of Service value.
func (mb *MsgBox) QoS() byte {
	return mb.qos
}

// Dir returns messeage Direction (i.e. Inbound or Outbound).
func (mb *MsgBox) Dir() protobase.MsgDir {
	return mb.dir
}

// Envelope returns the message content (i.e. Topic, Payload, and ....).
func (mb *MsgBox) Envelope() protobase.MsgEnvelopeInterface {
	return mb.envelope
}

// Clone deep-copies and returns current message and set its
// direction to argument `dir`.
func (mb *MsgBox) Clone(dir protobase.MsgDir) protobase.MsgInterface {
	var (
		e       protobase.MsgEnvelopeInterface = mb.envelope
		route   string
		payload []byte
	)
	switch e.(type) {
	case *MsgEnvelope:
		me := e.(*MsgEnvelope)
		route = me.route
		payload = me.payload
	default:
		route = e.Route()
		payload = e.Payload()
	}
	nme := NewMsgEnvelope(route, payload)
	nmb := NewMsgBox(mb.qos, mb.messageId, dir, nme)
	return nmb
}
