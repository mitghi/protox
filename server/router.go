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

package server

/**
* TODO
* . remove old subscribors from multi-level topics
* . move rest of the router from `Server`
* . remove maps
* . write unit tests
* . add retain storage
**/

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	buffpool "github.com/mitghi/lfpool"
	"github.com/mitghi/protox/client"
	"github.com/mitghi/protox/containers"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/protobase"
	"github.com/mitghi/protox/utils/strs"
)

var (
	// RRemoveErr indicates a problem associated wiht removing
	// a path from router.
	RRemoveErr error = errors.New("router: cannot remove")
	// SNoSubs indicates that topic does not exist.
	SNoSubs error = errors.New("subs: cannot locate subscriptions.")
)

type Router struct {
	sync.RWMutex

	routes map[string]map[net.Conn]*client.Client
	subs   *Subs
	retain *messages.Retain
}

type Subs struct {
	sync.RWMutex
	*containers.Radix
	*buffpool.BuffPool
	cache subcache
}

type subscription struct {
	topic  string
	uid    string
	qos    byte
	isLeaf bool
}

type subinfo struct {
	subs map[string]*subscription
}

// - MARK: Initializers.

func NewSubs(ralloc int) *Subs {
	result := &Subs{
		cache:    subcache{make(map[string]*subcacheline), 0, 0, 0, 0},
		Radix:    containers.NewRadix(ralloc),
		BuffPool: buffpool.NewBuffPool(),
	}
	return result
}

func NewSubsWithBuffer(ralloc int, buff *buffpool.BuffPool) *Subs {
	result := &Subs{
		cache:    subcache{make(map[string]*subcacheline), 0, 0, 0, 0},
		Radix:    containers.NewRadix(ralloc),
		BuffPool: buff,
	}
	return result
}

func NewRouter() *Router {
	newRouter := &Router{
		routes: make(map[string]map[net.Conn]*client.Client),
		subs:   NewSubs(10),
		retain: messages.NewRetain(),
	}

	return newRouter
}

func NewRouterWithBuffer(buff *buffpool.BuffPool) *Router {
	newRouter := &Router{
		routes: make(map[string]map[net.Conn]*client.Client),
		subs:   NewSubsWithBuffer(10, buff),
		retain: messages.NewRetain(),
	}

	return newRouter
}

func newSubinfo() *subinfo {
	return &subinfo{make(map[string]*subscription)}
}

// - MARK: subinfo section.

func (s *subinfo) insert(uid string, topic string, qos byte, isLeaf bool) (ret *subscription) {
	sb, ok := s.subs[uid]
	if !ok {
		ret = &subscription{topic, uid, qos, isLeaf}
		s.subs[uid] = ret
		return ret
	}
	sb.uid = uid
	sb.qos = qos
	sb.isLeaf = isLeaf
	return sb
}

func (s *subinfo) remove(uid string) bool {
	sb, ok := s.subs[uid]
	if !ok {
		return false
	}
	if sb == nil {
		return false
	}
	delete(s.subs, uid)
	return true
}

func (s *subinfo) get(uid string) (*subscription, bool) {
	sb, ok := s.subs[uid]
	if !ok {
		return nil, false
	}
	return sb, true
}

// - MARK: Subs section.

func (s *Subs) searchPath(path []byte, sep string, callback func(**containers.RDXNode, [][]byte, int, int)) error {
	var (
		buff  *buffpool.Buffer = s.GetBuffer(len(path))
		paths [][]byte
		plen  int
		curr  **containers.RDXNode
	)
	defer s.ReleaseBuffer(buff)
	paths, err := messages.TopicComponents(path)
	if err != nil {
		return errors.New("router: invalid topic.")
	}
	plen = len(paths) - 1
ML:
	for i, v := range paths {
		if i == 0 {
			var bstr string
			buff.Write(v)
			bstr = buff.String()
			_, tail, word := s.Find(bstr)
			if tail != nil && word == bstr {
				curr = &tail
			} else {
				return SNoSubs
			}
			buff.WriteByte('/')
			bstr = buff.String()
			_, ntail, nword := s.FindFrom(*curr, bstr)
			if ntail == nil {
				return SNoSubs
			}
			if n := (*ntail.Link); n != nil && nword == bstr {
				curr = &n
			} else {
				return SNoSubs
			}
			buff.Reset()
			continue
		}
		if *curr == nil {
			break
		}
		if string(v) != protobase.IdentWLCD {
			if vn := s.wldCheck(curr); vn != nil && callback != nil {
				callback(&vn, paths, i, plen)
			}
		}
		buff.WriteByte('/')
		buff.Write(v)
		bstr := buff.String()
		_, tail, word := s.FindFrom(*curr, bstr)
		if tail == nil {
			return SNoSubs
		}
		n := (*tail.Link)
		if n == nil {
			return SNoSubs
		}
		if word != bstr {
			return SNoSubs
		}
		curr = &n
		buff.Reset()
		if i == plen {
			if callback != nil {
				callback(curr, paths, i, plen)
				break ML
			}
		} else {
			buff.WriteString((*curr).Key)
			buff.WriteByte('/')
			_, tail, word = s.FindFrom(*curr, buff.String())
			if tail == nil {
				return SNoSubs
			}
			n = (*tail.Link)
			if n == nil {
				return SNoSubs
			}
			curr = &n
		}
		buff.Reset()
	}

	return nil
}

