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

package router

/**
* TODO
* . remove old subscribors from multi-level topics
* . move rest of the router from `Server`
* . remove maps
* . write unit tests
* . add retain storage
* . remove explicit buffpool dependency
* . REQUIRES FUNCTIONAL and WELL TESTED BUFFPOOL
**/

import (
	"fmt"

	"github.com/mitghi/protox/containers"
	"github.com/mitghi/protox/messages"
	"github.com/mitghi/protox/utils/strs"
)

func NewRouter() (r *Router) {
	r = &Router{
		subs:   NewSubs(10),
		retain: messages.NewRetain(),
	}
	return r
}

// // TODO
// func NewRouterWithBuffer(buff *buffpool.BuffPool) (r *Router) {
// 	r = &Router{
// 		subs:   NewSubsWithBuffer(10, buff),
// 		retain: messages.NewRetain(),
// 	}
// 	return r
// }

// - MARK: Router section.

func (r *Router) Add(client string, topic string, eticket byte) {
	r.Lock()
	r.subs.Lock()

	// drop cache lines
	r.subs.cache.RemoveCacheLines(topic)

	callback := func(node **containers.RDXNode, index int, endindex int) {
		n := (*node)
		if v := n.GetValue(); v != nil {
			m := v.(*subinfo)
			_ = m.insert(client, topic, eticket, index == endindex)
		} else {
			m := newSubinfo()
			_ = m.insert(client, topic, eticket, index == endindex)
			n.SetValue(m)
		}
		/* d e b u g */
		// r.subs.AddCacheLine(topic, ns)
		// logger.FTracef(1, "AddSub", "+ [Router] add subscription into the cahe line.", ns)
		/* d e b u g */
	}
	curr := r.subs.insertPath([]byte(topic), nil, Sep, callback)
	/* d e b u g */
	// var ns *subscription
	/* d e b u g */
	if curr != nil && (*curr) != nil {
		if v := (*curr).GetValue(); v == nil {
			nv := newSubinfo()
			_ = nv.insert(client, topic, eticket, true)
			(*curr).SetValue(nv)
		} else {
			/* d e b u g */
			// logger.FWarnf("AddSub", "* [Router] curr==nil for topic(%s) of client(%s).", topic, client)
			/* d e b u g */
			m := v.(*subinfo)
			mm, _ := m.get(client)
			mm.eticket = eticket
			_ = m.insert(client, topic, eticket, true)
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

func (r *Router) Find(topic string) (map[string]byte, error) {
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
				m[cls.uid] = cls.eticket
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
				if strs.Match(v.topic, topic, Sep, Wlcd) {
					ceticket, ok := m[k]
					if !ok {
						m[k] = v.eticket
					} else if v.eticket > ceticket {
						m[k] = v.eticket
					}
					r.subs.cache.AddCacheLines(topic, v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), Sep, callback)
	return m, err
}

func (r *Router) Remove(client string, topic string) error {
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
				if k == client && strs.Match(v.topic, topic, Sep, Wlcd) {
					/* d e b u g */
					fmt.Println("k==client")
					/* d e b u g */
					delete(s, k)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), Sep, callback)
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
				if strs.Match(v.topic, topic, Sep, Wlcd) {
					delete(s, k)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), Sep, callback)
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
				if strs.Match(v.topic, topic, Sep, Wlcd) {
					m[k] = append(m[k], v)
					r.subs.cache.AddCacheLines(topic, v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), Sep, callback)
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
				if strs.Match(v.topic, topic, Sep, Wlcd) {
					m[k] = append(m[k], v)
				}
			}
		}
	}
	err := r.subs.searchPath([]byte(topic), Sep, callback)
	return m, err
}
