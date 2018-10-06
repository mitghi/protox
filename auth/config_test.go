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

package auth

import (
	"testing"

	"github.com/mitghi/protox/protobase"
)

func defaultAuthConfig() *AuthConfig {
	var (
		c *AuthConfig = NewAuthConfig()
	)
	// start definition procedure.
	// STEP 1:
	// . define Authorization Gruops
	// . define their Access Control Mode
	c.AccessGroups = AuthGroups{
		Members: map[string][][3]string{
			"User": [][3]string{
				{"can", "publish", "self/inbox"},
				{"can", "subscribe", "self/notifications"},
			},
		},
		Type: protobase.ACLModeInclusive,
	}
	// STEP 2:
	// . add default credentials ( e.g. admins, managers,
	//   backup bots, .... ).
	c.Credentials = []AuthEntity{
		AuthEntity{
			Credential: &Creds{Username: "test", Password: "test", ClientId: "test"},
			Group:      "User",
		},
		AuthEntity{
			Credential: &Creds{Username: "test2", Password: "test2", ClientId: "test2"},
			Group:      "User",
		},
	}
	// STEP 3:
	// . associate Authentication Mode ( global )
	//   for initialized auth subsystem.
	// NOTE:
	// . AUTHModeStrict -> Authoentication Mode is
	//   `Strictly` defined, meaning that all dynamic
	//    modification requests during runtime are
	//    refused and rejected.
	c.Mode = protobase.AUTHModeStrict

	return c
}

func TestAuthConfig(t *testing.T) {
	var (
		c *AuthConfig = defaultAuthConfig()
	)
	if ok, err := c.IsValid(); !ok || err != nil {
		t.Fatalf("expected ok==true, err==nil. But got ok==%t, err==%+v.", ok, err)
	}
}
