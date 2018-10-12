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

package packet

import (
	"github.com/mitghi/protox/protobase"
)

// Ensure protocol (interface) conformance.
var _ protobase.PacketInterface = (*Packet)(nil)

// NewPacket crafts a new `Packet`. It is
// the default constructor, populates the
// structure with the given arguments and
// returns a pointer to it.
func NewPacket(data []byte, code byte, length int) (p *Packet) {
	p = &Packet{
		Data:   data,
		Code:   code,
		Length: length,
	}
	return p
}

// SetData sets internal byte slice pointer to `data` argument.
func (p *Packet) SetData(data []byte) {
	p.Data = data
}

// SetCode sets internal control packet code to `code` argument.
func (p *Packet) SetCode(code byte) {
	p.Code = code
}

// SetLength sets the total length of byte slice `data`.
func (p *Packet) SetLength(length int) {
	p.Length = length
}

// GetData returns a pointer to packet data.
func (p *Packet) GetData() []byte {
	return p.Data
}

// GetCode returns the associated Protocol Command Code.
func (p *Packet) GetCode() byte {
	return p.Code
}

// GetLength returns bytes total size.
func (p *Packet) GetLength() int {
	return p.Length
}

// IsValid returns wether a given control packet code is in the mapping or not.
func (p *Packet) IsValid() bool {
	// TODO
	// . 0x0 cmd is handled in protocodes
	if p.Code == 0 {
		return false
	}
	return IsValidCommand(p.Code)
}

// IsValidCommand returns wether a packet code is in the mapping or is invalid.
func IsValidCommand(cmd byte) bool {
	_, ok := PROTOCODES[cmd]
	return ok
}
