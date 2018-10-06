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

func TestConnack(t *testing.T) {
	c := NewConnack()
	c.Meta.HasSession = true
	c.ResultCode = RESPFAIL
	c.Encode()

	np := c.GetPacket().(*packet.Packet)
	nc := NewConnack()
	if err := nc.DecodeFrom(np.Data); err != nil {
		t.Fatal("err!=nil, expected nil. Unable to decode packet.")
	}
	if nc.ResultCode != c.ResultCode {
		t.Fatal("different resultcodes, expected same.", nc.ResultCode, c.ResultCode)
	}
	if nc.Meta.HasSession == false {
		t.Fatal("nc.Meta.HasSession == false, expected true")
	}
	fmt.Println("")
	for _, v := range c.Encoded.Bytes() {
		fmt.Printf("%x ", v)
	}
}
