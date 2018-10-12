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
	"fmt"
	"testing"

	"github.com/mitghi/protox/protocol/packet"
)

func TestPublishAndDecode(t *testing.T) {
	conn := NewRawPublish()
	conn.Topic = "a/awesome/topic"
	conn.Message = []byte("a/awesome/topic")
	err := conn.Encode()
	if err != nil {
		t.Fatal("err!=nil", err)
	}
	fmt.Println("")
	for _, v := range conn.Encoded.Bytes() {
		fmt.Printf("'% x' ", v)
	}
	b := conn.Encoded.Bytes()
	p := conn.GetPacket().(*packet.Packet)
	// header boundary
	fmt.Println("content", b, p.Data, p.Code, p.Length)

	nc := NewPublish(p)
	// nc.DecodeFrom(p.Data)
	if string(nc.Message) != string(conn.Message) {
		t.Fatal("nc.Message!=conn.Message, expected equal")
	}
	fmt.Println("----------------")
	fmt.Println(string(nc.Message), string(conn.Message), conn.Topic, nc.Topic)

	// length check
	var nb bytes.Buffer
	nb.Write((p.Data)[1:])
	dl := DecodeLength(&nb)
	fmt.Println("dl is :", dl)
	pl := len((p.Data)[2:])
	if dl != int32(pl) {
		t.Fatal("dl!=packLen, expected equal", dl, int32(pl))
	}
}

func TestPublishWithQoS(t *testing.T) {
	conn := NewRawPublish()
	conn.Topic = "a/great/simple/topic"
	conn.Message = []byte("a simple message")
	conn.Meta.Qos = 0x1
	if err := conn.Encode(); err != nil {
		t.Fatal("err!=nil", err)
	}

	for _, v := range conn.Encoded.Bytes() {
		fmt.Printf("%#x", v)
	}
	packet := conn.GetPacket().(*packet.Packet)
	conn2 := NewPublish(packet)
	if conn2 == nil {
		fmt.Println("cannot decode the packet from old package")
	}
	// if err := conn2.DecodeFrom(packet.Data); err != nil {
	// 	fmt.Println("cannot decode the packet from old package")
	// }
	fmt.Println("this is the data:", conn2.Message, conn2.Topic, conn2.Meta.Qos)
}

func TestDecodeRaw(t *testing.T) {
	p := NewRawPublish()
	p.Topic = "test"
	p.Message = []byte("test")
	p.Meta.Qos = 1
	p.Meta.MessageId = 1
	p.Encode()
	fmt.Println("")
	for _, v := range p.Encoded.Bytes() {
		fmt.Printf("%2x ", v) // b2 0e 00 01 00 04 74 65 73 74 00 04 74 65 73 74
	}
	np := p.Encoded.Bytes()
	nnp := NewRawPublish()
	err := nnp.DecodeFrom(np)
	if err != nil {
		t.Fatal("err!=nil, cannot decode", err)
	}
}
