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
	"testing"
)

const (
	usr       string = "test"
	psw              = "$2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC"
	clid             = "identifier_string"
	keepalive int    = 1
)

func makeConnectPacket() (c *Connect) {
	c = NewRawConnect()
	c.Username = usr
	c.Password = psw
	c.ClientId = clid
	c.KeepAlive = keepalive
	return c
}

func TestConnect(t *testing.T) {
	conn := makeConnectPacket()
	err := conn.Encode()
	if err != nil {
		t.Fatal("err!=nil", err)
	}
	fmt.Printf("% #x", conn.Encoded)
	fmt.Println("")
	for _, v := range conn.Encoded.Bytes() {
		fmt.Printf("'%#x' ", v)
	}
	b := conn.Encoded.Bytes()
	// header boundary
	hb := GetHeaderBoundary(b)
	fmt.Println("header boundary", hb, b[:hb])
	fmt.Println("----------------")
	conn2 := NewRawConnect()
	err = conn2.DecodeFrom(b)
	if err != nil {
		t.Fatal("err!=nil", err)
	}
	fmt.Println("encoded:", conn2.String())
	conn3 := NewConnect(conn.GetPacket())
	if conn3 == nil {
		t.Fatal("conn3==nil")
	}
	fmt.Println("uuid: ", conn3.Protocol.Id, conn3.Meta.MessageId)
	// NOTE:
	// . conn2 provides .GetPacket() iff encoded
	// var v *packet.Packet = conn.GetPacket().(*packet.Packet)
}

// TODO:
// . move to stash

// func NewConnectFrom(packet protobase.PacketInterface) (c *Connect) {
// 	c = NewConnect()
// 	if err := c.DecodeFrom(packet.GetData()); err != nil {
// 		return nil
// 	}
// 	return c
// }

// func TestNewConnectFrom(t *testing.T) {
// 	var (
// 		c *Connect = makeConnectPacket()
// 		n *Connect
// 	)
// 	if err := c.Encode(); err != nil {
// 		t.Fatal("packet encoding failed. err!=nil", err)
// 	}
// 	n = NewConnectFrom(c.GetPacket())
// 	if n == nil {
// 		t.Fatal("n==nil")
// 	}
// 	fmt.Println(n)
// }
