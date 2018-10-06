package server

import (
	"crypto/tls"  
)

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
