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

// Ensure protocol (interface) conformance.
var _ protobase.ClientInterface = (*Client)(nil)

func NewClient(username string, password string, cid string) *Client {
	client := &Client{
		RWMutex:     &sync.RWMutex{},
		ClientId:    cid,
		WillMessage: "",
		WillQos:     "",
		WillRetain:  "",
		Username:    username,
		Password:    password,
	}
	return client
}

func (c *Client) SetServer(server protobase.ServerInterface) {
	c.Server = server
}

func (c *Client) Connected(opts protobase.OptionInterface) bool {
	logger.Debug("+ [Client] Connected.")
	/* d e b u g */
	// (*c.Server).NotifyConnected(c)
	/* d e b u g */
	return true
}

func (c *Client) Disconnected(opts protobase.OptCode) {
	/* d e b u g */
	// (*c.Server).NotifyDisconnected(c)
	/* d e b u g */
	logger.Debug("- [Client] Disconnected.")
}

func (c *Client) Publish(msg protobase.MsgInterface) {
	switch msg.Dir() {
	case protobase.MDInbound:
		logger.Debug("+ [Client] Marked to notify status.")
	case protobase.MDOutbound:
		logger.Debugf("+ [Client] Publishing a message with QoS(%d).", int(msg.QoS()))
	}
}

func (c *Client) Subscribe(msg protobase.MsgInterface) {
	topic := msg.Envelope().Route()
	logger.Debug("+ [1][Client] Marked to receive updates.")
	logger.Debugf("+ [1][Client] Mark has QoS(%d)", int(msg.QoS()))
	/* d e b u g */
	// (*c.Server).NotifySubscribe(topic, c)
	/* d e b u g */
	c.AddTopic(topic)
}

//
func (c *Client) GetTopics() []string {
	return c.Topics
}

//
func (c *Client) AddTopic(topic string) {
	// TODO
	c.Topics = append(c.Topics, topic)
}

//
func (c *Client) SetAuthMechanism() {
	// TODO
}

//
func (c *Client) GetIdentifier() string {
	return c.Username
}

func (c *Client) SetCreds(creds protobase.CredentialsInterface) {
	c.creds = creds
}

func (c *Client) GetCreds() protobase.CredentialsInterface {
	return c.creds
}

func (c *Client) SetUser(user interface{}) {
	c.User = user
}

func (c *Client) GetUser() interface{} {
	return c.User
}

func (c *Client) Setup() error {
	// TODO
	logger.FDebug("Setup", "* [Client/Setup] invoked.")
	return nil
}
