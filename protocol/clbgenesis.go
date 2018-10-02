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
	"net"

	"github.com/mitghi/protox/protobase"
)

// Genesis is the initial and most important stage. All new connections can connect
// to broker iff they pass this stage. This stage only accepts `Connect` packets.
// Any other control packet results in immediate termination ( it can be adjusted using
// policies.
type CGenesis struct {
	constate

	Conn           *CLBConnection
	gotFirstPacket bool
}

// NewGenesis creates a pointer to a new `Gensis` packet.
func NewCGenesis(conn *CLBConnection) *CGenesis {
	result := &CGenesis{
		constate: constate{
			constatebase: constatebase{
				Conn: conn,
			},
			client: nil,
			server: nil,
		},
		Conn:           nil,
		gotFirstPacket: false,
	}
	result.Conn = result.constate.constatebase.Conn.(*CLBConnection)
	result.client = conn.GetClient()
	if result.Conn == nil {
		logger.FFatal("NewCGenesis", "- [FATAL|TODO]result.Conn==nil.")
		// TODO
		// . this is critical.
	}
	return result
}

// SetNextState pushes the state machine into its next stage.
// Initially it is from CGenesis to Online ( CGenesis -> Online -> .... ).
func (self *CGenesis) SetNextState() {
	newState := NewCOnline(self.Conn)
	self.Conn.State = newState

	logger.Debug("+ [CGenesis] CGenesis for client [Status] ready.")
}

// cleanUp is a routine which removes pointers from the struct.
func (self *CGenesis) cleanUp() {
	self.Conn = nil
	self.client = nil
}

// Shutdown terminates the state and calls the handlers to terminate
// and undo all side effects.
func (self *CGenesis) Shutdown() {
	// TODO
	go func() { self.client.Disconnected(protobase.PUForceTerminate) }()
}

// Handle is only a stub to satisfy interface requirements ( for CGenesis stage ).
func (self *CGenesis) Handle(packet *Packet) {
}

func (self *CGenesis) onCONNACK(packet *Packet) {
	logger.FTrace(1, "onCONNACK", "+ [ConnAck] packet received.")
	var p *Connack = NewConnack()
	if err := p.DecodeFrom((*packet).Data); err != nil {
		logger.FTrace(1, "onCONNACK", "- [Fatal] invalid connack packet.", err)
	}
	logger.FTrace(1, "onCONNACK", "* [ConnAck] connack content.", p)
	// NOTE
	// . this has changed
	// if p.ResultCode == TMP_RESPOK {
	if p.ResultCode == RESPFAIL {
		logger.FTrace(1, "onCONNACK", "- [Credentials] INVALID CREDENTIALS.")
		self.Conn.SetStatus(STATERR)
		return
	} else if p.ResultCode == RESPOK {
		self.Conn.SetStatus(STATONLINE)
		logger.FTrace(1, "onCONNACK", "+ [Credentials] are valid and [Client] is now (Online).")
		// set connack result to options
		// and pass to to next state.
		for k, _ := range self.Conn.stateOpts {
			delete(self.Conn.stateOpts, k)
		}
		caopt := NewConnackOpts()
		caopt.parseFrom(p)
		self.Conn.stateOpts[CCONNACK] = caopt
		self.SetNextState()
		self.cleanUp()
		return
	} else {
		// NOTE: TODO:
		// . THIS IS FATAL, check invalid codes
		logger.FTracef(1, "onCONNACK", "- [Connack/Resp] unknown response-code(%b) in packet.", p.ResultCode)
		return
	}
}

// HandleDefault is the first function invoked in `CGenesis` when a new state struct is created.
// It passes credentials from `Connect` packet to a `AuthInterface` implementor and upgrades
// from `CGenesis` to `Online` stage. It sends a `Connack` with appropirate status code, regardless.
func (self *CGenesis) HandleDefault(packet *Packet) (ok bool) {
	const fn = "HandleDefault"
	// TODO
	//  add defer to cleanUp and check its performance impact
	var (
		newcl   protobase.ClientInterface = self.Conn.GetClient()
		p       *Connect                  = NewConnect()
		Conn    *CLBConnection            = self.Conn
		nc      net.Conn
		err     error
		rpacket *Packet
		addr    string
		// by default, assume packet is invalid
		// valid   bool                      = false
	)
	defer func() {
		if !ok {
			Conn.protocon.Lock()
			if Conn.protocon.Conn != nil {
				Conn.protocon.Conn = nil
				Conn.protocon.Writer = nil
				Conn.protocon.Reader = nil
			}
			Conn.protocon.Unlock()
		}
	}()

	Conn.protocon.RLock()
	addr = Conn.protocon.addr
	Conn.protocon.RUnlock()
	if Conn.protocon.Conn == nil {
		// TODO
		// . pass connection options (tls) to dialRemoteAddr
		if c, err := dialRemoteAddr(addr, nil, false); err != nil {
			logger.FDebug(fn, "- [TCPConnect] cannot connect to remote addr.", "error", err)
			ok = false
			self.client.Disconnected(protobase.PUSocketError)
			return
		} else {
			nc = c
		}
	}
	if nc == nil {
		ok = false
		self.client.Disconnected(protobase.PUSocketError)
		return
	}
	Conn.protocon.Lock()
	Conn.SetNetConnection(nc)
	Conn.protocon.Unlock()

	p.Username, p.Password, p.ClientId = newcl.GetCreds().GetCredentials()
	if err = p.Encode(); err != nil {
		logger.FFatal("HandleDefault", "- [Encode] cannot encode in [CGenesis].", err)
		ok = false
		self.client.Disconnected(protobase.PURejected)
		return
	}
	rpacket = p.GetPacket().(*Packet)
	if err = Conn.protocon.send(rpacket); err != nil {
		logger.Debug("- [Send] send returned an error.", err)
		ok = false
		self.client.Disconnected(protobase.PUSocketError)
		return
	}

	ok = true
	return

	// TODO
	// . improve error handling
	// . do proper cleanup before exiting
	// self.cleanUp()
}

func (self *CGenesis) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", "* [Disconnect] packet received.")
	self.Conn.protocon.Conn.Close()
}

func (self *CGenesis) onPONG(packet *Packet) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}
