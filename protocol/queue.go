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
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mitghi/protox/protobase"
)

const (
	QAInitialize protobase.QAction = iota
	QADestroy
	QADrain
	QANone
)

type Queue struct {
	Protocol

	Message    []byte
	Mark       []byte
	Address    string
	ReturnPath string
	Action     protobase.QAction
}

// - MARK: Initializers.

func NewQueue() *Queue {
	return &Queue{
		Protocol:   NewProtocol(CQUEUE),
		Address:    "",
		ReturnPath: "",
		Action:     QANone,
	}
}

// - MARK: Queue section.

func ParseQVarOptions(opts byte) (hasId, hasAddress, hasReturnPath, hasMark bool) {
	hasId = ((opts >> 3) & 0x01) > 0
	hasAddress = ((opts >> 2) & 0x01) > 0
	hasReturnPath = ((opts >> 1) & 0x01) > 0
	hasMark = (opts & 0x01) > 0

	return hasId, hasAddress, hasReturnPath, hasMark
}

func ParseQCMDOptions(opts byte) (initialize, destroy, drain bool) {
	initialize = ((opts >> 3) & 0x01) > 0
	destroy = ((opts >> 2) & 0x01) > 0
	drain = ((opts >> 1) & 0x01) > 0

	return initialize, destroy, drain
}

func ParseQOptions(opts byte) (hasOpts, isDuplicate, hasPayload bool) {
	hasOpts = ((opts >> 3) & 0x01) > 0
	isDuplicate = ((opts >> 2) & 0x01) > 0
	hasPayload = ((opts >> 1) & 0x01) > 0
	return hasOpts, isDuplicate, hasPayload
}

func ParseQAction(opts byte) protobase.QAction {
	switch protobase.QAction(opts) {
	case QAInitialize:
		return QAInitialize
	case QADestroy:
		return QADestroy
	case QADrain:
		return QADrain
	default:
		return QANone
	}
}

func CreateQOpts(hasOpts, isDuplicate, hasPayload bool) (opts byte) {
	if hasOpts {
		opts |= 0x8
	}
	if isDuplicate {
		opts |= 0x4
	}
	if hasPayload {
		opts |= 0x2
	}
	return opts
}

func CreateQVarOpts(hasId, hasAddress, hasReturnPath, hasMark bool) (opts byte) {
	if hasId {
		opts |= 0x8
	}
	if hasAddress {
		opts |= 0x4
	}
	if hasReturnPath {
		opts |= 0x2
	}
	if hasMark {
		opts |= 0x1
	}

	return (opts << 4)
}

func CreateQCMDOpts(action protobase.QAction) (opts byte, err error) {
	switch action {
	case QAInitialize:
		opts |= 0x0
	case QADestroy:
		opts |= 0x1
	case QADrain:
		opts |= 0x2
	default:
		return 0x0, errors.New("protocol(queue): invalid option.")
	}
	return opts, nil
}

//
func (q *Queue) Encode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	if q.Encoded != nil {
		return err
	}

	copts, err := CreateQCMDOpts(q.Action)
	if err != nil {
		return err
	}

	var (
		varHeader     bytes.Buffer
		payload       bytes.Buffer
		hasPayload    bool = len(q.Message) > 0
		hasMessageId  bool = q.Meta.MessageId > 0
		hasAddress    bool = len(q.Address) > 0
		hasReturnPath bool = len(q.ReturnPath) > 0
		hasMark       bool = len(q.Mark) > 0
		hasOpts       bool = hasMessageId || hasAddress || hasReturnPath
		vopts         byte = CreateQVarOpts(hasMessageId, hasAddress, hasReturnPath, hasMark)
		opts          byte = CreateQOpts(hasOpts, q.Meta.Dup, hasPayload)
		// merge proto code and fixed options
		cmd byte = q.Command | opts
	)
	// merge variable and command options
	vopts |= copts
	err = q.Header.WriteByte(cmd)
	if err != nil {
		return err
	}

	SetUint16(uint16(vopts), &varHeader)
	if hasMessageId {
		SetUint16(q.Meta.MessageId, &varHeader)
	}
	if hasAddress {
		SetString(q.Address, &varHeader)
	}
	if hasReturnPath {
		SetString(q.ReturnPath, &varHeader)
	}
	if hasMark {
		SetBytes(q.Mark, &varHeader)
	}
	if hasPayload {
		SetBytes(q.Message, &payload)
	}
	varHeader.ReadFrom(&payload)
	EncodeLength(int32(varHeader.Len()), q.Header)
	q.Header.Write(varHeader.Bytes())
	q.Encoded = q.Header

	return err
}

//
func (q *Queue) Decode() (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()

	return err
}

//
func (q *Queue) DecodeFrom(buff *[]byte) (err error) {
	defer func() {
		err = RecoverError(err, recover())
	}()
	if len(*buff) == 0 {
		return InvalidHeader
	}

	var (
		hbnd            int    = GetHeaderBoundary(buff)
		header          []byte = (*buff)[:hbnd]
		isDuplicate     bool   = false
		hasOpts         bool   = false
		hasPayload      bool   = false
		hasMessageId    bool   = false
		hasAddress      bool   = false
		hasReturnPath   bool   = false
		hasMark         bool   = false
		action          protobase.QAction
		packets         []byte
		packetRemaining int32
		buffrd          *bytes.Reader
		opts            byte
	)

	hasOpts, isDuplicate, hasPayload = ParseQOptions(header[0] & 0x0F)
	q.Meta.Dup = isDuplicate
	packets = (*buff)[hbnd:]
	buffrd = bytes.NewReader(packets)
	packetRemaining = int32(len(packets))
	if hasOpts {
		opts = byte(GetUint16(buffrd, &packetRemaining))
		action = ParseQAction(opts & 0x0F)
		fmt.Printf("%b", opts)
		q.Action = action
		hasMessageId, hasAddress, hasReturnPath, hasMark = ParseQVarOptions((opts & 0xF0) >> 4)
		if hasMessageId {
			q.Meta.MessageId = GetUint16(buffrd, &packetRemaining)
		}
		if hasAddress {
			q.Address = GetString(buffrd, &packetRemaining)
		}
		if hasReturnPath {
			q.ReturnPath = GetString(buffrd, &packetRemaining)
		}
		if hasMark {
			q.Mark = GetBytes(buffrd, &packetRemaining)
		}
	}
	if hasPayload {
		q.Message = GetBytes(buffrd, &packetRemaining)
	}

	return err
}

func (q *Queue) Metadata() *ProtoMeta {
	return nil
}

func (q *Queue) String() string {
	return fmt.Sprintf("%+v", *q)
}

func (q *Queue) UUID() (uid uuid.UUID) {
	uid = (*q.Protocol.Id)
	return uid
}

// GetPacket creates a pointer to a new `Packet` created by using
// internal `Encoded` data.
func (q *Queue) GetPacket() protobase.PacketInterface {
	var (
		data []byte  = q.Encoded.Bytes()
		dlen int     = len(data)
		code byte    = q.Command
		pckt *Packet = NewPacket(&data, code, dlen)
	)

	return pckt
}
