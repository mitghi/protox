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

package networking

import (
	"github.com/mitghi/protox/protobase"
)

// MARK: Genesis

// Genesis is the initial and most important stage.
// All new connections can connect to broker iff
// they pass this stage. This stage only accepts
// `Connect` packets. Any other control packet results
// in immediate termination ( it can be adjusted using
// policies.
type Genesis struct {
	constate

	Conn           *Connection
	gotFirstPacket bool
}

// MARK: Genesis

// NewGenesis creates a pointer to a new `Gensis` packet.
func NewGenesis(conn *Connection) *Genesis {
	result := &Genesis{
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
	result.Conn = result.constate.constatebase.Conn.(*Connection)
	if result.Conn == nil {
		// TODO
		// . this is critical.
	}
	return result
}

// SetNextState pushes the state machine into its next stage.
// Initially it is from Genesis to Online ( Genesis -> Online -> .... ).
func (g *Genesis) SetNextState() {
	newState := NewOnline(g.Conn)
	newState.SetClient(g.client)
	newState.SetServer(g.server)
	g.Conn.State = newState

	logger.Debug("+ [Genesis] Genesis for client [Status] ready.")
}

// cleanUp is a routine which removes pointers from the struct.
func (g *Genesis) cleanUp() {
	g.Conn = nil
	g.client = nil
	g.server = nil
}

// Shutdown terminates the state and calls the handlers to terminate
// and undo all side effects.
func (g *Genesis) Shutdown() {
	g.client.Disconnected(protobase.PUForceTerminate)
}

// Handle is only a stub to satisfy interface requirements ( for Genesis stage ).
func (g *Genesis) Handle(packet protobase.PacketInterface) {
}

// HandleDefault is the first function invoked in `Genesis` when a new state struct is created.
// It passes credentials from `Connect` packet to a `AuthInterface` implementor and upgrades
// from `Genesis` to `Online` stage. It sends a `Connack` with appropirate status code, regardless.
func (g *Genesis) HandleDefault(packet protobase.PacketInterface) (status bool) {
	// TODO:
	// . add defer to cleanUp and check its performance impact
	var (
		// by default, packet is invalid
		valid   bool                    = false
		p       *Connect                = NewConnect(packet)
		cack    *Connack                = NewRawConnack()
		authsys protobase.AuthInterface = g.Conn.GetAuthenticator()
		creds   protobase.CredentialsInterface
		rpacket *Packet
		newcl   protobase.ClientInterface
	)
	logger.FDebug("HandleDefault", "* [Packet] content of raw packet.", packet.GetData())
	// terminate immediately if packet is malformed or invalid.
	if p == nil {
		logger.Debug("- [Fatal] invalid connection packet in [Genesis].", packet)
		// TODO
		//  undo side effects
		g.gotFirstPacket = false
		g.cleanUp()
		return false
	}
	// if err := p.DecodeFrom(packet.GetData()); err != nil {
	// 	logger.Debug("- [Fatal] invalid connection packet in [Genesis].", err)
	// 	// TODO
	// 	//  undo side effects
	// 	g.gotFirstPacket = false
	// 	g.cleanUp()
	// 	return false
	// }
	logger.FDebug("HandleDefault", "* [Packet] connection packet content.", p.String())
	// connection is established
	// push into the next state
	g.gotFirstPacket = true
	// TODO:
	// . improve by directly pass connect packet to auth subsystem
	creds, err := authsys.MakeCreds(p.Username, p.Password, p.ClientId)
	if err != nil {
		logger.Fatal("- [Fatal] cannot make credentials in [Genesis].", err)
		return false
	}
	// TODO/NOTICE
	//  do not create a new client until credentials are valid ( reduce memory alloc. overhead )
	newcl = g.Conn.clientDelegate(p.Username, p.Password, p.ClientId)
	// NOTE:
	// . check error explicitely
	if valid, err = authsys.CanAuthenticate(creds); valid {
		cack.SetResultCode(protobase.RESPOK)
		cack.Encode()
		rpacket = cack.GetPacket().(*Packet)
		g.Conn.SetClient(newcl)
		g.SetNextState() // Genesis -> Online
		if p.Meta.CleanStart {
			// drop queued packets
			g.Conn.storage.AddClient(p.Username)
			// TODO
			// . drop subscriptions
		}
		g.Conn.SendDirect(rpacket)
		// TODO
		//  these lines are moves to cleanUp, remove them when
		//  its finalized.
		/* d e b u g */
		//  g.Conn = nil
		//  g.client = nil
		/* d e b u g */
		g.cleanUp()
		return true
	} else {
		cack.SetResultCode(protobase.RESPFAIL)
		cack.Encode()
		rpacket = cack.GetPacket().(*Packet)
		g.Conn.SetClient(newcl)
		g.Conn.SendDirect(rpacket)
		// TODO
		//  improve error handling
		g.cleanUp()
		return false
	}
}

func (g *Genesis) OnPONG(packet protobase.PacketInterface) {
	// TODO
	logger.FDebug("onPONG", "* [Pong] packet received.")
}
