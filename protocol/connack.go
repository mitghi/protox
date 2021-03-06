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

	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol/packet"
)

// Encode is a routine for encoding `Connack` packet.
func (ca *Connack) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if ca.Encoded != nil {
		return
	}

	var (
		varHeader bytes.Buffer
		flags     uint8
		cmd       = ca.Command
	)

	// TODO
	// . add response options
	if ca.Meta.HasSession {
		flags |= packet.RHASSESSION // 0x8
	}
	if ca.SessionId != "" {
		flags |= 0x4
	}
	if ca.Meta.CleanStart {
		flags |= 0x2
	}
	cmd |= flags

	ca.Header.WriteByte(cmd)
	SetUint8(ca.ResultCode, &varHeader)
	if ca.SessionId != "" {
		SetString(ca.SessionId, &varHeader)
	}
	EncodeLength(int32(varHeader.Len()), ca.Header)
	ca.Header.Write(varHeader.Bytes())
	ca.Encoded = ca.Header

	return err
}

// SetSessionId sets the `SessionId` in the header.
func (ca *Connack) SetSessionId(sessionId string) {
	ca.SessionId = sessionId
}

// SetResultCode sets the `ResultCode` in the header.
func (ca *Connack) SetResultCode(resultCode byte) {
	ca.ResultCode = resultCode
}

// DecodeFrom decodes a packet from `buff` argument. It is not implemented
// because it is always the server responsibilty to send this packet.
func (ca *Connack) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if len(buff) == 0 {
		return InvalidHeader
	}
	var (
		hbnd            int = GetHeaderBoundary(buff)
		header          []byte
		packets         []byte
		packetRemaining int32
		opts            byte
		buffrd          *bytes.Reader
	)
	/* d e b u g */
	// TODO
	header = buff[:hbnd]
	opts = header[0] & 0x0f
	hasSession, hasSessionId, cleanStart := ParseHCOptions(opts)
	ca.Meta.HasSession = hasSession
	ca.Meta.CleanStart = cleanStart
	/* d e b u g */

	packets = buff[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	ca.ResultCode = GetUint8(buffrd, &packetRemaining)
	/* d e b u g */
	// if packetRemaining > 0 {
	// 	logger.FDebug("DecodeFrom", "packetRemaining>0", packetRemaining)
	// 	ca.SessionId = GetString(buffrd, &packetRemaining)
	// }
	/* d e b u g */
	if hasSessionId && packetRemaining > 0 {
		ca.SessionId = GetString(buffrd, &packetRemaining)
	}
	return err
}

type ConnackOpts struct {
	// TODO
	optcode    byte
	ResultCode byte
	SessionId  string
	HasSession bool
	CleanStart bool
}

func NewConnackOpts() *ConnackOpts {
	return &ConnackOpts{optcode: packet.CCONNACK}
}

func (ca *ConnackOpts) StateCode() protobase.OptCode {
	return (protobase.OptCode)(ca.optcode)
}
func (ca *ConnackOpts) Opts() interface{} {
	// TODO
	return nil
}
func (ca *ConnackOpts) Match(protobase.OptCode) bool {
	return false
}

func (ca *ConnackOpts) ParseFrom(cack *Connack) {
	ca.ResultCode = cack.ResultCode
	ca.SessionId = cack.SessionId
	ca.HasSession = cack.Meta.HasSession
	ca.CleanStart = cack.Meta.CleanStart
}
