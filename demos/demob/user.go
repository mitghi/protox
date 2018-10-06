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
	"sync"
	"time"

	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
)

type User struct {
	*client.CLBUser
}

type CustomClient struct {
	sync.RWMutex
	client.Client

	publishing bool
	user       *User
}

func NewUser(opts client.CLBOptions) *User {
	ncl, ok := client.NewCLBUser(opts)
	if !ok {
		return nil
	}
	ret := &User{ncl}
	return ret
}

func NewCustomClient(uid, pid, cid string) *CustomClient {
	return &CustomClient{
		sync.RWMutex{},
		client.Client{Username: uid, Password: pid, ClientId: cid},
		false,
		nil,
	}
}

func (self *CustomClient) Connected(opts protobase.OptionInterface) bool {
	logger.Infof("+ [USER] %s connected with opts %+v.\n", self.Username, opts.(*protocol.ConnackOpts))
	self.user.SetConnected(true)
	self.sendPubs()

	return true
}

func (self *CustomClient) Disconnected(opts protobase.OptCode) {
	logger.Infof("+ [USER] %s disconnected.\n", self.Username)
	self.user.SetConnected(true)
	self.sendPubs()
}

func (self *CustomClient) Subscribe(msg protobase.MsgInterface) {
	logger.Infof("+ [USER] %s subscribed to %s.\n", self.Username, msg.Envelope().Route())
}

func (self *CustomClient) Publish(msg protobase.MsgInterface) {
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
		logger.Infof("+ [USER][publish] %s is sending to topic [%s], message [%s].\n", self.Username, topic, string(message))
	case protobase.MDOutbound:
		logger.Infof("+ [USER][publish] %s has receive dtopic [%s], message [%s].\n", self.Username, topic, string(message))
	}
}

func (self *CustomClient) sendPubs() {
	var isp bool
	self.RLock()
	isp = self.publishing
	self.RUnlock()

	if isp {
		return
	}

	self.Lock()
	if !self.publishing {
		self.publishing = true
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			for _ = range ticker.C {
				self.user.Conn.Publish(*pr, []byte(*pm), byte(*qos), pcallback)
				if !self.user.IsRunning() {
					break
				}
			}
			ticker.Stop()
			return
		}()
	}
	self.Unlock()

}
