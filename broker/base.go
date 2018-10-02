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

// package broker provides Message broker implementation.
package broker

import (
	"os"
	"sync"
	"time"

	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/server"
)

// Broker status flags
const (
	BrokerNone uint32 = iota
	BrokerRunning
	BrokerStopping
)

// Default port is 0xcead.
const ADDR = ":52909" // :0xcead

var (
	// Default conneciton heartbeat
	HEARTBEAT int = 5
	// Exit timeout
	// Server will forcefully terminated
	// if it fails to exit before this
	// deadline.
	DSTDWN time.Duration = time.Second * 5
	// Logger provides logging facilities.
	logger protobase.LoggingInterface
)

// Init is the package level initializor.
func init() {
	logger = logging.NewLogger("Broker")
}

// ClientStore holds reference to allocated
// `protobase.ClientInterface` implementors.
// It reuses the structure in the mapping,
// if a client exists.
type ClientStore struct {
	sync.RWMutex
	clients map[string]protobase.ClientInterface
}

// Options holds configuration detail.
type Options struct {
	HeartBeat           int
	Auth                protobase.AuthInterface
	MsgStore            protobase.MessageStorage
	ClientStore         protobase.CLStoreInterface
	ClientDelegate      server.ClientDelegate
	ConnectionDelegate  server.ConnectionDelegate
  ShutdownDeadline    time.Duration
  Exit                chan struct{}
}

// TODO
// . tracking

// stats provides statistics for each
// connection including.
type stats struct {
	Conns uint64
	Send  uint64
	Recv  uint64
}

// TODO
// . create broker interface

// Broker implements a message broker.
type Broker struct {
	wg          sync.WaitGroup 
	server      *server.Server             // serving subsystem
	authsys     protobase.AuthInterface    // authentication subsystem
	msgstore    protobase.MessageStorage   // storage holding message data and metadata 
	clientstore protobase.CLStoreInterface // storage holding client data
	start       time.Time                  // startup delay
	shwddln     time.Duration              // maximum tolerable time for shutdown procedure
	opts        *Options                   // options
	heartbeat   int                        // maximum tolerable time for connection health check
	firstRun    uint32                     // initial startup flag
	running     uint32                     // running status flag
	stopping    uint32                     // stopping procedure flag
	exitch      <-chan struct{}            // exit channel 
	sigch       chan os.Signal
	E           chan struct{}
}