func (s *Subs) wldCheck(node **containers.RDXNode) *containers.RDXNode {
	_, tail, word := s.FindFrom((*node), protobase.TOPICWLCD)
	if tail != nil {
		if n := (*tail.Link); n != nil && word == protobase.TOPICWLCD {
			return n
		}
	}
	return nil
}

func (s *Subs) insertPath(path []byte, value interface{}, sep string, callback func(**containers.RDXNode, int, int)) **containers.RDXNode {
	var (
		buff  *buffpool.Buffer = s.GetBuffer(len(path))
		paths [][]byte
		plen  int
		curr  **containers.RDXNode
	)
	defer s.ReleaseBuffer(buff)
	// TODO
	//  check the error code
	paths, _ = messages.TopicComponents(path)
	plen = len(paths) - 1
	for i, v := range paths {
		if i == 0 {
			buff.Write(v)
			a, _ := s.InsertRFrom(s.GetRoot(), buff.String(), value)
			curr = &a
			a.NoMrg = true
			buff.Reset()
			buff.WriteString(a.Key)
			buff.WriteByte('/')
			b, _ := s.InsertRFrom(*curr, buff.String(), value)
			curr = &b
			b.NoMrg = true
			buff.Reset()
			continue
		}
		buff.WriteByte('/')
		buff.Write(v)
		a, _ := s.InsertRFrom(*curr, buff.String(), value)
		if a != nil {
			curr = &a
			a.NoMrg = true
			buff.Reset()
			buff.WriteString(a.Key)
		} else {
			logger.Fatal("(InsertPath) FATAl. TODO: FATAL SITUATION, CHECK.", "a", a)
			/* TODO */
			// NOTE: this event is fatal and should not occur
			// uncomment to crash the program in case of an error
			// and get stack trace dump.
			// panic("Subs: inconsistent state")
		}
		if buff.Len() > 0 {
			lch := buff.Bytes()[buff.Len()-1]
			/* d e b u g */
			// logger.Debug("len(buff)>0", lch, buff.String())
			// logger.Debug("insertPath()", "lch", lch)
			/* d e b u g */
			// TODO
			// . remove this later
			if strings.ContainsAny(string(lch), sep) == true {
			}
			n := *curr
			/* d e b u g */
			// NOTE: CHANGED:
			// if n != nil && n.Key == IdentWLCD {
			/* d e b u g */
			if n != nil {
				if n.Key == protobase.IdentWLCD {
					if callback != nil {
						callback(curr, i, plen)
					}
				} else if i == plen {
					/* d e b u g */
					// logger.FInfof("insertPath", "** [Router] in i==plen condition, plen(%d).", plen)
					/* d e b u g */
					if callback != nil {
						callback(curr, i, plen)
					}
				}
				// TODO: NOTE:
				//  this is for using last 1 byte in struct
				//  as option flags.
				/* d e b u g */
				// else {
				// 	n.SetOpts(0x1)
				// }
				/* d e b u g */
			}
		}
		if i != plen {
			buff.WriteByte('/')
			b, _ := s.InsertRFrom(*curr, buff.String(), value)
			curr = &b
			b.NoMrg = true
		}
		buff.Reset()
	}
	return curr
}

// - MARK: Router section.

func (r *Router) AddSub(client string, topic string, qos byte) {
	r.Lock()
	r.subs.Lock()

	// drop cache lines
	r.subs.cache.RemoveCacheLines(topic)

	callback := func(node **containers.RDXNode, index int, endindex int) {
		n := (*node)
		if v := n.GetValue(); v != nil {
			m := v.(*subinfo)
			_ = m.insert(client, topic, qos, index == endindex)
		} else {
			m := newSubinfo()
			_ = m.insert(client, topic, qos, index == endindex)
			n.SetValue(m)
		}
		/* d e b u g */
		// r.subs.AddCacheLine(topic, ns)
		// logger.FTracef(1, "AddSub", "+ [Router] add subscription into the cahe line.", ns)
		/* d e b u g */
	}
	curr := r.subs.insertPath([]byte(topic), nil, protobase.Sep, callback)
	/* d e b u g */
	// var ns *subscription
	/* d e b u g */
	if curr != nil && (*curr) != nil {
		if v := (*curr).GetValue(); v == nil {
			nv := newSubinfo()
			_ = nv.insert(client, topic, qos, true)
			(*curr).SetValue(nv)
		} else {
			/* d e b u g */
			// logger.FWarnf("AddSub", "* [Router] curr==nil for topic(%s) of client(%s).", topic, client)
			/* d e b u g */
			m := v.(*subinfo)
			mm, _ := m.get(client)
			mm.qos = qos
			_ = m.insert(client, topic, qos, true)
		}
	}
	/* d e b u g */
	// NOTE:
	// . this is wrong, cache lines should not be mutated in this routine
	//   at all. It must be repopulated by FindSub.
	// if ns != nil {
	// r.subs.cache.AddCacheLines(topic, ns)
	// logger.FTracef(1, "AddSub", "+ [Router] add cache line at the end of addsub routine.", ns)
	// }
	// NOTE:
	// print topic hierarchy from tree
	// r.subs.PrintV()
	/* d e b u g */

	r.subs.Unlock()
	r.Unlock()
}

