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
	"fmt"

	"github.com/google/uuid"

	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol/packet"
)

//
type Publish struct {
	Protocol

	Topic   string
	Message []byte
}

//
func NewPublish() *Publish {
	return &Publish{
		Protocol: NewProtocol(CPUBLISH),
		Topic:    "",
	}
}

//
func (self *Publish) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if self.Encoded != nil {
		return err
	}

	var (
		varHeader bytes.Buffer
		payload   bytes.Buffer
		opts      byte = CreateHOpts(self.Meta.Qos, self.Meta.Dup, self.Meta.Ret)
		// combine 0xF0 and 0x0F masks
		cmd byte = self.Command | opts
	)
	self.Header.WriteByte(cmd)
	if self.Meta.Qos > 0 {
		SetUint16(self.Meta.MessageId, &varHeader)
	}

	SetString(self.Topic, &varHeader)
	SetBytes(self.Message, &payload)
	varHeader.ReadFrom(&payload)

	// payload.Write(self.Message)
	// SetString(self.Message, &varHeader)
	EncodeLength(int32(varHeader.Len()), self.Header)
	self.Header.Write(varHeader.Bytes())
	self.Encoded = self.Header

	return err
}

//
func (self *Publish) Decode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

//
func (self *Publish) DecodeFrom(buff *[]byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if len(*buff) == 0 {
		return InvalidHeader
	}
	var (
		hbnd            int    = GetHeaderBoundary(buff)
		header          []byte = (*buff)[:hbnd]
		dup             bool
		ret             bool
		qos             byte
		packets         []byte //packet section
		packetRemaining int32
		buffrd          *bytes.Reader //end
	)
	dup, ret, qos = ParseHOptions(header[0] & 0x0F)
	self.Meta.Dup, self.Meta.Ret, self.Meta.Qos = dup, ret, qos
	packets = (*buff)[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	if self.Meta.Qos > 0 {
		self.Meta.MessageId = GetUint16(buffrd, &packetRemaining)
	}
	topic := GetString(buffrd, &packetRemaining)
	self.Topic = topic
	message := GetBytes(buffrd, &packetRemaining)
	self.Message = message

	return err
}

// TODO: complete this function, this is a stub implementation.
func (self *Publish) Metadata() *ProtoMeta {
	return nil
}

// TODO: complete this function, this is a stub implementation.
func (self *Publish) String() string {
	return fmt.Sprintf("%+v", *self)
}

// TODO: complete this function, this is a stub implementation.
func (self *Publish) UUID() (uid uuid.UUID) {
	uid = (*self.Protocol.Id)
	return uid
}

// GetPacket creates a pointer to a new `Packet` created by using
// internal `Encoded` data.
func (self *Publish) GetPacket() protobase.PacketInterface {
	var (
		data []byte         = self.Encoded.Bytes()
		dlen int            = len(data)
		code byte           = self.Command
		pckt *packet.Packet = packet.NewPacket(&data, code, dlen)
	)

	return pckt
}
