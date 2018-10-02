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

// Server contains bunch of interfaces and functionalities for serving
// connections compatible with Protox. It has APIs for interacting with
// the underlaying subsystems.
package server

/*
* TODO:
* . decouple Buffer Pool from direct import 
*/

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"

	buffpool "github.com/mitghi/lfpool" /* will be replaced with new lock-free implementation */
	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
)

// Logger is the default, structured package level logging service.
var (
	logger protobase.LoggingInterface
)

func init() {
	logger = logging.NewLogger("Server")
}

// Constants for server client types
const (
	// Client
	STCLIENT SConnType = iota
	// Manager
	STMANAGER
	// Router
	STROUTER
	// Resource provider
	STRESOURCE
	// Monitor and statics collector
	STMONITOR
)

// Constants for server service types (Protox over : HTTP/S, MQTT, WS, REDIS PROTOCOL, .... )
const (
	// Protox over original protocol
	ProtoOrigin byte = iota
	// Protox over HTTP
	ProtoHTTP
	// Protox over HTTPS
	ProtoHTTPS
	// Protox over WS
	ProtoWS
	// Protox over MQTT
	ProtoMQTT
	// Protox over RP
	ProtoRP
)

// Constants for server transport types (UNIX, TCP, SSL/TLS, .... )
const (
	ProtoUNIXSO byte = iota
	// Protox over TCP
	ProtoTCP
	// Protox over TLS
	ProtoTLS
	// Protox over SSL
	ProtoSSL
)

// Server error messages
var (
	SRVShutdownError  error = errors.New("server: cannot shutdown due to incompatible state.")
	SRVGeneralError   error = errors.New("server: %s")
	SRVInvalidAddr    error = errors.New("server: invalid address.")
	SRVInvalidMode    error = errors.New("server: invalid serving mode.")
	SRVMissingOptions error = errors.New("server: options are missing.")
	SRVTLSInvalidCA   error = errors.New("server: invalid caFile.")
)

// SConnTyp is server client type ( CLIENT, RESOURCE, ROUTER, MONITOR, .... )
type SConnType byte

// ServerHandlerFunc is the signature for `ClientInterface` delegation.
type ServerHandlerFunc func(net.Conn) protobase.ClientInterface

// ClientDelegate is a signature for `ClientInterface` delegation.
type ClientDelegate func(string, string, string) protobase.ClientInterface

// ConnectionDelegate is a signature for `ProtoConnection` delegation.
type ConnectionDelegate func(net.Conn) protobase.ProtoConnection

// Server is a main implementation of `protocol.ServerInterface`.
type Server struct {
	sync.RWMutex
	protobase.ServerInterface // explicitly implement interface

	Clients            map[net.Conn]protobase.ProtoConnection                 // section routing
	Router             map[string]map[string]protobase.ProtoConnection        // end
	onNewClient        func(string, string, string) protobase.ClientInterface //section delegates
	onNewConnection    ConnectionDelegate
	onNewMessage       ServerHandlerFunc
	permissionDelegate func(protobase.AuthInterface, ...string) bool
	Authenticator      protobase.AuthInterface
	Store              protobase.MessageStorage                               // end
	State              *serverState
	rt                 *Router
	listener           *net.Listener                                          //end
	buffer             *buffpool.BuffPool
	opts               *ServerConfigs
	corous             sync.WaitGroup                                         // coroutines
	heartbeat          int                                                    // state section
  Status             uint32
  StatusChan         chan uint32
	critical           chan struct{} // end. Critical/Fatal error
	// TODO: NOTE:
	//  . add timestamp and expiration date.
	//  . for heartbeat, uint can also be used.
}

// UNIXSOPtions contains necessary information required by unix socket server.
type UNIXSOptions struct {
	// TODO
}

// TCPOptions contains necessarry information required by TCP server.
type TCPOptions struct {
	// TODO
}

// TLSOptions contains neccessary information required by TLS server.
type TLSOptions struct {
	Curves           []tls.CurveID
	Ciphers          []uint16
	Cert             string
	Key              string
	Ca               string
	HandshakeTimeout int
	ShouldVerify     bool
}

// ServerConfigs is a struct for configuring the server.
type ServerConfigs struct {
	Config interface{}
	Addr   string
	Mode   byte
	Type   byte
	// TRate is the cycle interval in milliseconds 
	// for performing status check.
	TRate int
	// TODO
	// . add callbacks
}
