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
	"time"

	"github.com/mitghi/protox/protobase"
)

// ServeTCP is the main listening loop. It is responsible for
// accepting incoming TCP connections from clients.
func (s *Server) ServeTCP(address string) (err error) {
	const fn = "ServeTCP"
	s.State.mode = ProtoTCP
	var (
		server net.Listener
		ticker *time.Ticker
	)
	server, err = net.Listen("tcp", address)
	if err != nil {
		logger.Debug("- [Fatal] Cannot listen for incomming connections.")
    s.StatusChan <- protobase.ServerStopped
    _ = s.SetStatus(protobase.ServerStopped)    
		return err
	}
	defer server.Close()
	s.listener = &server
  s.StatusChan <- protobase.ServerRunning
	_ = s.SetStatus(protobase.ServerRunning)  
	ticker = time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	defer s.SetStatus(protobase.ServerStopped)
  /* d e b u g */
	// defer s.corous.Done()
  /* d e b u g */  
	go func() {
		for _ = range ticker.C {
			if stat := s.GetStatus(); stat == protobase.ForceShutdown {
				if err := (*s.listener).Close(); err != nil {
					logger.FError(fn, "- [TCP Handler] cannot close the server.", err)
				}
				logger.FDebug(fn, "* [Coro] waiting for coroutines ....")
				break
        /* d e b u g */        
				// s.corous.Wait()
        /* d e b u g */        
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
      /* d e b u g */      
			// Wait for all corous to finish
			// s.corous.Wait()
      /* d e b u g */      
			logger.FDebug(fn, "* [Coro] Returning from TcpListener", "error")
			logger.FDebug(fn, "* [Coro] finished \t\t finished. ")
			logger.FDebug(fn, "* [Coro] {{ breaking ML }}")
      /* d e b u g */      
			// tell other side of chan because shit hit the fan!
			// s.critical <- struct{}{}
			// return err
      /* d e b u g */      
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
