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

// TODO
// . add qos to suback

//
func (sa *Suback) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if sa.Encoded != nil {
		return err
	}

	var (
		varHeader bytes.Buffer
	)
	sa.Header.WriteByte(sa.Command)
	SetUint16(sa.Meta.MessageId, &varHeader)
	EncodeLength(int32(varHeader.Len()), sa.Header)
	sa.Header.Write(varHeader.Bytes())
	sa.Encoded = sa.Header

	return err
}

//
func (sa *Suback) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if len(buff) == 0 {
		return InvalidHeader
	}
	var (
		hbnd int = GetHeaderBoundary(buff)
		// header byte = buff[:hbnd]
		packets         []byte
		packetRemaining int32
		buffrd          *bytes.Reader
		code            uint16
	)
	packets = buff[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	code = GetUint16(buffrd, &packetRemaining)
	sa.Meta.MessageId = code

	return err
}
