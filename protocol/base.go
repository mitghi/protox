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

// Protocol is the main package for handling low level operations. It
// is the underlaying subsystem for each server willing to use Protox.
package protocol

import (
	"bytes"

	"github.com/google/uuid"

	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
)

// Log is central logger
var logger protobase.LoggingInterface

// Init is package level initializer.
func init() {
	logger = logging.NewLogger("ProtoConnection")
}

// Protocol is protocol structure embedded in each packet. It has functionalities
// for parsing and crafting packets.
type Protocol struct {
	protobase.EDProtocol
	// TODO
	//  check alignment
	Command byte
	Header  *bytes.Buffer
	Encoded *bytes.Buffer
	Meta    *ProtoMeta
	Id      *uuid.UUID
}

// ProtoMeta is meta information embedded in each packet. It contains common
// information such as Quality of Service, Duplicate flag, Retain flag and
// an ID for message identification.
type ProtoMeta struct {
	// TODO
	//  add rest of fields
	Qos        byte
	Dup        bool
	Ret        bool
	MessageId  uint16
	CleanStart bool
	HasSession bool
}

type ConStateInterface interface {
	onCONNECT(packet *Packet)
	onCONNACK(packet *Packet)
	onPUBLISH(packet *Packet)
	onPUBACK(packet *Packet)
	onSUBSCRIBE(packet *Packet)
	onSUBACK(packet *Packet)
	onPING(packet *Packet)
	onPONG(packet *Packet)
	onDISCONNECT(packet *Packet)
	onQUEUE(packet *Packet)
	HandleDefault(packet *Packet) (status bool)
	Handle(packet *Packet)
	Run()
	SetNextState()
}

// ConnectionState is a interface for status of a connection. Each state
// must implement all of its functionalities, during different stages in
// the program, data will be passed between states which changes the behavior
// of its underlaying functionalities. For example during `Genesis` stage, any
// control packet besides `Connect` results in immediate disconnection from the
// broker. After `Genesis`, data will be passed to `Online` state which is opposite
// of `Genesis` state ( `Connect` results in immediate termination ).
type ConnectionState interface {
	ConStateInterface
	SetClient(client protobase.ClientInterface)
	SetServer(server protobase.ServerInterface)
	Shutdown()
}

type baseControlInterface interface {
	Shutdown()
}
