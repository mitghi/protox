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
	"github.com/mitghi/protox/protobase"
)

// MsgEnvelope is a struct that conforms to `protobse.MsgEnvelopeInterface`.
// It is used to abstract message content and passed to client functions
// for consumption.
type MsgEnvelope struct {
	route   string
	payload []byte
	// TODO
	// . add meta information
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

// Section: Initializers.

// NewMsgEnvelope is a function that allocates and initializes a new
// `MsgEnvelope` and returns a pointer to it.
func NewMsgEnvelope(route string, payload []byte) *MsgEnvelope {
	return &MsgEnvelope{route, payload}
}

// NewMsgBox is a function that allocates and initializes a new `MsgBox`
// and return a pointer to it.
func NewMsgBox(qos byte, messageId uint16, dir protobase.MsgDir, envelope protobase.MsgEnvelopeInterface) *MsgBox {
	return &MsgBox{qos, messageId, dir, envelope}
}

// Section: MsgEnvelope receiver methods.

// Route returns destination route.
func (me *MsgEnvelope) Route() string {
	return me.route
}

// Payload returns data in bytes.
func (me *MsgEnvelope) Payload() []byte {
	return me.payload
}

// Section: MsgBox receiver methods.

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
