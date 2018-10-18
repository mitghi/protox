package server

import (
	"crypto/tls"
	"net"
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

func precheckOpts(opts *ServerConfigs) error {
	// TODO
	// . precheck server options before proceeding
	_, _, err := net.SplitHostPort(opts.Addr)
	if err != nil {
		return SRVInvalidAddr
	}
	if opts.Config != nil {
		switch opts.Mode {
		case ProtoTCP:
			if _, ok := opts.Config.(TCPOptions); !ok {
				return SRVMissingOptions
			}
			return nil
		case ProtoTLS:
			if _, ok := opts.Config.(TLSOptions); !ok {
				return SRVMissingOptions
			}
			return nil
		case ProtoSSL:
			// TODO
			return SRVInvalidMode
		case ProtoUNIXSO:
			// TODO
			return SRVInvalidMode
		default:
			return SRVInvalidMode
		}
	}
	return nil
}
