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
	"errors"
	"net"
	"sync"

	/* will be replaced with new lock-free implementation */
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

type ServerType byte

// Constants for server service types (Protox over : HTTP/S, MQTT, WS, REDIS PROTOCOL, .... )
const (
	// Protox over original protocol
	ProtoOrigin ServerType = iota
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

// Defaults
var (
	DefaultHeartbeat int = 1
)

// Server is a main implementation of `protocol.ServerInterface`.
type Server struct {
	sync.RWMutex
	protobase.ServerInterface

	corous             sync.WaitGroup
	Clients            map[net.Conn]protobase.ProtoConnection
	onNewClient        func(string, string, string) protobase.ClientInterface
	onNewConnection    ConnectionDelegate
	onNewMessage       ServerHandlerFunc
	permissionDelegate func(protobase.AuthInterface, ...string) bool
	Authenticator      protobase.AuthInterface
	Store              protobase.MessageStorage
	Router             protobase.RouterInterface
	State              *serverState
	listener           *net.Listener
	opts               *ServerConfigs
	StatusChan         chan uint32
	critical           chan struct{}
	heartbeat          int
	Status             uint32
	// TODO: NOTE:
	//  . add timestamp and expiration date.
	//  . for heartbeat, uint can also be used.
}
