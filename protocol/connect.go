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

// Connect is the initial connection packet. This packet is always
// sent by a client to the broker. Note - this might change in
// future versions.
type Connect struct {
	Protocol

	ClientId    string
	KeepAlive   int
	WillTopic   string
	WillMessage string
	WillQos     int
	WillRetain  int
	Username    string
	Password    string
	CleanStart  bool
	Version     string
}

// NewConnect returns a new `Connect` control packet. It is used during initial `Genesis`
// stage to decode packet data sent by a client.
func NewConnect() *Connect {
	return &Connect{
		Protocol:    NewProtocol(CCONNECT),
		ClientId:    "",
		KeepAlive:   0,
		WillTopic:   "",
		WillMessage: "",
		WillQos:     0,
		WillRetain:  0,
		Username:    "",
		Password:    "",
		CleanStart:  false,
		Version:     "",
		// TODO
		// . add control byte options ( 0x0f )
	}
}

// Encode is for encoding `Connect` data. It is not implemented because
// this packet is always sent from a client to the broker. Note - this
// might change in future versions.
func (cn *Connect) Encode() (err error) {
  // TODO:
  // . process credentials in another
  //   method.
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
	/* d e b u g */  
	// cmd = cn.Command | flags
	// NOTE: TODO:
	// . this has changed
	// original
	// cmd = cn.Command
	// new
	/* d e b u g */  
	cmd = cn.Command
	if cn.Meta.CleanStart {
		var opts byte = 0x8 // clean-start bit
		cmd |= opts
	}
	cn.Header.WriteByte(cmd)
	/* d e b u g - s t a r t*/
	hasPassword := (flags & 0x1) != 0
	hasClientId := (flags & 0x2) != 0
	hasKeepalive := (flags & 0x4) != 0
	hasUsername := (flags & 0x8) != 0
	logger.FInfo("Encode", "* [Connection] --OPTIONS[keepalive, clid, clusrname, clpasswd]=(", hasKeepalive, hasClientId, hasUsername, hasPassword, ")--")
	/* d e b u g -  e n d  */
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
	/* d e b u g */  
	// varHeader = append(varHeader, byte(vhLength & 0xff00 >> 8))
	// varHeader = append(varHeader, byte(vhLength & 0x00ff))
	// varHeader = append(varHeader, bstr...)
	// EncodeLength(int32(len(varHeader)), cn.Header)
	// cn.Header.Write(varHeader)
	// cn.Header.Write(payload)
	/* d e b u g */  
	return err
}

// Decode is decoding `Connect` data. It is not implemented yet.
func (cn *Connect) Decode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

// DecodeFrom decodes packet data from function argument `buff`. It
// saves its result to the initialized structure.
func (cn *Connect) DecodeFrom(buff *[]byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	/* d e b u g */
	// lenLen := 1
	// for buff[lenLen] & 0x80 != 0{
	//   lenLen += 1
	// }
	/* d e b u g */  
	headerBoundary := GetHeaderBoundary(buff)
	packets := (*buff)[headerBoundary:]
	/* d e b u g */
	// TODO
	// . parse the rest of options
	header := (*buff)[:headerBoundary]
	opts := header[0] & 0x0f
	if ((opts & 0x8) >> 3) == 1 {
		cn.Meta.CleanStart = true
	}
	/* d e b u g */
	buffreader := bytes.NewReader(packets)
	packetRemaining := int32(len(packets))
	versionStr := GetString(buffreader, &packetRemaining)
	cn.Version = versionStr
	flags := GetUint8(buffreader, &packetRemaining)
	hasPassword := (flags & 0x1) != 0
	hasClientId := (flags & 0x2) != 0
	hasKeepalive := (flags & 0x4) != 0
	hasUsername := (flags & 0x8) != 0
	logger.Debug("--OPTIONS[keepalive, clid, clusrname, clpasswd]=(", hasKeepalive, hasClientId, hasUsername, hasPassword, ")--")
	// TODO: add new headers
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


func (cn *Connect) Metadata() *ProtoMeta {
	return nil
}


func (cn *Connect) String() string {
	return fmt.Sprintf("\n\tUsername(%s), Password(NONE),\n\t ClientId(%s), KeepAlive(%d), Version(%s).",
		cn.Username, cn.ClientId, cn.KeepAlive, cn.Version)
}


func (cn *Connect) UUID() (uid uuid.UUID) {
	uid = (*cn.Protocol.Id)
	return uid
}

// GetPacket creates a pointer to a new `Packet` created by using
// internal `Encoded` data.
func (cn *Connect) GetPacket() protobase.PacketInterface {
	var (
		data []byte  = cn.Encoded.Bytes()
		dlen int     = len(data)
		code byte    = cn.Command
		pckt *packet.Packet = packet.NewPacket(&data, code, dlen)
	)

	return pckt
}
