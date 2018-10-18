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
)

// TODO

// generateTLSConfig generates a tailored configuration for TLS server
// and returns a pointer to a `tls.Config`. It returns an error to indicate
// a problem.
func (s *Server) generateTLSConfig(opts *TLSOptions) (*tls.Config, error) {
	return generateTLSConfig(opts)
}

func generateTLSConfig(opts *TLSOptions) (*tls.Config, error) {
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
