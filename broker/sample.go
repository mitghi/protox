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

package broker

import (
	"net"

	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
)

func NewClientStore() *ClientStore {
	return &ClientStore{clients: make(map[string]protobase.ClientInterface)}
}

// Add updates/adds a client `cid` to its struct pointer.
func (cls *ClientStore) Add(cid string, ptr protobase.ClientInterface) {
	cls.Lock()
	cls.clients[cid] = ptr
	cls.Unlock()
}

// Get fetches a client `cid` from the mapping. It returns a null
// pointer when client doesn't exist.
func (cls *ClientStore) Get(cid string) protobase.ClientInterface {
	cls.RLock()
	v, ok := cls.clients[cid]
	cls.RUnlock()
	if ok == true && v != nil {
		return v
	}
	return nil
}

// ClientDelegate creates a new handler for each new client and returns a structure
// compatible with `protocol.ClientInterface`. Most of high-level business logic should
// be implemented by customizing/providing a compatible `protocol.ClientInterface` structure.
// Protocol notifications such as Subscribe, Publish, Disconnect, Presence, Request, Broadcast
// and Proposals are delivered by calling delegate routines on the structure returned by this function.
// It will reuse the memory if a client struct is already in the storage, otherwise it allocate and
// returns a new one.
func (self *Broker) clientDelegate(username string, password string, cid string) protobase.ClientInterface {
	logger.Debug("* [clientDelegate] client ", username, " joined.")
	if ret := self.clientstore.Get(username); ret != nil {
		logger.Debug("** [clientDelegate] reusing existing struct for client: ", username)
		return ret
	}
	var cl *client.Client = client.NewClient(username, password, cid)
	self.clientstore.Add(username, cl)
	return cl
}

// ConnectionDeleagte creates a new connection for each new client
// and returns a compatible structure with `protocol.ProtoConnection` interface.
func (self *Broker) connectionDelegate(cl net.Conn) protobase.ProtoConnection {
	var proto *protocol.Connection = protocol.NewConnection(cl)
	return proto
}
