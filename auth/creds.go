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
	"strings"

	"github.com/mitghi/protox/protobase"
)

// check interface ( protocol ) conformations
var _ protobase.CredentialsInterface = (*Creds)(nil)

// GetCredentials returns data associated with
// authenication method.
func (self *Creds) GetCredentials() (username string, password string, clientId string) {
	username = self.Username
	password = self.Password
	clientId = self.ClientId

	return
}

// GetUID returns a string used for
// user identification ( i.e. used id ).
func (self *Creds) GetUID() string {
	return self.Username
}

// Copy returns a new instance of a compatible
// `protobase.CredentialsInterface`.
func (self *Creds) Copy() protobase.CredentialsInterface {
	return &Creds{self.Username, self.Password, self.ClientId}
}

// cleanInput sanitizes input.
func (self *Creds) cleanInput(cred protobase.CredentialsInterface) bool {
	// TODO:
	// . this is a critical method and should
	//   be reimplemented in a sane way. This
	//   version is only a dummy.
	if (self.Username == "") || (self.Password == "") {
		return false
	}
	self.Username = strings.TrimSpace(self.Username)
	self.Password = strings.TrimSpace(self.Password)
	self.ClientId = strings.TrimSpace(self.ClientId)

	return true
}

//  Match is a receiver method that compares two
// `protobase.CredentialsInterface` and returns
// a boolean to indicate whether both are identical
// or not. It is used to match stored credentials
// against user-given credentials usually during
// initial handshake and initialization stage.
func (self *Creds) Match(cred protobase.CredentialsInterface) (ret bool) {
	if cred == nil {
		return false
	}
	var (
		uidok  bool
		pswok  bool
		clidok bool
	)
	switch cred.(type) {
	case *Creds:
		nc, _ := cred.(*Creds)
		uid, passwd, clid := nc.GetCredentials()
		uidok = self.Username == uid
		pswok = self.Password == passwd
		clidok = self.ClientId == clid
		ret = (uidok && pswok) && clidok
		break
	default:
		uid, passwd, clid := cred.GetCredentials()
		uidok = self.Username == uid
		pswok = self.Password == passwd
		clidok = self.ClientId == clid
		ret = (uidok && pswok) && clidok
		break
	}

	return ret
}

// IsValid returns a boolean indicating that
// whether the actual credentials are properly
// formatted and checks edge cases ( e.g. empty
// strings ).
func (self *Creds) IsValid() (ok bool) {
	// TODO
	// . implement format checks
	// . add edge cases
	ok = ((len(self.Username) > 0) &&
		(len(self.Password) > 0))
	ok = ok && ((len(strings.TrimSpace(self.Username)) > 0) &&
		(len(strings.TrimSpace(self.Password)) > 0))

	return ok
}
