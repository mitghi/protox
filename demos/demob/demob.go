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

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/protocol"
//	"github.com/pkg/profile"
)

var (
	logger *logging.Logging

	username *string
	password *string
	pr       *string
	pm       *string
	addr     string
	cl       *client.Client
	qos      *int
	user     *User
)

func init() {
	logger = logging.NewLogger("DemoB")
}

func pcallback(opts protobase.OptionInterface, msg protobase.MsgInterface) {
	logger.Debug("+ [Publish/Callback] ++++ INSIDE PUBLISH PCALLBACK ++++")
}

func main() {
  //	defer profile.Start(profile.CPUProfile).Stop()
	log.Println("[+] started")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)

	username = flag.String("username", "", "username")
	password = flag.String("password", "", "password")
	qos = flag.Int("qos", 0, "quality of service")
	pr = flag.String("pr", "a/simple/route", "publish route")
	pm = flag.String("pm", "hello", "publish message")
	flag.Parse()

	if raddr := os.Getenv("PROTOX_ADDR"); raddr != "" {
		addr = raddr
	} else {
		addr = ":52909"
	}

	opts := client.CLBOptions{
		Addr:      addr,
		MaxRetry:  10,
		HeartBeat: 4,
		ClientDelegate: func() protobase.ClientInterface {
			cl := NewCustomClient(*username, *password, "")
			cl.SetCreds(&auth.Creds{*username, *password, ""})
			return cl
		},
		StorageDelegate: messages.NewMessageBox(),
		Conn:            protocol.NewClientConnection(addr),
		SecMRS:          2,
		CFCallback:      nil,
	}

	user = NewUser(opts)
	if user == nil {
		panic("unable to create user")
	}
	cl := user.Cl.(*CustomClient)
	cl.user = user
	if err := user.Setup(); err != nil {
		panic("unable to setup user")
	}

	user.Connect()
	<-sigs
	user.Disconnect()
	log.Println("received signal, exiting....")
	<-user.Exch
}
