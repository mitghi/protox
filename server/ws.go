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

/**
*
* IMPORTANT
*
* this file is only for experimenting web compatibility.
*
**/

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// ServeWS is the main listening loop for serving protocol connections
// over WebSockets. Notice - this is experimental.
func (s *Server) ServeWS(address string) (err error) {
	http.Handle("/stat", websocket.Handler(wsHandler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("FAILED TO START WEBSOCKET")
		return err
	}

	return nil
}

// WsHandler is a connection handling routine used by `ServeWS`.
func wsHandler(ws *websocket.Conn) {
	reader, writer := io.Reader(ws), io.Writer(ws)
	for {
		buff := make([]byte, 1024)
		_, err := reader.Read(buff)
		if err != nil {
			log.Println("- error in reading line from websocket")
		}
		log.Println(buff)
		out := "this is working"
		b := []byte(out)
		writer.Write(b)
	}
}
