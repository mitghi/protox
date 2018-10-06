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
	"fmt"
	"github.com/mitghi/protox/protocol/packet"
	"testing"
)

func TestPing(t *testing.T) {
	conn := NewPing()
	err := conn.Encode()
	if err != nil {
		t.Fatal("err!=nil", err)
	}
	fmt.Println("")
	for _, v := range conn.Encoded.Bytes() {
		fmt.Printf("'%#x' ", v)
	}
	b := conn.Encoded.Bytes()
	p := conn.GetPacket().(*packet.Packet)
	// header boundary
	fmt.Println("content", b, p.Data, p.Code, p.Length)
	fmt.Println("----------------")
}
