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
	"sync"

	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/utils/strs"
)

// TODO
// . finish test suits
// . implement inclusive, exclusive auth system
// . finish documentation
// . add multi-level inheritence

var _ protobase.ACLInterface = (*ACL)(nil)

type ACL struct {
	*sync.RWMutex
	roles map[string]protobase.ACLPermInterface
}

type ACLNodeBase struct {
	Name  string
	nodes map[string]protobase.ACLNodeInterface
}

type Action struct {
	*ACLNodeBase
}

type Ability struct {
	*ACLNodeBase
}

type Perms struct {
	*ACLNodeBase
}

type Resource struct {
	*ACLNodeBase
	Providers []string
	role      *Role
}

type Entity struct {
	Name        string
	Descendants []string
}

type Role struct {
	*sync.RWMutex
	Name       string
	Mode       protobase.ACLMode
	entity     *Entity
	permission *Perms
}

type RoleUser struct {
	*Role
	Parent     protobase.ACLPermInterface
	Permission *Perms
	Mode       protobase.ACLMode
}

func NewACL() *ACL {
	return &ACL{RWMutex: &sync.RWMutex{}, roles: make(map[string]protobase.ACLPermInterface)}
}

func NewACLNodeBase() *ACLNodeBase {
	return &ACLNodeBase{nodes: make(map[string]protobase.ACLNodeInterface)}
}

func NewPerms(name string) *Perms {
	return &Perms{
		ACLNodeBase: &ACLNodeBase{
			Name:  name,
			nodes: make(map[string]protobase.ACLNodeInterface),
		},
	}
}

func NewRole(name string) *Role {
	return &Role{
		RWMutex:    &sync.RWMutex{},
		Name:       name,
		entity:     &Entity{Name: name},
		permission: NewPerms(name),
	}
}

func NewRoleUser(name string, role protobase.ACLPermInterface) protobase.ACLPermInterface {
	return &RoleUser{
		Role:       NewRole(name),
		Parent:     role,
		Permission: NewPerms(name),
		Mode:       protobase.ACLModeNormal,
	}
}

func (acl *ACL) getRole(name string) protobase.ACLPermInterface {
	if role, ok := acl.roles[name]; ok {
		return role
	}
	return nil
}

func (acl *ACL) makeRole(name string) (role protobase.ACLPermInterface, err error) {
	if role = acl.getRole(name); role == nil {
		role = NewRole(name)
		acl.roles[name] = role
		return role, nil
	}
	return nil, EACLInvalid
}

func (acl *ACL) GetRole(name string) protobase.ACLPermInterface {
	acl.RLock()
	defer acl.RUnlock()
	return acl.getRole(name)
}

func (acl *ACL) MakeRole(name string) (role protobase.ACLPermInterface, err error) {
	acl.Lock()
	defer acl.Unlock()
	return acl.makeRole(name)
}

func (acl *ACL) GetOrCreate(name string) (role protobase.ACLPermInterface, isNew bool) {
	acl.Lock()
	defer acl.Unlock()
	isNew = true
	role, err := acl.makeRole(name)
	if err != nil {
		isNew = false
		role = acl.getRole(name)
	}
	return role, isNew
}

func (acl *ACL) HasRole(name string) (hasRole bool) {
	acl.RLock()
	defer acl.RUnlock()
	_, hasRole = acl.roles[name]
	return hasRole
}

func (r *Role) SetPerm(ability string, action string, resource string) error {
	r.Lock()
	defer r.Unlock()
	var permission *Perms = r.permission
	l := []string{ability, action, resource}
	return permission.Add(l...)
}

func (r *Role) UnsetPerm(ability string, action string, resource string) error {
	r.Lock()
	defer r.Unlock()
	var permission *Perms = r.permission
	l := []string{ability, action, resource}
	_, err := permission.Unset(l...)
	return err
}

func (r *Role) HasPerm(ability string, action string, resource string) bool {
	r.Lock()
	defer r.Unlock()
	l := []string{ability, action, resource}
	return r.permission.CanDo(true, l...)
}

func (r *Role) HasExactPerm(ability string, action string, resource string) bool {
	r.Lock()
	defer r.Unlock()
	l := []string{ability, action, resource}
	return r.permission.CanDo(false, l...)
}

func (r *Role) SetMode(mode protobase.ACLMode) bool {
	// TODO
	// . return proper boolean response to indicate
	//   success status.
	r.Mode = mode
	return true
}

