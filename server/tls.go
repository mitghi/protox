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

package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"time"

	"github.com/mitghi/protox/protobase"
)

// TODO

// ServeTLS is the main listening loop for serving clients over TLS sockets. It is
// responsible for accepting incoming TCP connection from clients. It is important
// to ensure that key paths are not non-existing.
func (s *Server) ServeTLS() (err error) {
	const fn = "ServeTCP"
	// TODO
	s.State.mode = ProtoTLS
	var (
		server net.Listener
		ticker *time.Ticker
	)

	// NOTE: IMPORTANT:
	// . opts must be checked for having appropirate Config type
	opts := s.opts.Config.(TLSOptions)
	tlsconfigs, err := s.generateTLSConfig(&opts)
	if err != nil {
		return err
	}
	server, err = tls.Listen("tcp", s.opts.Addr, tlsconfigs)
	defer server.Close()
	if err != nil {
		logger.Debug("- [Fatal] Cannot listen for incomming connections.")
		return err
	}
	s.listener = &server
	_ = s.SetStatus(protobase.ServerRunning)
	ticker = time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	defer s.SetStatus(protobase.ServerStopped)
	// defer s.corous.Done()
	go func() {
		for _ = range ticker.C {
			if stat := s.GetStatus(); stat == protobase.ForceShutdown {
				if err := (*s.listener).Close(); err != nil {
					logger.FError(fn, "- [TCP Handler] cannot close the server.", err)
				}
				logger.FDebug(fn, "* [Coro] waiting for coroutines ....")
				break
				// s.corous.Wait()
			}
		}
	}()
ML:
	for {
		var (
			conn net.Conn
		)
		conn, err = server.Accept()
		if err != nil {
			// Wait for all corous to finish
			// s.corous.Wait()
			logger.FDebug(fn, "Returning from TcpListener", "error")
			logger.FDebug(fn, "* [Coro] finished \t\t finished. ")
			logger.FDebug(fn, "{{ breaking ML }}")
			// tell other side of chan because shit hit the fan!
			// s.critical <- struct{}{}
			// return err
			break ML
		}
		logger.FInfo(fn, "* [Genesis] Participation request accepted.")
		s.corous.Add(1)
		go s.handleIncomingConnection(conn)
	}
	s.disconnectAll()
	// Wait for all corous to finish
	s.corous.Wait()
	// tell other side of chan because shit hit the fan!
	s.critical <- struct{}{}
	return err
}

// generateTLSConfig generates a tailored configuration for TLS server
// and returns a pointer to a `tls.Config`. It returns an error to indicate
// a problem.
func (s *Server) generateTLSConfig(opts *TLSOptions) (*tls.Config, error) {
	var (
		err      error
		cert     tls.Certificate
		config   *tls.Config
		certpool *x509.CertPool
	)
	cert, err = tls.LoadX509KeyPair(opts.Cert, opts.Key)
	if err != nil {
		logger.Error("no cert is given.", err)
		return nil, err
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		logger.Error("failed to parse certificate.")
		return nil, err
	}
	cert.Leaf = leaf
	if opts.Ciphers == nil {
		opts.Ciphers = make([]uint16, len(defaultCiphers))
		opts.Ciphers = append(opts.Ciphers, defaultCiphers...)
	}
	if opts.Curves == nil {
		opts.Curves = make([]tls.CurveID, len(defaultCurves))
		opts.Curves = append(opts.Curves, defaultCurves...)
	}
	// NOTE: `InsecureSkipVerify` option should be checked in testing (InsecureSkipVerify: true)
	config = &tls.Config{
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         opts.Curves,
		CipherSuites:             opts.Ciphers,
	}
	if opts.ShouldVerify {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	if opts.Ca != "" {
		certpool = x509.NewCertPool()
		rpem, err := ioutil.ReadFile(opts.Ca)
		if rpem == nil || err != nil {
			return nil, err
		}
		if !certpool.AppendCertsFromPEM([]byte(rpem)) {
			return nil, SRVTLSInvalidCA
		}
		config.ClientCAs = certpool
	}

	return config, nil
}
