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

// package packet contains PDUs
package packet

// Control packet map
var (
	PROTOCODES map[byte]string = map[byte]string{
		0x01: "PCONNECT",
		0x02: "PCONNACK",
		0x04: "PQUEUE",
		0x05: "PQUEUEACK",
		0x06: "PPUBACK",
		0x07: "PSUBSCRIBE",
		0x08: "PSUBACK",
		0x09: "PUNSUBSCRIBE",
		0x0A: "PUNSUBACK",
		0x0B: "PPUBLISH",
		0x0C: "PPING",
		0x0D: "PPONG",
		0x0E: "PDISCONNECT",
		// TODO
		// 0x03: "PRESPONSE",
		// 0x04: "PREQACK",
		// 0x05: "PREQUEST",
		// 0x06 : "PRESACK",
	}
)

// Packet represents a PDU ( protocol data unit ).
type Packet struct {
	Data   []byte
	Code   byte
	Length int
}
