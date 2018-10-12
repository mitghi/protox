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
)

//
func (p *Publish) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if p.Encoded != nil {
		return err
	}
	var (
		varHeader bytes.Buffer
		payload   bytes.Buffer
		opts      byte = CreateHOpts(p.Meta.Qos, p.Meta.Dup, p.Meta.Ret)
		// combine 0xF0 and 0x0F masks
		cmd byte = p.Command | opts
	)
	p.Header.WriteByte(cmd)
	if p.Meta.Qos > 0 {
		SetUint16(p.Meta.MessageId, &varHeader)
	}
	SetString(p.Topic, &varHeader)
	SetBytes(p.Message, &payload)
	varHeader.ReadFrom(&payload)
	EncodeLength(int32(varHeader.Len()), p.Header)
	p.Header.Write(varHeader.Bytes())
	p.Encoded = p.Header

	return err
}

//
func (p *Publish) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if len(buff) == 0 {
		return InvalidHeader
	}
	var (
		hbnd            int    = GetHeaderBoundary(buff)
		header          []byte = buff[:hbnd]
		buffrd          *bytes.Reader
		dup             bool
		ret             bool
		qos             byte
		packets         []byte
		packetRemaining int32
	)
	dup, ret, qos = ParseHOptions(header[0] & 0x0F)
	p.Meta.Dup, p.Meta.Ret, p.Meta.Qos = dup, ret, qos
	packets = buff[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	if p.Meta.Qos > 0 {
		p.Meta.MessageId = GetUint16(buffrd, &packetRemaining)
	}
	topic := GetString(buffrd, &packetRemaining)
	p.Topic = topic
	message := GetBytes(buffrd, &packetRemaining)
	p.Message = message

	return err
}
