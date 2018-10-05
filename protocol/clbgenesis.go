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

/*
* TODO:
* . Add policies ( retry policy, ip policy, 
*   delay threshold, connection lag )
* . Implement this using state machine
*/

// Genesis is the initial and most important stage. All new connections can connect
// to broker iff they pass this stage. This stage only accepts `Connect` packets.
// Any other control packet results in immediate termination ( it can be adjusted using
// policies ). CGenesis is Genesis stage from client perspective. 
type CGenesis struct {
	constate

	Conn           *CLBConnection
  // TODO
  // . refactor into struct holding
  //   stage specific flags.
	gotFirstPacket bool
}

// NewCGenesis allocates and initializes `CGensis` state and
// returns the pointer, `nil` on failure. It implements the
// genesis procedure from client perspective, validates
// connection acknowledge packet and pushes the state into
// `COnline` when succesful.
func NewCGenesis(conn *CLBConnection) *CGenesis {
	var (
    cg *CGenesis =  &CGenesis{
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
  )
	cg.Conn = cg.constate.constatebase.Conn.(*CLBConnection)
	cg.client = conn.GetClient()
	if cg.Conn == nil {
		logger.FFatal("NewCGenesis", "- [FATAL|TODO]result.Conn==nil.")
    return nil
	}
	return cg
}

// SetNextState pushes the state machine into its next stage.
// Initially it is from CGenesis to Online ( CGenesis -> Online -> .... ).
func (cg *CGenesis) SetNextState() {
	newState := NewCOnline(cg.Conn)
	cg.Conn.State = newState

	logger.Debug("+ [CGenesis] CGenesis for client [Status] ready.")
}

// cleanUp is the deinitialization routine.
func (cg *CGenesis) cleanUp() {
  // TODO
  // . ensure garbage collection
  
  // no further mutation neccessary,
  // remove the connection struct provided from the server.
  // remove the associated client provided from the server.
	cg.Conn = nil
	cg.client = nil
}

// Shutdown terminates the current state and forwards the flow to
// handler routine, performing termination and undoing side effects.
func (cg *CGenesis) Shutdown() {
	// TODO
  // . prevent double calls
	go func() { cg.client.Disconnected(protobase.PUForceTerminate) }()
}

// Handle is a stub routine. It is to satisfy interface
// requirements ( for CGenesis stage ).
func (cg *CGenesis) Handle(packet *Packet) {
  // NOP
}

// onCONNACK parses 'connack' packet. It performs
// status validation check and pushes the connection 
// state into next stage iff status is assigned to 
// 'STATONLINE', indicating server accepted the 
// connection.
func (cg *CGenesis) onCONNACK(packet *Packet) {
  const fn string = "onCONNACK"
	var (
    p *Connack = NewConnack() // connection acknowledge packet
    caopt *ConnackOpts        // connection acknowledge options subpacket
  )  
	logger.FTrace(1, fn, "+ [ConnAck] packet received.")
	if err := p.DecodeFrom((*packet).Data); err != nil {
		logger.FTrace(1, fn, "- [Fatal] invalid connack packet.", err)
    // TODO
    // . forward to error handler
	}
	logger.FTrace(1, fn, "* [ConnAck] connack content.", p)
	if p.ResultCode == RESPFAIL {
    // TODO
    // . refactor into separate routine performing pattern matching
    //   and set the status from the provided table
    // Fail. invalid credentials.
		logger.FTrace(1, fn, "- [Credentials] INVALID CREDENTIALS.")
		cg.Conn.SetStatus(STATERR)
		return
	} else if p.ResultCode == RESPOK {
    // successfully validated.
		cg.Conn.SetStatus(STATONLINE)
		logger.FTrace(1, fn, "+ [Credentials] are valid and [Client] is now (Online).")
    // safe to remove connection state options
		for k, _ := range cg.Conn.stateOpts {
			delete(cg.Conn.stateOpts, k)
		}
    // deserialize the options and populate
    // the packet.
		caopt = NewConnackOpts()
		caopt.parseFrom(p)
    // store packet options for current CCONNACK state
		cg.Conn.stateOpts[CCONNACK] = caopt
    // push to next state
		cg.SetNextState()
    // clean current state
		cg.cleanUp()
		return
	} else {
		// NOTE:
    // TODO:
		// . THIS IS FATAL. check invalid codes
		logger.FTracef(1, fn, "- [Connack/Resp] unknown response-code(%b) in packet.", p.ResultCode)
    // clean current state
    cg.cleanUp()
		return
	}
}

// HandleDefault is the first function invoked in `CGenesis` when a
// new state struct is created. It passes credentials from `Connect`
// packet to a `AuthInterface` implementor and upgrades from
// `CGenesis` to `Online` stage. It sends a `Connack` with
// appropirate status code, regardless.
func (cg *CGenesis) HandleDefault(packet *Packet) (ok bool) {
	const fn = "HandleDefault"
	var (
		newcl   protobase.ClientInterface = cg.Conn.GetClient()
		p       *Connect                  = NewConnect()
		Conn    *CLBConnection            = cg.Conn
		nc      net.Conn
		err     error
		rpacket *Packet
		addr    string
		// by default, assume packet is invalid
		// valid   bool                      = false
	)
	//  defer clean up
	defer func() {
		if !ok {
      /* critical section */      
			Conn.protocon.Lock()
			if Conn.protocon.Conn != nil {
				Conn.protocon.Conn = nil
				Conn.protocon.Writer = nil
				Conn.protocon.Reader = nil
			}
			Conn.protocon.Unlock()
      /* critical section - end */
		}
	}()
  // get destination address
  /* critical section */
	Conn.protocon.RLock()  
	addr = Conn.protocon.addr
	Conn.protocon.RUnlock()
  /* critical section - end */
	if Conn.protocon.Conn == nil {
		// TODO
		// . pass connection options (tls) to dialRemoteAddr
		if c, err := dialRemoteAddr(addr, nil, false); err != nil {
			logger.FDebug(fn, "- [TCPConnect] cannot connect to remote addr.", "error", err)
			ok = false
			cg.client.Disconnected(protobase.PUSocketError)
			return
		} else {
			nc = c
		}
	}
	if nc == nil {
		ok = false
		cg.client.Disconnected(protobase.PUSocketError)
		return
	}
	Conn.protocon.Lock()
	Conn.SetNetConnection(nc)
	Conn.protocon.Unlock()

	p.Username, p.Password, p.ClientId = newcl.GetCreds().GetCredentials()
	if err = p.Encode(); err != nil {
		logger.FFatal("HandleDefault", "- [Encode] cannot encode in [CGenesis].", err)
		ok = false
		cg.client.Disconnected(protobase.PURejected)
		return
	}
	rpacket = p.GetPacket().(*Packet)
	if err = Conn.protocon.send(rpacket); err != nil {
		logger.Debug("- [Send] send returned an error.", err)
		ok = false
		cg.client.Disconnected(protobase.PUSocketError)
		return
	}

	ok = true
	return

	// TODO
	// . improve error handling
	// . do proper cleanup before exiting
	// cg.cleanUp()
}

func (cg *CGenesis) onDISCONNECT(packet *Packet) {
	// TODO
	logger.FDebug("onDISCONNECT", "* [Disconnect] packet received.")
	cg.Conn.protocon.Conn.Close()
}

func (cg *CGenesis) onPONG(packet *Packet) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}
