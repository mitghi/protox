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

// Broker contains a sample broker and more.
package broker

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/server"
)

// Ensure interface (protocol) conformance.
var (
	_ protobase.BrokerInterface = (*Broker)(nil)
)

// NewBroker returns a configured Broker
// from provided Options.
func NewBroker(opts Options) protobase.BrokerInterface {
	var (
		ret *Broker = &Broker{}
		err error
	)
	if opts.ServerConf.Config != nil {
		ret.server, err = server.NewServerWithConfigs(opts.ServerConf)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else {
		ret.server = server.NewServer()
	}
	ret.exitch = ret.server.GetErrChan()
	if opts.Auth != nil {
		ret.server.SetAuthenticator(opts.Auth)
		ret.authsys = opts.Auth
	} else {
		ret.authsys = auth.NewAuthenticator()
		ret.server.SetAuthenticator(ret.authsys)
	}
	if opts.ConnectionDelegate != nil {
		ret.server.SetConnectionHandler(opts.ConnectionDelegate)
	} else {
		ret.server.SetConnectionHandler(ret.connectionDelegate)
	}
	if opts.HeartBeat != 0 {
		ret.heartbeat = opts.HeartBeat
		ret.server.SetHeartBeat(opts.HeartBeat)
	} else {
		ret.heartbeat = HEARTBEAT
		ret.server.SetHeartBeat(HEARTBEAT)
	}
	if opts.ClientDelegate != nil {
		ret.server.SetClientHandler(opts.ClientDelegate)
	} else {
		ret.server.SetClientHandler(ret.clientDelegate)
	}
	if opts.MsgStore != nil {
		ret.msgstore = opts.MsgStore
		ret.server.SetMessageStore(opts.MsgStore)
	} else {
		ret.msgstore = messages.NewInitedMessageStore()
		ret.server.SetMessageStore(ret.msgstore)
	}
	if opts.ClientStore != nil {
		ret.clientstore = opts.ClientStore
	} else {
		ret.clientstore = NewClientStore()
	}
	if opts.Exit != nil {
		ret.E = opts.Exit
	} else {
		ret.E = make(chan struct{})
	}
	if opts.ShutdownDeadline > 0 {
		ret.shwddln = opts.ShutdownDeadline
	} else {
		ret.shwddln = DSTDWN
	}
	ret.sigch = make(chan os.Signal, 1)
	signal.Notify(ret.sigch, syscall.SIGINT, syscall.SIGKILL)

	return ret
}

func (brk *Broker) RegisterClients(clients []*auth.Creds) {
	for _, v := range clients {
		brk.authsys.Register(v)
	}
}

func (brk *Broker) Status() byte {
	// TODO
	return 0x0
}

func (brk *Broker) handleSignals() {
	select {
	case <-brk.sigch:
		fmt.Printf("[X] received SIGINT, shutting down ....\n")
	case <-brk.exitch:
		fmt.Printf("[X] server exiting ( fatal ? ) .\n")
	}
	brk.Stop()
}

func (brk *Broker) Start() (ok bool) {
	var (
		serverStatus uint32
		statusChan   <-chan uint32
	)
	if atomic.LoadUint32(&brk.running) == BrokerRunning {
		return false
	} else if status := brk.server.GetStatus(); status != protobase.ServerNone {
		return false
	} else if atomic.LoadUint32(&brk.firstRun) == 1 {
		return false
	} else if atomic.LoadUint32(&brk.stopping) == BrokerStopping {
		return false
	}
	atomic.StoreUint32(&brk.firstRun, 1)
	logger.Info("[+] starting server....")
	// spawn handler coroutines
	go brk.handleSignals()
	go brk.server.ServeTCP(ADDR)
	statusChan = brk.server.GetStatusChan()
	serverStatus = <-statusChan
	switch serverStatus {
	case protobase.ServerRunning:
		atomic.StoreUint32(&brk.running, BrokerRunning)
		ok = true
	case protobase.ServerStopped:
		atomic.StoreUint32(&brk.running, BrokerStopping)
		ok = false
	default:
		ok = false
	}
	return ok
}

func (brk *Broker) corou(fn func()) {
	brk.wg.Add(1)
	go fn()
}

func (brk *Broker) Stop() bool {
	if atomic.LoadUint32(&brk.running) == BrokerNone {
		return false
	} else if atomic.LoadUint32(&brk.stopping) == BrokerStopping {
		return false
	}
	atomic.StoreUint32(&brk.stopping, 1)
	// handle statuses
	if stat := brk.server.GetStatus(); stat == protobase.ServerRunning {
		var (
			ch  <-chan struct{}
			err error
		)
		ch, err = brk.server.Shutdown()
		if err != nil {
			fmt.Println("[-----unable-to-shutdown-gracefully-----]")
			fmt.Println("[-] Shutdown failed.")
			atomic.StoreUint32(&brk.running, BrokerNone)
			select {
			case brk.E <- struct{}{}:
			default:
			}
			return false
		}
		// terminate with timeout
		select {
		case <-time.After(DSTDWN):
			fmt.Println("[-----unable-to-shutdown-before-timeout-----]")
			fmt.Println("[-] Shutdown failed.")
			atomic.StoreUint32(&brk.running, BrokerNone)
			select {
			case brk.E <- struct{}{}:
			default:
			}
			return false
		case <-ch:
			break
		}
	}
	fmt.Printf("[+] Shutdown completed.")

	select {
	case brk.E <- struct{}{}:
	default:
	}
	atomic.StoreUint32(&brk.running, BrokerNone)
	return true
}
