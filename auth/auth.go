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
	"fmt"
	"time"

	"github.com/mitghi/protox/protobase"
)

/**
* TODO
* . add auth modes
* . represent route permissions in serialized form
*   to reduce memory footprints.
**/

var (
	EAUTHUnknownMode error = fmt.Errorf(eFMT, "auth", "unable to set mode")
)

// NewAuthenticator allocates and initializes  a new `Authentication`
// instance and returns a pointer to it.
func NewAuthenticator() *Authentication {
	result := &Authentication{
		accounts:    make(map[string]*AuthInfo),
		permissions: NewACL(),
		mode:        protobase.AUTHModeNone,
	}
	return result
}

// NewAuthenticatorFromConfig allocate and initializes a new `Authentication`
// instance and config it according to `config` argument. It returns an error
// in case of unsuccessfull operation or invalid configuration.
func NewAuthenticatorFromConfig(config *AuthConfig) (a *Authentication, err error) {
	if ok, err := config.IsValid(); (!ok) || (err != nil) {
		return nil, err
	}
	var (
		acl protobase.ACLInterface
	)
	a = NewAuthenticator()
	acl = a.GetACL()
	a.mode = config.Mode
	// register groups
	for group, perms := range config.AccessGroups.Members {
		role, err := acl.MakeRole(group)
		if err != nil {
			return nil, err
		}
		if !role.SetMode(config.AccessGroups.Type) {
			return nil, EAUTHUnknownMode
		}
		for _, perm := range perms {
			if err := role.SetPerm(perm[0], perm[1], perm[2]); err != nil {
				return nil, err
			}
		}
	}
	// register clients
	for _, cred := range config.Credentials {
		if (cred.Credential == nil) || (!cred.Credential.IsValid()) {
			return nil, ECREDINVAL
		}
		ok, err := a.RegisterToGroup(cred.Group, cred.Credential)
		if !ok || err != nil {
			return nil, ECREDINVAL
		}
	}

	return a, nil
}

// NewAuthInfo allocates and initializes a new `AuthInfo` with the given
// credentials and returns a pointer to it.
func NewAuthInfo(creds protobase.CredentialsInterface) *AuthInfo {
	return &AuthInfo{creds: creds, userType: protobase.AuthUserNormal}
}

// - MARK: Authtentication section.

// canAuth evaluates given credential validity and returns its associates information
// in case of successfull evaluation. It returns an error with `reason` set to `false`
// when an error occures or invalid credentials is given.
func (self *Authentication) canAuth(creds protobase.CredentialsInterface) (user *AuthInfo, reason bool, err error) {
	uid := creds.GetUID()
	if user, ok := self.getUserWithIdentifier(&uid); ok {
		if user == nil {
			return nil, false, NonExistingUser
		}
		user.RLock()
		usrcreds := user.creds
		user.RUnlock()
		if creds.Match(usrcreds) {
			return user, true, nil
		}
		return nil, false, BadPassword
	}

	return nil, false, NonExistingUser
}

// CanAuthenticate returns a boolean indicating validity of the given credentials. It returns
// an error propogated from lower levels.
func (self *Authentication) CanAuthenticate(creds protobase.CredentialsInterface) (ok bool, err error) {
	user, ok, err := self.canAuth(creds)
	if !ok && user != nil {
		user.Lock()
		user.stat.faults += 1
		user.Unlock()
	}
	logger.Debugf("* [AuthSys] auth status for [uid] %s is %t .", creds.GetUID(), ok)
	return ok, err
}

func (self *Authentication) GetUserType(uid string) (utype protobase.AuthUserType, err error) {
	if uinfo, ok := self.getUserWithIdentifier(&uid); !ok {
		return "", fmt.Errorf(eFMT, "auth", "unable to find user with given id")
	} else {
		return uinfo.userType, nil
	}
}

// Register takes a `protobase.CredentialsInterface` struct and tries to register it.
// It returns true iff the given credential has not been registered prior to current
// attempt.
func (self *Authentication) Register(creds protobase.CredentialsInterface) (result bool) {
	uid := creds.GetUID()
	if !self.hasUserWithIdentifier(&uid) {
		authinfo := NewAuthInfo(creds)
		self.Lock()
		self.accounts[uid] = authinfo
		self.Unlock()
		return true
	}

	return false
}

// RegisterToGroup takes a `protobase.CredentialsInterface` struct and
// tries to register it. It returns true iff the given credential has
// not been registered prior to current attempt and iff given `group`
// exists.
func (self *Authentication) RegisterToGroup(group string, creds protobase.CredentialsInterface) (ok bool, err error) {
	uid := creds.GetUID()
	self.Lock()
	defer self.Unlock()
	if !self.permissions.HasRole(group) {
		return false, fmt.Errorf(eFMT, "auth", "attempt adding user to non-existing group.")
	}
	_, ok = self.accounts[uid]
	if !ok {
		authinfo := NewAuthInfo(creds)
		authinfo.userType = (protobase.AuthUserType)(group)
		self.accounts[uid] = authinfo
		return true, nil
	}

	return false, EAUTHUserReadd
}

