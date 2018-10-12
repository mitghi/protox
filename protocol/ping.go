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

func (self *Ping) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if self.Encoded != nil {
		return err
	}
	var (
		vh bytes.Buffer
	)
	self.Header.WriteByte(self.Command)
	EncodeLength(int32(vh.Len()), self.Header)
	// self.Header.Write(vh.Bytes())
	self.Encoded = self.Header

	return err
}

func (self *Ping) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

// TODO: complete this function, this is a stub implementation.
func (self *Ping) UUID() (uid [16]byte) {
	uid = (self.Protocol.Id)
	return uid
}

//
func (self *Pong) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if self.Encoded != nil {
		return err
	}
	var (
		vh bytes.Buffer
	)
	self.Header.WriteByte(self.Command)
	EncodeLength(int32(vh.Len()), self.Header)
	// self.Header.Write(vh.Bytes())
	self.Encoded = self.Header
	return err
}

//
func (self *Pong) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

// TODO: complete this function, this is a stub implementation.
func (self *Pong) UUID() (uid [16]byte) {
	uid = (self.Protocol.Id)
	return uid
}
