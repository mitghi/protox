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
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/mitghi/protox/protocol/packet"
)

// Protox Error messages
var (
	BadMsgTypeError        = errors.New("protox: message type is invalid")
	BadQosError            = errors.New("protox: QoS is invalid")
	BadWillQosError        = errors.New("protox: will QoS is invalid")
	BadLengthEncodingError = errors.New("protox: remaining length field exceeded maximum of 4 bytes")
	BadReturnCodeError     = errors.New("protox: is invalid")
	DataExceedsPacketError = errors.New("protox: data exceeds packet length")
	BsgTooLongError        = errors.New("protox: message is too long")
)

// PanicErr is a wrapper for `error`.
type panicErr struct {
	err error
}

// Error returns the description.
func (p panicErr) Error() string {
	return p.err.Error()
}

// RaiseError raises a panic with a given `error`.
func RaiseError(err error) {
	panic(panicErr{err})
}

// RecoverError catches `panic` inflight and handles them accordingly.
// Using this function must be as follow:
//  func panicableFunc(buff []byte) (err error) {
//    defer func() {
//      err = RecoverError(err, recover())
//    }()
//    ...
//    return err
//  }
func RecoverError(existingErr error, recovered interface{}) error {
	if recovered != nil {
		if pErr, ok := recovered.(panicErr); ok {
			return pErr.err
		} else {
			panic(recovered)
		}
	}
	return existingErr
}

// GetUint8 reads a single `byte` from a `io.Reader` and moves the pointer
// forward.
func GetUint8(r io.Reader, packetRemaining *int32) uint8 {
	if *packetRemaining < 1 {
		RaiseError(DataExceedsPacketError)
	}
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		RaiseError(err)
	}
	*packetRemaining--

	return b[0]
}

// GetUint16 reads a `uint16` from a `io.Reader` and moves the pointer
// forward.
func GetUint16(r io.Reader, packetRemaining *int32) uint16 {
	if *packetRemaining < 2 {
		RaiseError(DataExceedsPacketError)
	}
	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		RaiseError(err)
	}
	*packetRemaining -= 2

	return uint16(b[0])<<8 | uint16(b[1])
}

// GetString reads a `string` from a `io.Reader` and moves the pointer
// forward.
func GetString(r io.Reader, packetRemaining *int32) string {
	strLen := int(GetUint16(r, packetRemaining))
	if int(*packetRemaining) < strLen {
		RaiseError(DataExceedsPacketError)
	}
	b := make([]byte, strLen)
	if _, err := io.ReadFull(r, b); err != nil {
		RaiseError(err)
	}
	*packetRemaining -= int32(strLen)

	return string(b)
}

// GetBytes reads n bytes from the buffer where n
// is total length as uint16.
func GetBytes(r io.Reader, packetRemaining *int32) []byte {
	buffLen := int(GetUint16(r, packetRemaining))
	if int(*packetRemaining) < buffLen {
		RaiseError(DataExceedsPacketError)
	}
	b := make([]byte, buffLen)
	if _, err := io.ReadFull(r, b); err != nil {
		RaiseError(err)
	}
	*packetRemaining -= int32(buffLen)

	return b
}

// SetUint8 writes a single byte to `buf` argument.
func SetUint8(val uint8, buf *bytes.Buffer) {
	buf.WriteByte(byte(val))
}

// SetUint16 writes a `uint16` to `buf` argument as LSB, MSB.
func SetUint16(val uint16, buf *bytes.Buffer) {
	buf.WriteByte(byte(val & 0xff00 >> 8))
	buf.WriteByte(byte(val & 0x00ff))
}

// SetString writes a string as bytes into the `buf` argument.
func SetString(val string, buf *bytes.Buffer) {
	length := uint16(len(val))
	SetUint16(length, buf)
	buf.WriteString(val)
}

// SetBytes writes a byte slice into the `buf` argument.
func SetBytes(val []byte, buf *bytes.Buffer) {
	length := uint16(len(val))
	SetUint16(length, buf)
	buf.Write(val)
}

// BoolToBytes converts a boolean value to a byte.
func BoolToByte(val bool) byte {
	if val {
		return byte(1)
	}
	return byte(0)
}

// DecodeLength decodes total packet length present in the heder.
func DecodeLength(r io.Reader) int32 {
	var v int32
	var buf [1]byte
	var shift uint
	for i := 0; i < 4; i++ {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			RaiseError(err)
		}
		b := buf[0]
		v |= int32(b&0x7f) << shift
		if b&0x80 == 0 {
			return v
		}
		shift += 7
	}
	RaiseError(BadLengthEncodingError)
	panic("unreachable")
}

// EncodeLength encodes `buf` length as LSB, MSB and writes it to the
// buffer given by `buf` argument.
func EncodeLength(length int32, buf *bytes.Buffer) {
	if length == 0 {
		buf.WriteByte(0)
		return
	}
	for length > 0 {
		digit := length & 0x7f
		length = length >> 7
		if length > 0 {
			digit = digit | 0x80
		}
		buf.WriteByte(byte(digit))
	}
}

// GetHeaderBoundary reads the header part of a packet and writes it to a
// byte slice pointer given by `buf` argument, then returns header length.
func GetHeaderBoundary(buf []byte) int {
	lenLen := 1
	for buf[lenLen]&0x80 != 0 {
		lenLen += 1
	}
	lenLen += 1

	return lenLen
}

// ReadPacket reads packets from a `*bufio.Reader` stream. It implements MQTT receive
// algorithm.
//  TODO
//   Write test cases
//   Check reliability against slow send attacks
//   Limit maximum read time
//   Stress test
func ReadPacket(reader *bufio.Reader, pack *[]byte, rl *uint32) error {
	// TODO: write test cases
	var mp uint32 = 1
	for {
		msg, err := reader.ReadByte()
		if err != nil {
			return BadLengthEncodingError
		}
		*pack = append(*pack, msg)
		*rl += uint32(msg&0x7F) * mp
		if msg&0x80 == 0 {
			break
		}
		mp *= 128
		// prevent this trick which is for keeping the connection open
		if mp > 0x200000 {
			return BadLengthEncodingError
		}
	}
	if *rl > 0 {
		var remaining []byte = make([]byte, *rl)
		if _, err := io.ReadFull(reader, remaining); err != nil { // Read 'length' remaining bytes
			return BadLengthEncodingError
		}
		*pack = append(*pack, remaining...) // Add remaining bytes to the network packet
	}

	return nil // OK
}

// CreateHOptions creates a header containing options for Quality of Service, Duplicate
// packets and Retain.
func CreateHOpts(qos byte, dup, ret bool) byte {
	var opts byte = 0x00
	opts |= (qos & 0x03)
	opts = (opts << 1)

	if dup == true {
		opts |= 0x08
	}
	if ret == true {
		opts |= 0x1
	}

	return opts
}

// calcMinQoS returns a feasible QoS value.
func calcMinQoS(prev byte, nev byte) byte {
	if prev > packet.MAXQoS || nev > packet.MAXQoS {
		return packet.MAXQoS
	}
	if nev <= prev {
		return nev
	}
	return prev
}
