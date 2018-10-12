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

	"github.com/mitghi/protox/protobase"
)

// Encode encodes `Connect` data.
func (cn *Connect) Encode() (err error) {
	// TODO:
	// . process credentials in another
	//   method.
	const fn string = "Encode"
	defer func() {
		err = RecoverError(err, recover())
	}()
	var (
		flags uint8 = 0
		vh    bytes.Buffer
		pl    bytes.Buffer
		cmd   byte
	)
	if cn.Password != "" {
		flags |= 0x1
	}
	// client id
	flags |= 0x2
	if cn.KeepAlive > 0 {
		flags |= 0x4
	}
	if cn.Username != "" {
		flags |= 0x8
	}
	cmd = cn.Command
	if cn.Meta.CleanStart {
		cmd |= 0x8 // clean-start bit
	}
	cn.Header.WriteByte(cmd)
	hasPassword := (flags & 0x1) != 0
	hasClientId := (flags & 0x2) != 0
	hasKeepalive := (flags & 0x4) != 0
	hasUsername := (flags & 0x8) != 0
	logger.FInfo(fn, "* [Connection] --OPTIONS[keepalive, clid, clusrname, clpasswd]=(",
		hasKeepalive, hasClientId, hasUsername, hasPassword, ")--")
	vhProtocol := []byte(protobase.ProtoVersion)
	SetString(string(vhProtocol), &vh)
	SetUint8(flags, &vh)
	if hasPassword {
		logger.FTrace(1, "Encode", "setting string")
		SetString(cn.Password, &pl)
	}
	if hasClientId {
		logger.FTrace(1, "Encode", "setting clientid")
		SetString(cn.ClientId, &pl)
	}
	if hasKeepalive {
		logger.FTrace(1, "Encode", "setting keepalive")
		SetUint16(uint16(cn.KeepAlive), &pl)
	}
	if hasUsername {
		logger.FTrace(1, "Encode", "setting username")
		SetString(cn.Username, &pl)
	}
	if _, err := vh.Write(pl.Bytes()); err != nil {
		return err
	}
	EncodeLength(int32(vh.Len()), cn.Header)
	logger.FDebugf("Encode", "% x\n", cn.Header.Bytes())
	cn.Header.Write(vh.Bytes())
	cn.Encoded = cn.Header
	return err
}

// DecodeFrom decodes packet data from function argument `buff`. It
// saves its result to the initialized structure.
func (cn *Connect) DecodeFrom(buff []byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	headerBoundary := GetHeaderBoundary(buff)
	packets := buff[headerBoundary:]
	header := buff[:headerBoundary]
	opts := header[0] & 0x0f
	if ((opts & 0x8) >> 3) == 1 {
		cn.Meta.CleanStart = true
	}
	buffreader := bytes.NewReader(packets)
	packetRemaining := int32(len(packets))
	versionStr := GetString(buffreader, &packetRemaining)
	cn.Version = versionStr
	flags := GetUint8(buffreader, &packetRemaining)
	hasPassword := (flags & 0x1) != 0
	hasClientId := (flags & 0x2) != 0
	hasKeepalive := (flags & 0x4) != 0
	hasUsername := (flags & 0x8) != 0
	logger.Debug("--OPTIONS[keepalive, clid, clusrname, clpasswd]=(",
		hasKeepalive, hasClientId, hasUsername, hasPassword, ")--")
	if hasPassword {
		password := GetString(buffreader, &packetRemaining)
		cn.Password = password
	}
	if hasClientId {
		clientId := GetString(buffreader, &packetRemaining)
		cn.ClientId = clientId
	}
	if hasKeepalive {
		keepAlive := GetUint16(buffreader, &packetRemaining)
		cn.KeepAlive = int(keepAlive)
	}
	if hasUsername {
		username := GetString(buffreader, &packetRemaining)
		cn.Username = username
	}
	return err
}

func (cn *Connect) String() string {
	return fmt.Sprintf("\n\tUsername(%s), Password(NONE),\n\t ClientId(%s), KeepAlive(%d), Version(%s).",
		cn.Username, cn.ClientId, cn.KeepAlive, cn.Version)
}