func (ru *RoleUser) HasPerm(ability string, action string, resource string) bool {
	ru.Lock()
	l := []string{ability, action, resource}
	lr := ru.permission.CanDo(true, l...)
	ru.Unlock()
	lr = lr || ru.Parent.HasPerm(ability, action, resource)
	if ru.Mode == protobase.ACLModeExclusive {
		lr = !lr
	}
	return lr
}

func (ru *RoleUser) HasExactPerm(ability string, action string, resource string) bool {
	ru.Lock()
	l := []string{ability, action, resource}
	lr := ru.permission.CanDo(false, l...)
	ru.Unlock()
	lr = lr || ru.Parent.HasExactPerm(ability, action, resource)
	if ru.Mode == protobase.ACLModeExclusive {
		lr = !lr
	}
	return lr
}

func (ru *RoleUser) SetMode(mode protobase.ACLMode) bool {
	// TODO
	// . return proper boolean response to indicate
	//   success status.
	ru.Mode = mode
	return true
}

func (anb *ACLNodeBase) SetValue(key string, value protobase.ACLNodeInterface) bool {
	if _, ok := anb.nodes[key]; !ok {
		anb.nodes[key] = value
		return true
	}
	return false
}

func (anb *ACLNodeBase) Len() int {
	return len(anb.nodes)
}

func (anb *ACLNodeBase) RemoveValue(key string) bool {
	if _, ok := anb.nodes[key]; ok {
		delete(anb.nodes, key)
		return true
	}
	return false
}

func (anb *ACLNodeBase) HasWildIdentifier(item string) (ok bool) {
	for k, _ := range anb.nodes {
		if strs.Match(k, item, protobase.Sep, protobase.Wlcd) {
			return true
		}
	}
	return false
}

func (anb *ACLNodeBase) CanDo(wildmatch bool, args ...string) bool {
	if len(args) == 0 {
		return false
	} else if len(args) == 1 {
		if !wildmatch {
			node := anb.GetIdentifier(args[0])
			if node == nil {
				return false
			}
			return node.IsResource(args[0])
		} else {
			return anb.HasWildIdentifier(args[0])
		}
	}
	ident := args[0]
	n := anb.GetIdentifier(ident)
	if n == nil {
		return false
	}
	return n.CanDo(wildmatch, args[1:]...)
}

func (anb *ACLNodeBase) GetIdentifier(ident string) protobase.ACLNodeInterface {
	if n, ok := anb.nodes[ident]; ok {
		return n
	}
	return nil
}

func (anb *ACLNodeBase) Add(args ...string) error {
	if len(args) == 0 {
		return nil
	}
	ident := args[0]
	n := anb.GetIdentifier(ident)
	if n == nil {
		n = anb.MakeChild(len(args), ident)
		anb.nodes[ident] = n
		if len(args) == 1 {
			n.SetValue(ident, nil)
		}
	} else {
		if len(args) == 1 {
			return EACLViolation
		}
	}

	return n.Add(args[1:]...)
}

func (anb *ACLNodeBase) Unset(args ...string) (bool, error) {
	if len(args) == 0 {
		return true, nil
	}
	ident := args[0]
	n := anb.GetIdentifier(ident)
	if n == nil {
		return false, EACLInvalid
	}
	prune, err := n.Unset(args[1:]...)
	if err != nil {
		return false, err
	}
	if prune {
		delete(anb.nodes, ident)
	}
	prune = false
	if len(anb.nodes) == 0 {
		prune = true
	}
	return prune, nil
}

func (anb *ACLNodeBase) HasIdentifier(ident string) bool {
	_, status := anb.nodes[ident]
	return status
}

func (anb *ACLNodeBase) MakeChild(level int, ident string) protobase.ACLNodeInterface {
	switch level {
	case 3:
		return &Ability{&ACLNodeBase{Name: ident, nodes: make(map[string]protobase.ACLNodeInterface)}}
	case 2:
		return &Action{&ACLNodeBase{Name: ident, nodes: make(map[string]protobase.ACLNodeInterface)}}
	case 1:
		return &Resource{ACLNodeBase: &ACLNodeBase{Name: ident, nodes: make(map[string]protobase.ACLNodeInterface)}}
	default:
		return nil
	}
}

func (anb *ACLNodeBase) IsResource(ident string) bool {
	return anb.Name == ident
}