// hasUserWithIdentifier takes a `string` pointer and returns a boolean indicating
// existence of current identifier.
func (self *Authentication) hasUserWithIdentifier(username *string) (ok bool) {
	if username == nil || *username == "" {
		return false
	}
	self.RLock()
	_, ok = self.accounts[(*username)]
	self.RUnlock()
	return ok
}

// HasClient returns a boolen to indicate whether a client
// with given identifier exists or not.
func (self *Authentication) HasClient(uid string) (ok bool) {
	self.RLock()
	defer self.RUnlock()
	_, ok = self.accounts[uid]
	return ok
}

// getUserWithIdentifier takes a `string` pointer and returns the associated information
// in case of existing identifier and a boolean indicating success status.
func (self *Authentication) getUserWithIdentifier(username *string) (user *AuthInfo, ok bool) {
	if username == nil || *username == "" {
		return nil, false
	}
	self.RLock()
	user, ok = self.accounts[(*username)]
	self.RUnlock()
	return user, ok
}

// RemoveWithIdentifier takes a `string` pointer and tries to remove the entry
// associated with the given identifier when it exists and indicate its success
// with a boolean. It also returns an error when unsuccessfull.
func (self *Authentication) RemoveWithIdentifier(identifier *string) (result bool, err error) {
	if identifier == nil || (identifier != nil && *identifier == "") {
		// TODO
		// . add proper error msg
		return false, nil
	}
	if self.hasUserWithIdentifier(identifier) {
		self.Lock()
		delete(self.accounts, *identifier)
		self.Unlock()
		return true, nil
	}
	return false, NonExistingUser
}

// MakeCreds takes standard `protobase.CredentialsInterface` arguments and
// creates a new `protobase.CredentialsInterface`.
func (self *Authentication) MakeCreds(uid string, pid string, cid string, args ...interface{}) (creds protobase.CredentialsInterface, err error) {
	// TODO
	// . make creds with internal providers, if available.
	// . user args to add additional meta infos.
	// . return error in case of invalid args.
	// . make it generic as much as possible.
	return &Creds{uid, pid, cid}, nil
}

/* TODO */
//
func (self *Authentication) Authenticate(creds protobase.CredentialsInterface) bool {
	// TODO
	// . unify authentication methods
	return false
}

//
func (self *Authentication) HasSession(clientId string) (result bool) {
	// TODO
	// . refactor session information into a
	//   single struct and return proper
	//   information from this method.
	return false
}

/* TODO */

// TryAuthenticate evaluates the given credentials and tries to authenticate with it.
// It returns a boolean indicating its success status.
func (self *Authentication) TryAuthenticate(creds protobase.CredentialsInterface) bool {
	user, status, _ := self.canAuth(creds)
	if !status {
		return status
	} else if user == nil {
		return false
	}
	t := time.Now()
	user.Lock()
	user.SetAuthorized()
	user.stat.succs += 1
	user.lstacc = &t
	user.Unlock()
	return status
}

// TryUnAuthenticate takes an identifier and tries to unauthenticate the entry
// associated with it. It returns a boolean indicating its success status.
func (self *Authentication) TryUnAuthenticate(uid string) bool {
	if user, ok := self.getUserWithIdentifier(&uid); ok && user != nil {
		user.Lock()
		user.UnsetAuthorized()
		user.Unlock()
		return ok
	}
	return false
}

// CreateGroup creates a new ACL group with the given permissions.
func (self *Authentication) CreateGroup(name string, permissions [][3]string) (err error) {
	// TODO
	// . check if duplicate permission can occur
	if len(permissions) == 0 {
		return EAUTHInvalidPerms
	}
	role, _ := self.permissions.GetOrCreate(name)
	if role == nil {
		return EAUTHGeneralFailure
	}
	for _, v := range permissions {
		err = role.SetPerm(v[0], v[1], v[2])
		if err != nil {
			return err
		}
	}

	return nil
}

// SetMode is a receiver method that sets the authorization
// mode.
func (self *Authentication) SetMode(mode protobase.AuthMode) {
	self.mode = mode
}

// GetMode is a getter for authentication mode.
func (self *Authentication) GetMode() protobase.AuthMode {
	return self.mode
}

// GetACL returns internal ACL subsystem. It is important to
// ensure returned value is not null ( in absence of ACL ).
func (self *Authentication) GetACL() protobase.ACLInterface {
	return self.permissions
}

// - MARK: AuthInfo section.

// IsAuthorized returns whether the current entry is authorized.
func (self *AuthInfo) IsAuthorized() byte {
	return self.status
}

// SetAuthorized sets the authorization status to true.
func (self *AuthInfo) SetAuthorized() {
	self.status = UAuthorized
}

// UnsetAuthorized unauthorizes the current entry by setting
// authorization flag to false.
func (self *AuthInfo) UnsetAuthorized() {
	self.status = UNotAuthorized
}

// SetType sets user type flag.
func (self *AuthInfo) SetType(t protobase.AuthUserType) {
	self.userType = t
}

// GetType returns associated `protobase.AuthUserType` of the
// given entry.
func (self *AuthInfo) GetType() protobase.AuthUserType {
	return self.userType
}
