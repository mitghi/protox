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

func TestAuth(t *testing.T) {
	var (
		roles     [][3]string
		testCases []struct {
			perm   [3]string
			expect bool
		}
		a   *Authentication = NewAuthenticator()
		acl *ACL            = a.GetACL().(*ACL)
	)
	// register following permissions for group "user"
	roles = [][3]string{
		{"can", "publish", "general/info"},
		{"can", "subscribe", "general/info"},
		{"can", "publish", "self/inbox"},
		{"can", "publish", "general/news"},
		{"can", "subscribe", "general/news"},
	}
	// assert these test cases
	testCases = []struct {
		perm   [3]string
		expect bool
	}{
		{perm: [3]string{"can", "subscribe", "general/news"}, expect: true},
		{perm: [3]string{"can", "subscribe", "self/inbox"}, expect: false},
	}
	// set authentication mode
	a.SetMode(protobase.AUTHModeStrict)
	// create a new role
	role, err := acl.MakeRole("user")
	if err != nil {
		t.Fatalf("expected err==nil, got %+v", err)
	}
	// add permissions
	for _, v := range roles {
		err = role.SetPerm(v[0], v[1], v[2])
		if err != nil {
			t.Fatalf("expected err==nil, got %+v", err)
		}
	}
	// check test assertions
	for _, v := range testCases {
		perm := v.perm
		if role.HasExactPerm(perm[0], perm[1], perm[2]) != v.expect {
			t.Fatalf("assertion failed, expected %t for permission ( %s %s %s ).", v.expect, perm[0], perm[1], perm[2])
		}
	}
}

func TestAuthRegisterToGroup(t *testing.T) {
	var (
		c   *AuthConfig = defaultAuthConfig()
		a   *Authentication
		err error
	)
	// create new auth subsystem from default configuration
	a, err = NewAuthenticatorFromConfig(c)
	if err != nil {
		t.Fatalf("assertion failed, expected err==nil, got %+v", err)
	}
	// register new user, this case must not fail
	ok, err := a.RegisterToGroup("User", &Creds{Username: "test3", Password: "test3", ClientId: ""})
	if err != nil {
		t.Fatalf("inconsistent state, expected err==nil, got %+v ( ok == %t ).", err, ok)
	}
	// re-register existing user, this case must fail
	ok, err = a.RegisterToGroup("User", &Creds{Username: "test3", Password: "test3", ClientId: ""})
	if err == nil {
		t.Fatalf("inconsistent state, expected err==nil, got %+v ( ok == %t ).", err, ok)
	}
}
