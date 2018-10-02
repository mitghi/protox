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

	"github.com/mitghi/protox/protobase"
)

// Disconnect is a control packet. It temrinates the connection.
type Disconnect struct {
	Protocol
}

// NewDisconnect returns a new `Disconnect` control packet.
func NewDisconnect() *Disconnect {
	result := &Disconnect{
		Protocol: NewProtocol(CDISCONNECT),
	}

	return result
}

// Encode is a routine for encoding `Disconnect` packet.
func (self *Disconnect) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if self.Encoded != nil {
		return
	}
	var (
		varHeader bytes.Buffer
	)
	self.Header.WriteByte(self.Command)
	EncodeLength(int32(varHeader.Len()), self.Header)
	self.Encoded = self.Header

	return err
}

// DecodeFrom decodes a packet from `buff` argument. It is not implemented
// because it is always the server responsibilty to send this packet.
func (self *Disconnect) DecodeFrom(buff *[]byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

// Decode decodes the internal data. It is not implemented because
// it is always server responsibility to send th is packet.
func (self *Disconnect) Decode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

// TODO: complete this function, this is a stub implementation.
func (self *Disconnect) Metadata() *ProtoMeta {
	return nil
}

// TODO: complete this function, this is a stub implementation.
func (self *Disconnect) String() string {
	return ""
}

// TODO: complete this function, this is a stub implementation.
func (self *Disconnect) UUID() (uid uuid.UUID) {
	uid = (*self.Protocol.Id)
	return uid
}

// GetPacket creates a pointer to a new `Packet` created by using
// internal `Encoded` data.
func (self *Disconnect) GetPacket() protobase.PacketInterface {
	var (
		data []byte  = self.Encoded.Bytes()
		dlen int     = len(data)
		code byte    = self.Command
		pckt *Packet = NewPacket(&data, code, dlen)
	)

	return pckt
}

// TODO:
//
// func (self *Disconnect) SetCode(code byte) {
//   self.SetCode(code)
// }
