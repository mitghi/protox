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

	"github.com/mitghi/protox/auth"
	"github.com/mitghi/protox/broker"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/protobase"
)

var (
	brk         *broker.Broker
	credentials []*auth.Creds
)

func clientDelegate(uid string, pid string, cid string) protobase.ClientInterface {
	ret := User{client.Client{Username: uid, Password: pid, ClientId: cid}}
	return &ret
}

func main() {

	credentials = []*auth.Creds{
		{Username: "test", Password: "$2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC", ClientId: ""},
		{Username: "test2", Password: "$2a$06$2uqusEvRMcpla2KXph8sBuBXO4WVOgIVbIgfRjk5y01UXxxgR9z6O", ClientId: ""},
		{Username: "test3", Password: "$2a$06$sgQ9yjjVvRxQhLqWKSGv4OTE2EF4ojUu1sEHnGUJdimmn.5M9M7/.", ClientId: ""},
		{Username: "test4", Password: "$2a$06$9wavlAtmNZ66Whe2wturDO7yIBdE41/Zcn4c5z4ydzJ/ydVJIZwJK", ClientId: ""},
	}

	brk = broker.NewBroker(broker.Options{ClientDelegate: clientDelegate}).(*broker.Broker)
	brk.RegisterClients(credentials)

	ok := brk.Start()
	if !ok {
		os.Exit(1)
	}

	<-brk.E

	os.Exit(0)
}
