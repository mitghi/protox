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

package protobase

// Version
const (
	ProtoVersion = "\x50\x52\x58\x31"
)

// Maximum supported Quality of Service
const (
	MAXQoS byte = 0x1
)

// QoS (Quality of Service) codes
const (
	LQOS0 = 0x00
	LQOS1 = 0x01
	LQOS2 = 0x02
)

// Duplicate option
const (
	NDUP = 0
	YDUP = 1
)

// Retain option ( N = no, Y = yes . ex. NRET = no retain, YRET = retain)
const (
	NRET = 0
	YRET = 1
)

// Authorization status codes
const (
	RESULTFAIL = 0x01
	RESULTERR  = 0x02
	RESULTNOP  = 0x03
	RESULTOK   = 0x04
)

// Connection status codes
const (
	// STATDISCONNECT is the initial state
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

// Connection response header options
const (
	// TODO
	// . add the rest
	RHASSESSION = 0x08
	RCLEANSTART = 0x04
)

// Control packet raw codes
const (
	PNULL        = 0x00
	PCONNECT     = 0x01
	PCONNACK     = 0x02
	PQUEUE       = 0x04
	PQUEUEACK    = 0x05
	PPUBACK      = 0x06
	PSUBSCRIBE   = 0x07
	PSUBACK      = 0x08
	PUNSUBSCRIBE = 0x09
	PUNSUBACK    = 0x0A
	PPUBLISH     = 0x0B
	PPING        = 0x0C
	PPONG        = 0x0D
	PDISCONNECT  = 0x0E
	// TODO
	//  PRESACK      = 0x06
	// NOTE: new control codes should be included
	//
	// RREQUEST
	// RRESPONSE
	// RREQBCST
	// REQRBCST
	// PPROPOS
	// PRPROPS
	// PRESPONSE    = 0x03
	// PREQACK      = 0x04
	// PREQUEST     = 0x05
)

// Server status
const (
	ServerNone    uint32 = 1
	ServerRunning uint32 = 2
	ServerStopped uint32 = 3
	ServerIdle    uint32 = 4
)

// Server command codes
const (
	ForceShutdown uint32 = 5
	Restart       uint32 = 6
)

// Message Direction
const (
	MDInbound MsgDir = iota
	MDOutbound
	MDNone
)

// Client disconnect flags
const (
	PUNone OptCode = iota
	PUSocketError
	PURejected
	PUDisconnect
	PUForceTerminate
	PUAckDeadline
)

// Access Control List mode flags
const (
	ACLModeNormal ACLMode = iota
	ACLModeInclusive
	ACLModeExclusive
)

// Authorization mode flags
const (
	AUTHModeNone AuthMode = iota
	AUTHModeDynamic
	// AUTHModeStrict is defined as follow:
	// - A set of predefined permissions for user groups
	//   which 'Auth' authorizes by validating a route or
	//   action.
	AUTHModeStrict
	AUTHModeRouter
)

// Client mode flags
const (
	CClient ClientMode = iota
	CRouter
	CAgent
)
