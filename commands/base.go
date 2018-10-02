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

// Commands is the package that provides text-based protocol and
// other facilities for building, parsing and packing commands.
package commands

import "errors"

// ParseNodeInterface is a interface which contains required
// methods that must be implemented fora any implementor
// willing to represent itself as a ParseNode compatible
// container.
type ParseNodeInterface interface {
	Value() interface{}
	Next() ParseNodeInterface
	AddValue(ParseElementInterface)
	SetNext(ParseNodeInterface)
}

// ParseElementInterface is an interface that must be implemented
// by any struct willing to conform parse element.
type ParseElementInterface interface {
	// Gettype returns the element type when given as `ParseNodeType`
	// enum flag.
	GetType() ParseNodeType
	// GetElement returns the current internal value of element pointer.
	GetElement() ParseElementInterface
	// SetType sets underlaying type flag to the given argument
	// and returns true iff underlaying element has no type flag.
	SetType(ParseNodeType) bool
	// SetElement sets the internal element pointer to the given argument.
	SetElement(ParseElementInterface)
}

// Error messages
var (
	EINVARRH error = errors.New("commands: invalid array header.")
	EINVAL   error = errors.New("commands: invalid value.")
	EINLEN   error = errors.New("commands: invalid length.")
	EVIOLEN  error = errors.New("commands: violation, inconsistent length and payload.")
	EVIOTOK  error = errors.New("commands: violation, invalid token.")
	EPARSEF  error = errors.New("commands: unable to parse the given input.")
	EVIOGEN  error = errors.New("commands: protocol violation.")
)

// String constants
const (
	// cErrf is a string format for error representation.
	cErrf string = "commands: %s"
)

// Token represents a single identifier in text stream
type Token byte

// ParseNodeType is the identifier used to identify
// a node and retrive its related information.
type ParseNodeType byte

// Constants for literals represented in parse nodes
const (
	PNDNone ParseNodeType = iota
	PNDArrayLit
	PNDStringLit
	PNDNumeralLit
)

// Format constants
const (
	// Array header string format
	ARRHFormat string = "*%d:%d\r\n"
	// Array header string format with append
	ARRHFormatApnd string = "*%d:%d\r\n%s"
	// String value format
	STRFormat string = "$%d\r\n%s\r\n"
)

// Token constants
const (
	TOK_ART = '*'
	TOK_SEP = ':'
	TOK_CR  = '\r'
	TOK_LF  = '\n'
	TOK_LEN = '$'
)

// Operation constants
const (
	// No operation
	OP_NONE = iota
	// Length start
	OP_LEN
	// Array  start
	OP_ARR_S
	// Array Length
	OP_ARR_L
	// Array Total Size
	OP_ARR_TS
	// Array Total Elements
	OP_ARR_ES
	// Payload start
	OP_PAYLOAD
	// Payload end
	OP_PAYLOAD_END
	// Character
	OP_CH
	// Carriage Return
	OP_CR
	// Line Feed
	OP_LF
	// Carriage Return/Line Feed (\r\n)
	OP_CLRF
)

// arrHeaderInfo is a struct that contains array header
// information such as number of elements, total length
// and their representation as bytes.
type arrHeaderInfo struct {
	c   [][]byte
	ts  int
	es  int
	end int
}

// parseNodeInfo is a struct that contains node informations
// such as capcity, length and element type.
type parseNodeInfo struct {
	cap    int
	length int
	etype  ParseNodeType
}

// arrinfo is a local struct that contains all necessary information
// for an array represented in the textual protocol.
type arrinfo struct {
	curr     int
	expected int
	values   []interface{}
}

// element is a struct that contains all necessary information
// to represent a element in the abstract syntax tree.
type element struct {
	next  *element
	etype ParseNodeType
}
