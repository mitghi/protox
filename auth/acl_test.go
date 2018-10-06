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

func TestACLBasic(t *testing.T) {
	acl := NewACL()
	r, err := acl.MakeRole("guest")
	if err != nil {
		t.Fatal("cannot add role.")
	}
	_, err = acl.MakeRole("guest")
	if err == nil {
		t.Fatal("expected an error indicating write violation, got nil.")
	}
	r.SetPerm("can", "subscribe", "a/simple/demo")
	r.SetPerm("can", "publish", "a/simple/demo")
	r.SetPerm("can", "read", "a/simple/demo")
	r.SetPerm("cannot", "delete", "a/simple/demo")
	if !r.HasExactPerm("can", "subscribe", "a/simple/demo") {
		t.Fatal("invalid permission status.")
	}
	if r.HasExactPerm("non", "existing", "action and resource") {
		t.Fatal("invalid permission status, accessed non existing resource.")
	}
	if !r.HasExactPerm("cannot", "delete", "a/simple/demo") {
		t.Fatal("invalid permission status.")
	}
	if err := r.UnsetPerm("cannot", "delete", "a/simple/demo"); err != nil {
		t.Fatal(err)
	}
	if err := r.UnsetPerm("can", "read", "a/simple/demo"); err != nil {
		t.Fatal(err)
	}
	if err := r.UnsetPerm("can", "publish", "a/simple/demo"); err != nil {
		t.Fatal(err)
	}
	if err := r.UnsetPerm("can", "subscribe", "a/simple/demo"); err != nil {
		t.Fatal(err)
	}
	r.SetPerm("can", "publish", "a/simple/demo")
	if !r.HasExactPerm("can", "publish", "a/simple/demo") {
		t.Fatal("invalid permission status.")
	}
	if err := r.UnsetPerm("cannot", "delete", "a/simple/demo"); err == nil {
		t.Fatal("deleted nonexisting permission.")
	}
}

type ADUser struct {
	role *Role
}

func TestACLUserRole(t *testing.T) {
	acl := NewACL()
	r, err := acl.MakeRole("user")
	if err != nil {
		t.Fatal("unable to make a role.")
	}
	err = r.SetPerm("can", "read", "inbox")
	if err != nil {
		t.Fatal("unable to set permission.")
	}
	err = r.SetPerm("cannot", "delete", "inbox")
	if err != nil {
		t.Fatal("unable to set permission.")
	}
	ru := NewRoleUser("test", r)
	err = ru.SetPerm("can", "eat", "potato")
	if err != nil {
		t.Fatal("unable to set permission.")
	}
	if !ru.HasExactPerm("can", "eat", "potato") {
		t.Fatal("inconsistent state, expected true value for existing permission.")
	}
	if !ru.HasExactPerm("cannot", "delete", "inbox") {
		t.Fatal("inconsistent state, expected true value for existing permission.")
	}
	if ru.HasExactPerm("can", "read", "transactions") {
		t.Fatal("inconsistent state, expected false value for nonexisting permission.")
	}
	err = ru.UnsetPerm("can", "eat", "potato")
	if err != nil {
		t.Fatal("unable to unset existing permission.")
	}
	if ru.HasExactPerm("can", "eat", "potato") {
		t.Fatal("inconsistent state, expected false value for non existing value.")
	}
	err = ru.UnsetPerm("can", "read", "inbox")
	if err == nil {
		t.Fatal("inconsistent state, attempt to unset nonexisting permission ( permission does not exists in child node .).")
	}
	if !ru.HasExactPerm("can", "read", "inbox") {
		t.Fatal("inconsistent state, received false value for existing permission.")
	}
	// set exclusive mode
	ru.SetMode(protobase.ACLModeExclusive)
	if ru.HasExactPerm("can", "read", "inbox") {
		t.Fatal("inconsistent state, received true for existing permission in exclusive mode.")
	}
	ru.SetMode(protobase.ACLModeInclusive)
	if !ru.HasExactPerm("can", "read", "inbox") {
		t.Fatal("inconsistent state, received false for existing permission in inclusive mode.")
	}
	// add resource containing wildcard level
	err = ru.SetPerm("can", "publish", "a/generic/*/topic")
	if err != nil {
		t.Fatal("unable to set permission.")
	}
	// check wildcard matching
	if !ru.HasPerm("can", "publish", "a/generic/specific/topic") {
		t.Fatal("inconsistent state, received false for existing permission ( with wildcard ) in inclusive mode.")
	}
}
