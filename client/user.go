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
	"errors"
	"sync"
	"time"

	"github.com/mitghi/protox/protobase"
)

/**
* simple user implementation
**/

// Interface check
var _ protobase.CLBUserInterface = (*CLBUser)(nil)

// Error messages
var (
	// TODO
	// . add individual error messages for each option
	CLBUserInvalid error = errors.New("CLBUser: invalid.")
)

type CLBUser struct {
	sync.RWMutex

	Conn       protobase.ProtoClientConnection
	Cl         protobase.ClientInterface
	Storage    protobase.MessageBox
	Exch       chan struct{}
	exconnch   chan struct{}
	CFCallback func(*CLBUser)
	Addr       string
	SecMRS     int
	MaxRetry   int
	HeartBeat  int
	hadSetup   bool
	Running    bool
	Connected  bool
}

type CLBOptions struct {
	Addr            string
	MaxRetry        int
	HeartBeat       int
	ClientDelegate  func() protobase.ClientInterface
	StorageDelegate protobase.MessageBox
	Conn            protobase.ProtoClientConnection
	SecMRS          int // Seconds of Maximum Retry Sleep
	CFCallback      func(*CLBUser)
}

func checkOpts(opts CLBOptions) bool {
	if opts.Addr == "" {
		return false
	} else if opts.ClientDelegate == nil {
		return false
	} else if opts.Conn == nil {
		return false
	} else if opts.StorageDelegate == nil {
		return false
	}
	return true
}

func NewCLBUser(opts CLBOptions) (*CLBUser, bool) {
	if !checkOpts(opts) {
		return nil, false
	}

	u := &CLBUser{
		Exch:       make(chan struct{}, 1),
		exconnch:   make(chan struct{}),
		Running:    false,
		Conn:       opts.Conn,
		Addr:       opts.Addr,
		CFCallback: opts.CFCallback,
		SecMRS:     opts.SecMRS,
		Cl:         opts.ClientDelegate(),
		HeartBeat:  opts.HeartBeat,
		Storage:    opts.StorageDelegate,
		hadSetup:   false,
		Connected:  false,
	}

	return u, true
}

func (u *CLBUser) Setup() error {
	if u.Addr == "" {
		return CLBUserInvalid
	} else if u.Cl == nil {
		return CLBUserInvalid
	} else if u.Conn == nil {
		return CLBUserInvalid
	} else if u.Storage == nil {
		return CLBUserInvalid
	} else if u.hadSetup {
		return CLBUserInvalid
	}
	u.Conn.SetClient(u.Cl)
	u.Conn.SetMessageStorage(u.Storage)
	if u.HeartBeat >= 1 {
		u.Conn.SetHeartBeat(u.HeartBeat)
	}
	// NOTE
	// . check correctness
	u.hadSetup = true

	return nil
}

func (u *CLBUser) IsRunning() (ret bool) {
	u.RLock()
	defer u.RUnlock()
	ret = u.Running
	return ret
}

func (u *CLBUser) IsConnected() (ret bool) {
	u.RLock()
	defer u.RUnlock()
	ret = u.Connected
	return ret
}

func (u *CLBUser) GetExitCh() chan struct{} {
	return u.Exch
}

func (u *CLBUser) SetConnected(b bool) {
	u.Lock()
	u.Connected = b
	u.Unlock()
}

func (u *CLBUser) SetRunning(b bool) {
	u.Lock()
	u.Running = b
	u.Unlock()
}

func (u *CLBUser) Disconnect() {
	// TODO
	// . check for closed channel
	// . return error code
	// NOTE
	// . DONOT CALL THIS RECEIVER MANUALLY
	if u.IsRunning() {
		err := u.Conn.Disconnect()
		logger.Debugf("* [Client/User(CLBConnector)] Disconnecting. Error code (%+v).", err)
		// call delegate method manually ( because its forcefully
		// terminating either because of SIGNALS or FATAL conditions.
		u.Cl.Disconnected(protobase.PUForceTerminate)
		u.exconnch <- struct{}{}
	}
}

func (u *CLBUser) Connect() {
	// TODO
	// . return error code
	u.SetRunning(true)
	go func() {
		clexch := make(chan struct{}, 1)
		go func() {
			termch := u.Conn.GetTermChan()
			<-clexch
			select {
			case termch <- struct{}{}:
			default:
			}
			u.Conn.ContinueFlag(false)
			u.SetRunning(false)
		}()
		go func() {
			var (
				// TODO:
				// . refactor and set sleep cycles
				//   by corresponding receiver
				//   methods.
				minslp = time.Millisecond * 10
				maxslp = time.Second * time.Duration(u.SecMRS)
				dur    = minslp
			)
			// TODO
			// . add maximum retry
			for u.IsRunning() {
				u.Conn.Handle()
				time.Sleep(dur)
				dur *= 2
				if dur > maxslp {
					dur = minslp
				}
			}
			// TODO
			// . send disconnect packet
		}()
		// <-u.Exch
		<-u.exconnch
		clexch <- struct{}{}
		u.Exch <- struct{}{}
	}()
}
