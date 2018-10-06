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

	"github.com/mitghi/protox/protobase"
)

const (
	cACName string = "auth(config)"
)

/**
*
* TODO:
* . implement router mode
* . refactor boilerplate with interface
* . add easy mechanism to get preconfiged
*   auth subsystem ( factory ).
* . test suite for edge cases
*
**/

// NewAuthConfig is a function that allocate and
// initializes `AuthConfig` and returns a pointer
// to it. The default authentication mode is
// `protobase.AUTHModeNone` which returns an error
// during validity checks intentionally to prevent
// complications during development & debugging.
func NewAuthConfig() *AuthConfig {
	// NOTE:
	// . available authorization modes:
	//   check `protobase.constants` for
	//   changes.
	// ..
	// .. AUTHModeNone
	// .. AUTHModeDynamic
	// .. AUTHModeStrict
	// .. AUTHModeRouter
	var ac *AuthConfig = &AuthConfig{
		Mode: protobase.AUTHModeNone,
	}
	return ac
}

// - MARK: Auth configuration.

// IsValid is a receiver method that checks validity
// of underlaying data and returns an error in case of
// unsuccessfull operation. It can be used manually to
// ensure integrity, but it is mainly used by Auth
// subsystem.
func (ac *AuthConfig) IsValid() (ok bool, err error) {
	// TODO
	// . implement rest of authorization
	//   modes.
	switch ac.Mode {
	case protobase.AUTHModeNone:
		err = EACInconsistentConfig
		goto ERROR
	case protobase.AUTHModeDynamic:
		ok = true
		goto OK
	case protobase.AUTHModeStrict:
		ok = ((ac.AccessGroups.Len() > 0) && len(ac.Credentials) > 0)
		agok, agerr := ac.AccessGroups.IsValid()
		if agerr != nil {
			err = agerr
			goto ERROR
		}
		ok = ok && agok
		if !ok {
			err = fmt.Errorf(eFMT, cACName, "cannot validate Access Groups")
			goto ERROR
		}
		lok, lerr := ac.hasValidCreds()
		if !lok || lerr != nil {
			err = lerr
			goto ERROR
		}
		ok = ok && lok
		goto OK
	case protobase.AUTHModeRouter:
		// NOTE:
		// . temporarily return an error to indicate
		//   that this mode is still not implemented
		//   ( for development ).
		err = EAUTHNotImplemented
		goto ERROR
	default:
		err = EACINVAL
		goto ERROR
	}

OK:
	return ok, nil
ERROR:
	return false, err
}

// hasValidCreds checks validity of implementors
// of `protobase.CredentialsInterface` and returns
// an error in case of invalid entry.
func (ac *AuthConfig) hasValidCreds() (ok bool, err error) {
	for k, v := range ac.Credentials {
		if !v.Credential.IsValid() {
			err = fmt.Errorf(eFMT, cACName,
				fmt.Sprintf("(%d)th contains invalid entry ( %+v )", (k+1), v),
			)
			return false, err
		}
	}

	return true, nil
}

// SetMode sets Authorization mode globally for Auth
// subsystem. It returns false in case `mode` argument
// is invalid.
func (ac *AuthConfig) SetMode(mode protobase.AuthMode) (ok bool) {
	ok = (mode == protobase.AUTHModeNone) ||
		(mode == protobase.AUTHModeDynamic) ||
		(mode == protobase.AUTHModeStrict) ||
		(mode == protobase.AUTHModeRouter)
	if !ok {
		return false
	}
	ac.Mode = mode
	return ok
}

// AddCredentials is a receiver method which adds a new
// entry to its storage. The `group` argument is used to
// associate a given entry to the corresponding Auth Group
// in `AccessGroups`. It returns an error in case of unsucc-
// sessfull operation.
func (ac *AuthConfig) AddCredential(group string, cred protobase.CredentialsInterface) (err error) {
	// early return to avoid unneccessary stack or heap allocations.
	if cred == nil {
		return fmt.Errorf(eFMT, "auth(config)", "argument is a null pointer")
	} else if !cred.IsValid() {
		return EACINVAL
	}
	lac := *ac
	lac.Credentials = append(lac.Credentials, AuthEntity{Credential: cred, Group: group})
	*ac = lac

	return nil
}

// - MARK: Authorization Groups ( Partitions ).

// IsValid checks validity of underlaying data and
// returns an error in case of violation. It is
// used by `AuthConfig` and invoked prior to `AuthConfig`'s
// own validation procedure.
func (ag *AuthGroups) IsValid() (ok bool, err error) {
	// NOTE:
	// . available Access Control modes:
	//   check `protobase.constants` for
	//   changes.
	// ..
	// .. ACLModeNormal
	// .. ACLModeInclusive
	// .. ACLModeExclusive
	ok = (ag.Type == protobase.ACLModeNormal) ||
		(ag.Type == protobase.ACLModeInclusive) ||
		(ag.Type == protobase.ACLModeExclusive)
	if !ok {
		err = EACInconsistentConfig
		goto ERROR
	}

	return ok, nil
ERROR:
	return false, err
}

// Add is a receiver method that creates a new
// group when neccessary and adds the given
// permission line to it. It returns an error
// to indicate conformance violation.
func (ag *AuthGroups) Add(name string, perm ...string) (err error) {
	lperm := len(perm)
	// early return to avoid unneccessary stack and heap allocations.
	if (lperm > AuthACLPermLength) || (lperm < AuthACLPermLength) {
		return fmt.Errorf(eFMT, "auth(config)", "Add(....) requires exactly "+string(AuthACLPermLength))
	}
	var (
		entity  string = perm[0]
		action  string = perm[1]
		ability string = perm[2]
	)
	if _, ok := ag.Members[name]; !ok {
		ag.Members[name] = make([][3]string, 8)
	}
	ag.Members[name] = append(ag.Members[name], [3]string{
		entity,
		action,
		ability,
	})

	return nil
}

// HasGroup returns whether a given Auth Group is
// registered.
func (ag *AuthGroups) HasGroup(name string) (ok bool) {
	if _, ok = ag.Members[name]; ok {
		return true
	}
	return false
}

// Len returns number of total registered groups.
func (ag *AuthGroups) Len() int {
	return len(ag.Members)
}
