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
type UnSubscribe struct {
	Protocol

	Topic string
}

//
func NewUnSubscribe() *UnSubscribe {
	return &UnSubscribe{
		Protocol: NewProtocol(CUNSUBSCRIBE),
		Topic:    "",
	}
}

//
func (self *UnSubscribe) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if self.Encoded != nil {
		return err
	}

	var (
		varHeader bytes.Buffer
		opts      byte = CreateHOpts(self.Meta.Qos, self.Meta.Dup, self.Meta.Ret)
		cmd       byte = self.Command | opts // combine 0xF0 and 0x0F masks
	)
	self.Header.WriteByte(cmd)
	if self.Meta.Qos > 0 {
		SetUint16(self.Meta.MessageId, &varHeader)
	}
	SetString(self.Topic, &varHeader)
	EncodeLength(int32(varHeader.Len()), self.Header)
	self.Header.Write(varHeader.Bytes())
	self.Encoded = self.Header
	/* d e b u g */
	// var varHeader []byte
	// var payload []byte
	// c.Header.WriteByte(CCONNECT)
	// vhProtocol := "PRX"
	// vhLength := uint16(len(vhProtocol))
	// bstr := []byte(vhProtocol)
	// varHeader = append(varHeader, byte(vhLength & 0xff00 >> 8))
	// varHeader = append(varHeader, byte(vhLength & 0x00ff))
	// varHeader = append(varHeader, bstr...)
	// EncodeLength(int32(len(varHeader)), c.Header)
	// c.Header.Write(varHeader)
	// c.Header.Write(payload)
	/* d e b u g */
	return err
}

//
func (c *UnSubscribe) Decode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

//
func (self *UnSubscribe) DecodeFrom(buff *[]byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if len(*buff) == 0 {
		return InvalidHeader
	}
	var (
		hbnd   int    = GetHeaderBoundary(buff)
		header []byte = (*buff)[:hbnd]
		dup    bool
		ret    bool
		qos    byte

		packets         []byte
		packetRemaining int32
		buffrd          *bytes.Reader
		topic           string
	)
	dup, ret, qos = ParseHOptions(header[0] & 0x0F)
	self.Meta.Dup, self.Meta.Ret, self.Meta.Qos = dup, ret, qos
	packets = (*buff)[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	if self.Meta.Qos > 0 {
		self.Meta.MessageId = GetUint16(buffrd, &packetRemaining)
	}
	topic = GetString(buffrd, &packetRemaining)
	self.Topic = topic

	return err
}

// TODO: complete this function, this is a stub implementation.
func (self *UnSubscribe) Metadata() *ProtoMeta {
	return nil
}

// TODO: complete this function, this is a stub implementation.
func (self *UnSubscribe) String() string {
	return fmt.Sprintf("%+v", *self)
}

// TODO: complete this function, this is a stub implementation.
func (self *UnSubscribe) UUID() (uid uuid.UUID) {
	uid = (*self.Protocol.Id)
	return uid
}

// GetPacket creates a pointer to a new `Packet` created by using
// internal `Encoded` data.
func (self *UnSubscribe) GetPacket() protobase.PacketInterface {
	var (
		data []byte         = self.Encoded.Bytes()
		dlen int            = len(data)
		code byte           = self.Command
		pckt *packet.Packet = packet.NewPacket(&data, code, dlen)
	)

	return pckt
}
