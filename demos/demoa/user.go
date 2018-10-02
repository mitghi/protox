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

package main

import (
	"fmt"

	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
)

type User struct {
	client.Client
}

func (self *User) Connected(opts protobase.OptionInterface) bool {
	fmt.Printf("[USER] %s connected.\n", self.Username)
	return true
}

func (self *User) Disconnected(opts protobase.OptCode) {
	fmt.Printf("[USER] %s disconnected.\n", self.Username)
}

func (self *User) Subscribe(msg protobase.MsgInterface) {
	fmt.Printf("[USER] %s subscribed to %s.\n", self.Username, msg.Envelope().Route())
}

func (self *User) Publish(msg protobase.MsgInterface) {
	var (
		dir      protobase.MsgDir               = msg.Dir()
		envelope protobase.MsgEnvelopeInterface = msg.Envelope()
	)

	var (
		topic   string = envelope.Route()
		message []byte = envelope.Payload()
	)
	switch dir {
	case protobase.MDInbound:
		fmt.Printf("[USER][publish] %s is sending to topic [%s], message [%s].\n", self.Username, topic, string(message))
	case protobase.MDOutbound:
		fmt.Printf("[USER][publish] %s has received tpic [%s], message [%s].\n", self.Username, topic, string(message))
	}
}