func (r *Router) FindSub(topic string) (map[string]byte, error) {
	var (
		m map[string]byte = make(map[string]byte)
	)

	r.RLock()
	r.subs.Lock()
	defer r.subs.Unlock()
	defer r.RUnlock()

	if cline := r.subs.cache.GetCacheLines(topic); len(cline) != 0 {
		for _, v := range cline {
			for _, cls := range v {
				m[cls.uid] = cls.qos
			}
		}
		if len(cline) != 0 {
			/* d e b u g */
			// fmt.Println("returning from cline", len(cline), cline, m)
			/* d e b u g */
			return m, nil
		}
	}

	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
		n := (*node)
		if n.Value != nil {
			s := n.Value.(*subinfo).subs
			/* d e b u g */
			// fmt.Println("* [SUBS] after cache:", s)
			/* d e b u g */
			for k, v := range s {
				if strs.Match(v.topic, topic, protobase.Sep, protobase.Wlcd) {
					cqos, ok := m[k]
					if !ok {
						m[k] = v.qos
					} else if v.qos > cqos {
						m[k] = v.qos
					}
					r.subs.cache.AddCacheLines(topic, v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), protobase.Sep, callback)
	return m, err
}

func (r *Router) RemoveSub(client string, topic string) error {
	r.Lock()
	r.subs.Lock()
	defer r.subs.Unlock()
	defer r.Unlock()

	r.subs.cache.RemoveCacheLines(topic)

	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
		n := (*node)
		if n.Value != nil {
			s := n.Value.(*subinfo).subs
			for k, v := range s {
				/* d e b u g */
				fmt.Println("inside removesub", client, topic, k, v)
				/* d e b u g */
				if k == client && strs.Match(v.topic, topic, protobase.Sep, protobase.Wlcd) {
					/* d e b u g */
					fmt.Println("k==client")
					/* d e b u g */
					delete(s, k)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), protobase.Sep, callback)
	return err
}

func (r *Router) PruneSub(topic string) error {
	r.Lock()
	r.subs.Lock()
	defer r.subs.Unlock()
	defer r.Unlock()

	r.subs.cache.RemoveCacheLines(topic)

	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
		n := (*node)
		if n.Value != nil {
			s := n.Value.(*subinfo).subs
			for k, v := range s {
				if strs.Match(v.topic, topic, protobase.Sep, protobase.Wlcd) {
					delete(s, k)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), protobase.Sep, callback)
	return err
}

func (r *Router) FindSubC(topic string) (map[string][]*subscription, error) {
	var (
		m map[string][]*subscription = make(map[string][]*subscription)
	)

	r.RLock()
	r.subs.Lock()
	defer r.subs.Unlock()
	defer r.RUnlock()

	if cline := r.subs.cache.GetCacheLines(topic); len(cline) != 0 {
		for k, v := range cline {
			for _, cls := range v {
				/* d e b u g */
				fmt.Println("loop:", v, cls, k)
				/* d e b u g */
				m[cls.uid] = append(m[cls.uid], cls)
			}
		}
		if len(cline) != 0 {
			/* d e b u g */
			fmt.Println("returning from cline", len(cline), cline, m)
			/* d e b u g */
			return m, nil
		}
	}
	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
		n := (*node)
		if n.Value != nil {
			s := n.Value.(*subinfo).subs
			for k, v := range s {
				if strs.Match(v.topic, topic, protobase.Sep, protobase.Wlcd) {
					m[k] = append(m[k], v)
					r.subs.cache.AddCacheLines(topic, v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), protobase.Sep, callback)
	return m, err
}

func (r *Router) FindRawSubscribers(topic string) (map[string][]*subscription, error) {
	r.RLock()
	r.subs.Lock()
	defer r.subs.Unlock()
	defer r.RUnlock()

	m := make(map[string][]*subscription)
	callback := func(node **containers.RDXNode, paths [][]byte, level int, plen int) {
		n := (*node)
		if n.Value != nil {
			s := n.Value.(*subinfo).subs
			for k, v := range s {
				if strs.Match(v.topic, topic, protobase.Sep, protobase.Wlcd) {
					m[k] = append(m[k], v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), protobase.Sep, callback)
	return m, err
}
