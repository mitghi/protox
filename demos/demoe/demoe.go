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
	"os"

	"fmt"

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/broker"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/server"
)


var (
	logger *logging.Logging
)

func init(){
	logger = logging.NewLogger("DemoE")
}

// clientDelegate is the delegate used by server to
// create a compatible `protobase.ClientInterface`
// struct.
func clientDelegate(uid string, pid string, cid string) protobase.ClientInterface {
	ret := User{client.Client{Username: uid, Password: pid, ClientId: cid}}
	return &ret
}

// defaultAuthConfig returns default configuration
// of Auth subsystem ( Strict Mode ) with a few
// dummy clients.
func defaultAuthConfig() *auth.AuthConfig {
	// initialize auth configuration container
	var c *auth.AuthConfig = auth.NewAuthConfig()
	// set authorization mode
	c.Mode = protobase.AUTHModeStrict
	// set access group partitions
	c.AccessGroups = auth.AuthGroups{
		Members: map[string][][3]string{
			"User": [][3]string{
				{"can", "publish", "self/inbox"},
				{"can", "publish", "a/simple/demo"},
				{"can", "subscribe", "self/notifications"},
			},
			"Bot": [][3]string{
				{"can", "publish", "self/location"},
				{"can", "request", "self/access/upgrade"},
				{"can", "subscribe", "a/simple/demo"},
				{"can", "publish", "a/simple/demo"},
			},
			"Reader": [][3]string{
				{"can", "subscribe", "a/simple/demo"},
			},
		},
		// set access control list authorization
		// type.
		Type: protobase.ACLModeInclusive,
	}
	// define sample credentials
	c.Credentials = []auth.AuthEntity{

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test",
				Password: "$2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC",
				ClientId: "",
			},
			// define permission partition
			Group: "User",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test2",
				Password: "$2a$06$2uqusEvRMcpla2KXph8sBuBXO4WVOgIVbIgfRjk5y01UXxxgR9z6O",
				ClientId: "",
			},
			// define permission partition
			Group: "User",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test3",
				Password: "$2a$06$sgQ9yjjVvRxQhLqWKSGv4OTE2EF4ojUu1sEHnGUJdimmn.5M9M7/.",
				ClientId: "",
			},
			// define permission partition
			Group: "Bot",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test4",
				Password: "$2a$06$9wavlAtmNZ66Whe2wturDO7yIBdE41/Zcn4c5z4ydzJ/ydVJIZwJK",
				ClientId: "",
			},
			// define permission partition
			Group: "Bot",
		},

		auth.AuthEntity{
			Credential: &auth.Creds{
				Username: "test5",
				Password: "$2a$06$9wavlAtmNZ6623re2wturDO7yIBdEewfZcn4c5z4ydzJ/ydVJIZwJK",
				ClientId: "",
			},
			// define permission partition
			Group: "Reader",
		},
	}

	return c
}

var (
	c       *auth.AuthConfig        = defaultAuthConfig() // authentication configs
	authsys protobase.AuthInterface                       // authentication subsystem
)

func main() {
	var err error
	// setup authentication subsystem from configurations
	authsys, err = auth.NewAuthenticatorFromConfig(c)
	if err != nil {
		panic(err)
	}
	// initialize the broker with configured
	// authentication subsystem.
	brk := broker.NewBroker(broker.Options{
		Auth: authsys,
		ServerConf: server.ServerConfigs{
			Config: server.TLSOptions{
				Cert: "/Users/mitghi/go/src/github.com/mitghi/protox/config/cert/server.pem",
				Key:  "/Users/mitghi/go/src/github.com/mitghi/protox/config/cert/key.pem",
			},
			Mode: server.ProtoTLS,
			Addr: ":52909",
		},
		ClientDelegate: clientDelegate,
	}).(*broker.Broker)
	// run the broker
	ok := brk.Start()
	if !ok {
		fmt.Println("[-] unable to start.")
		os.Exit(1)
	}
	// wait for termination conditions such
	// as (KILL SIGNAL, FATAL ERRORS, .... ).
	<-brk.E

	os.Exit(0)
}
