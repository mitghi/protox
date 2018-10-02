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

package client

import (
	"sync"

	"github.com/mitghi/protox/protobase"
)

/*
* TODO:
* . remove `GetTopics` and `AddTopics` from Client
 */

var _ protobase.ClientInterface = (*Client)(nil)

//
func NewClient(username string, password string, cid string) *Client {
	client := &Client{
		RWMutex:     &sync.RWMutex{},
		ClientId:    cid,
		WillMessage: "",
		WillQos:     "",
		WillRetain:  "",
		Username:    username,
		Password:    password,
		// Server: nil,
	}
	// client.Conn.SetClient(client)
	return client
}

//
func (self *Client) SetServer(server ServerInterface) {
	self.Server = server
}

//
func (self *Client) Connected(opts protobase.OptionInterface) bool {
	logger.Debug("+ [Client] Connected.")
	// (*self.Server).NotifyConnected(self)
	return true
}

//
func (self *Client) Disconnected(opts protobase.OptCode) {
	// (*self.Server).NotifyDisconnected(self)
	logger.Debug("- [Client] Disconnected.")
}

func (self *Client) Publish(msg protobase.MsgInterface) {
	switch msg.Dir() {
	case protobase.MDInbound:
		logger.Debug("+ [Client] Marked to notify status.")
	case protobase.MDOutbound:
		logger.Debugf("+ [Client] Publishing a message with QoS(%d).", int(msg.QoS()))
	}
}

func (self *Client) Subscribe(msg protobase.MsgInterface) {
	topic := msg.Envelope().Route()
	logger.Debug("+ [1][Client] Marked to receive updates.")
	logger.Debugf("+ [1][Client] Mark has QoS(%d)", int(msg.QoS()))
	// (*self.Server).NotifySubscribe(topic, self)
	self.AddTopic(topic)
}

//
func (self *Client) GetTopics() []string {
	return self.Topics
}

//
func (self *Client) AddTopic(topic string) {
	// TODO
	self.Topics = append(self.Topics, topic)
}

//
func (self *Client) SetAuthMechanism() {
	// TODO
}

//
func (self *Client) GetIdentifier() string {
	return self.Username
}

func (self *Client) SetCreds(creds protobase.CredentialsInterface) {
	self.creds = creds
}

func (self *Client) GetCreds() protobase.CredentialsInterface {
	return self.creds
}

func (self *Client) SetUser(user interface{}) {
	self.User = user
}

func (self *Client) GetUser() interface{} {
	return self.User
}

func (self *Client) Setup() error {
	// TODO
	logger.FDebug("Setup", "* [Client/Setup] invoked.")
	return nil
}
