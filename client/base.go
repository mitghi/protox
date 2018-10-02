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

	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
)

// TODO
// . add state

// logger is the logging facility
var logger protobase.LoggingInterface

func init() {
	logger = logging.NewLogger("Client")
}

// ClientErrorHandlerFunc is the callback for error handling.
type ClientErrorHandlerFunc func(client *ClientBase)

// ClientBase is the base Client interface.
type ClientBase interface {
	Connected(protobase.OptionInterface) bool
	Disconnected(protobase.OptCode)
	Publish(protobase.MsgInterface)
	Subscribe(protobase.MsgInterface)
	GetIdentifier() string
	GetTopics() []string
	GetCreds() protobase.CredentialsInterface
	SetCreds(protobase.CredentialsInterface)
	SetUser(interface{})
	GetUser() interface{}
	// TODO: NOTE:
	// . investigate addition of auth mechanism to Client struct
	// SetAuthMechanism()
}

// ServerInterface is the interface containing
// limited server methods for Client.
type ServerInterface interface {
	NotifyDisconnected(prc ProtoConnection)
	NotifyConnected(prc ProtoConnection)
	NotifySubscribe(msgbox protobase.MsgInterface, prc ProtoConnection)
	NotifyPublish(msg protobase.MsgInterface, prc ProtoConnection)
	RegisterClient(prc ProtoConnection)
}

// ProtoConnection is the interface containing
// limited connection methods for Client.
type ProtoConnection interface {
	Handle()
	SetServer(sv ServerInterface)
}

// Client is the structure for high-level and client-side
// logics.
type Client struct {
	*sync.RWMutex

	creds       protobase.CredentialsInterface
	User        interface{}
	ClientId    string
	Username    string
	Password    string
	WillMessage string
	WillQos     string
	WillRetain  string
	Server      ServerInterface
	Topics      []string
	// Conn        protocol.Connection
	// TODO:
	// . make a new interface for this client
	// . add group and userRole and make the interface
	//   compatible.
	// e.g.:
	// group       string
	// userRole    protobase.PermissionInterface
}
