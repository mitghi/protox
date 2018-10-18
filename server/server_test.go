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
	"net"
	"testing"
)

const (
	logmode bool = false
)

func TestNewServerRouter(t *testing.T) {
	var (
		s   *Server = NewServer()
		err error
	)
	s.Router.Add("test", "a/simple/topic", 10)
	m, err := s.Router.Find("a/simple/topic")
	if err != nil {
		t.Fatal("failed to find topic.", err)
	}
	if logmode {
		t.Log("router: associated map to the supplied topic.", m)
	}
	err = s.Router.Remove("test", "a/simple/topic")
	if err != nil {
		t.Fatal("failed to remove topic.", err)
	}
}

func TestServeTCP(t *testing.T) {
	var (
		s        *Server
		listener net.Listener
		err      error
	)
	s, err = NewServerWithConfigs(ServerConfigs{
		Config: TLSOptions{
			Cert: "/playground/cert/server.pem",
			Key:  "/playground/cert/key.pem",
		},
		Addr: "0.0.0.0:52909",
		Mode: ProtoTLS,
	})
	if err != nil {
		t.Fatal(cERR, err)
	}
	listener, err = s.serverInstance(":52909")
	if err != nil {
		t.Fatal(cERR, err)
	}
	if logmode {
		t.Log(listener)
	}
	// s.SetClientHandler(func(user, pass, uid string) protobase.ClientInterface {
	// 	fmt.Println("new client")
	// 	return nil
	// })
	// go s.ServeTCP(":52909")
	// time.Sleep(time.Second * 4)
	// ch, err := s.Shutdown()
	// if err != nil {
	// 	t.Fatal(cERR, err)
	// }
	// <-ch
	// fmt.Println("server stopped.")
}
