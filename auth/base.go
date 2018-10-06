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
	"errors"
	"sync"
	"time"

	"github.com/mitghi/protox/logging"
	"github.com/mitghi/protox/protobase"
)

/**
* TODO:
* . check struct alignments
* . implement endpoints
* . add a mechanism to dump `Auth` into
*   3rd party services.
* . add end points to config
**/

// Generic error format
const (
	eFMT string = "%s: %s"
)

// Permission constants
const (
	// Permission requires exact length of following
	// constant.
	AuthACLPermLength = 3
)

// Authenication modes
const (
	// ModUSRPSWD is a flag to indicate the usage of username/password
	ModUSRPSWD = iota
	// ModSIG is a flag to indicate the usage of signature
	ModSIG
)

// Authorization flags
const (
	UNotAuthorized = iota
	UAuthorized
)

// - MARK: Error codes.

// ACL and Config error messages
var (
	EACLInvalid           error = errors.New("permissions: attempt to unset non existing node")
	EACLViolation         error = errors.New("permissions: attempt to readd resource")
	EACInconsistentConfig error = errors.New("auth(config): invalid configuration for mode")
	EACINVAL              error = errors.New("auth(config): invalid/unknown value/flag in configuration")
)

// Error messages
var (
	NonExistingUser     error = errors.New("permissions: user does not exist")
	BadPassword         error = errors.New("permissions: invalid password")
	EAUTHInvalidPerms   error = errors.New("permissions: invalid or insufficent permission list")
	EAUTHGeneralFailure error = errors.New("permissions: general operation failure")
	EAUTHUserReadd      error = errors.New("auth: attempt to re-registering existing user")
	ECREDINVAL          error = errors.New("credentials: missing or invalid credentials")
)

// Debug codes ( for development )
var (
	EAUTHNotImplemented error = errors.New("auth: not implemented")
)

// logger is the logging facility.
var logger protobase.LoggingInterface

func init() {
	logger = logging.NewLogger("Auth")
}

// stats is a struct that contains statistics associated
// with an entry.
type stats struct {
	succs  uint64
	faults uint64
}

// AuthInfo is a struct that is associated to each registered
// identifier in `Authentication`. It contains informations
// such as access times, statistics, ip address, permissions
// and etc .... .
type AuthInfo struct {
	sync.RWMutex
	creds    protobase.CredentialsInterface // credentials
	stat     stats                          // statistics
	lstacc   *time.Time                     // last access time
	lstdeacc *time.Time                     // last deauth time
	lstip    string                         // last ip addr
	perms    *RoleUser                      // client permissions
	acb      byte                           // access control bits
	status   byte                           // connection status
	userType protobase.AuthUserType         // user type
}

// Creds is a basic credential container.
type Creds struct {
	Username string
	Password string
	ClientId string
}

// Authentication is a `protobase.AuthInterface` compatible struct.
type Authentication struct {
	sync.RWMutex
	accounts    map[string]*AuthInfo
	permissions *ACL
	mode        protobase.AuthMode
}

// - MARK: Config structs group.

// AuthGroup is a struct used to define individual access
// setting. It is used for initial Auth subsystem configuration.
type AuthGroups struct {
	Members map[string][][3]string // Groups contains default feasible permissions
	Type    protobase.ACLMode      // Type is default Access type ( i.e. Inclusive, Exclusive or Undefined )
}

type AuthEntity struct {
	Credential protobase.CredentialsInterface // Credential contains entity cred. (e.g. id, passwd, .... )
	Group      string                         // Gruop specifies an association to certain Authorization group
}

// AuthConfig is a struct used to config Auth subsystem during/ and
// after initialization. It defines global access rules such as
// Authenication mode.
type AuthConfig struct {
	AccessGroups AuthGroups
	Credentials  []AuthEntity
	Mode         protobase.AuthMode
}
