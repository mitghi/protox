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

// package networking provides protox connection facilities.
package networking

import (
	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol/packet"
)

// Log is central logger
var logger protobase.LoggingInterface

// Init is package level initializer.
func init() {
	logger = logging.NewLogger("Networking")
}

// Authorization status codes
const (
	RESULTFAIL = 0x01
	RESULTERR  = 0x02
	RESULTNOP  = 0x03
	RESULTOK   = 0x04
)

// Connection status codes
const (
	STATDISCONNECT   uint32 = 1
	STATCONNECTING   uint32 = 2
	STATONLINE       uint32 = 3
	STATDISCONNECTED uint32 = 4
	STATERR          uint32 = 5
	STATFATAL        uint32 = 6
	STATGODOWN       uint32 = 7
)

// Connection response codes
const (
	TMP_RESPOK = 0x10 // TODO: change this later
	RESNON     = 0x00 // no response code is set yet
	RESPFAIL   = 0x01
	RESPOK     = 0x02
	RESPNOK    = 0x03
	RESPERR    = 0x04
)

// Packet is alias type for 'packet.Packet'
type Packet = packet.Packet
