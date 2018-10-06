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

	"github.com/google/uuid"
)

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
