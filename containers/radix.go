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

package containers

/*
*
* This package offers efficient data structure with low memory footprints.
* All memory allocated by these structures are reused and recycled internally via free list.
* It is the building block for other components such as Router, Message Storage and
* Subscribe Management.
*
 */

import (
	"bytes"
	"fmt"
	"time"
)

// Tick is used by the sliding window algorithm.
// It is the delay between subsequents executions
// in Milliseconds. default value : 100ms
const tick = 100

// Memory Pool flags
const (
	RDXFoccupied byte = iota
	RDXFfree
)

// RDXNode containers relevant information for each leaf and node.
type RDXNode struct {
	// size: 32 bytes
	//
	// link is head of childs
	// next is next node in chain
	// Isbnd is a boolean to indicate that this node is end of a word
	// value points to a generic data structure ( like void* )
	Key   string      //8
	Value interface{} //8
	Link  *RDXNode    //4
	Next  *RDXNode    //4
	size  int         //4
	h     byte        //1
	Isbnd bool        //1
	NoMrg bool        //1
	Opts  byte        //1
}

// RDXResult is a struct that chains a branch of tree
// which contains the result. `Link` is a double pointer
// pointing to the actual node in the tree, `Next` points
// to next element in the chain.
//   ----     ----     ----     ----
//  |Link|   |Link|   |Link|   |Link|
// -|Next|.. |Next|.. |Next|.. |Next|->....
//   ----     ----     ----     ----
//    ^        ^        ^        ^
//    |        |        |        |
//    |        |        |        |
//   ----     ----     ----     ----
//  |Link|   |Link|   |Link|   |Link|
// -|Next|-->|Next|-->|Next|-->|Next|-->....
//   ----     ----     ----     ----
// Important: `Link` and `Next` do not have the same semantic meaning
// as `RDXNode` struct.
type RDXResult struct {
	// size: 8 bytes
	//
	Link **RDXNode
	Next *RDXResult
}

// NewRDXResult creates a pointer to a new `RDXResult` and initializes
// its `Link` double pointer.
func NewRDXResult() *RDXResult {
	var result *RDXResult = &RDXResult{}
	result.Link = new(*RDXNode)

	return result
}

// Radix contains the root node, serving as entry and a pool which is
// for reusing memory when a node is removed.
type Radix struct {
	// size: 12 bytes
	//
	root *RDXNode
	pool *Pool

	onRemove func(*RDXNode)
}

// GetRoot returns the root node of the current radix.
func (self *Radix) GetRoot() *RDXNode {
	return self.root
}

// TODO
//  *check padding*
//  implement resize
//  implement rolling window
//  optimize chunks based on statics
//  reuse deleted nodes
type Pool struct {
	// size: 64 bytes
	//
	// NOTE: numbers are for debug purposes.
	// It indicates amount of memory per item.
	// For performance it is very important
	// to word align all the items so that
	// `padded` bytes inbetween items reduces
	// to the most minimized value possible.
	current  time.Time  //16
	chunks   [3]int     //12
	nodes    []*RDXNode //12
	index    int        //4
	capacity int        //4
	ralloc   int        //4
	counter  int        //4
	tindex   int        //4
}

//
func (self *Pool) checkStat() bool {
	duration := time.Since(self.current)
	diff := int(duration.Seconds() * 1000)
	// NOTE:
	//  checks every `tick`ms
	if diff >= tick {
		self.current = time.Now()
		// TODO
		//  adapt horizon size dynamically
		//
		// current sliding window has a horizon of length '3'
		self.chunks[0] = self.chunks[1]
		self.chunks[1] = self.chunks[2]
		self.chunks[2] = self.counter
		self.counter = 0
		if self.tindex < 2 {
			self.tindex += 1
		}
		return true
	}
	return false
}

//
func (self *Pool) Insert() {
	self.checkStat()
	self.counter += 1
}

//
func (self *Pool) Counter() int {
	self.checkStat()
	return self.counter
}

//
func (self *Pool) InsertAStat() int {
	self.checkStat()
	self.counter += 1
	return self.chunks[self.tindex]
}

//
func (self *Pool) Chunks() int {
	self.checkStat()
	var (
		start = self.tindex - 3
		end   = self.tindex
	)
	if start < 0 {
		if self.tindex > 0 {
			return self.chunks[end]
		}
	}
	if self.tindex == 0 {
		return 0
	}
	return self.chunks[end]
}

// RPrint recurses the tree and prints the information to
// stdout.
func rprint(node *RDXNode, l string, val bool) {
	if node != nil {
		if node.Isbnd {
			switch val {
			case true:
				fmt.Printf("|\x1b[36m%-48s\033[0m %-s|%-4v|(%s)\n", l+node.Key, "+", node.Value, string(node.h))
			case false:
				fallthrough
			default:
				fmt.Printf("|\x1b[36m%-48s\033[0m %-s|\n", l+node.Key, "+")
			}
		} else {
			switch val {
			case true:
				fmt.Printf("|%-48s %-s|%-4v|\n", l+node.Key, " ", node.Value)
			case false:
				fallthrough
			default:
				fmt.Printf("|%-48s %-s|\n", l+node.Key, " ")
			}
		}
		if node.Link != nil {
			rprint(node.Link, l+"..", val)
		}
		rprint(node.Next, l, val)
	}
}

// NewRadix returns a pointer to a new `Radix` struct.
func NewRadix(ralloc int) *Radix {
	r := new(Radix)
	node := new(RDXNode)
	node.Key = ""
	node.Value = nil
	// add memory pool
	r.pool = NewPool(ralloc)
	r.root = node
	r.onRemove = nil
	return r
}

// NewPool creates a memory pool and recycle nodes.
// NOTE: this is note ready to use.
func NewPool(ralloc int) *Pool {
	return &Pool{
		time.Now(),
		[3]int{0, 0, 0},
		make([]*RDXNode, ralloc),
		-1,
		ralloc,
		ralloc,
		0,
		0,
	}
}

// Print recursively prints the internals.
func (self *Radix) Print() {
	rprint(self.root, "", false)
}

//
func (self *Radix) PrintV() {
	rprint(self.root, "", true)
}

// NewNode returns a pointer to a `RDXNode` from the pool.
func (self *Radix) newNode() *RDXNode {
	// TODO
	//  integrate event counter
	//  dynamically adapt allocations
	if self.pool.index == -1 || self.pool.index == self.pool.capacity {
		fresl := self.SlotStat(RDXFfree)
		nds := make([]RDXNode, fresl)
		for i, j := 0, 0; i < self.pool.capacity && j < fresl; i++ {
			if self.pool.nodes[i] == nil {
				self.pool.nodes[i] = &nds[j]
				j++
			}
		}
		n := self.pool.nodes[0]
		self.pool.index = 1
		self.pool.nodes[0] = nil
		return n
	}
	self.pool.index++
	n := self.pool.nodes[self.pool.index-1]
	if index := self.pool.index - 1; index > 0 {
		self.pool.nodes[index] = nil
	}
	return n
}

// SlotStat returns either number of free slots available
// or number of occupied slots.
func (self *Radix) SlotStat(flag byte) (count int) {
	switch flag {
	case RDXFoccupied:
		for i := 0; i < self.pool.capacity; i++ {
			if self.pool.nodes[i] != nil {
				count++
			}
		}
	case RDXFfree:
		for i := 0; i < self.pool.capacity; i++ {
			if self.pool.nodes[i] == nil {
				count++
			}
		}
	}
	return count
}

func (self *Radix) SetRemoveCB(callback func(*RDXNode)) {
	self.onRemove = callback
}

// Insert inserts into a particular node.
// func (self *Radix) Insert(str string, value interface{}) {
func (self *Radix) Insert(str string, value interface{}) {
	if self.root.Next == nil {
		var res *RDXNode = self.newNode()
		res.Key = str
		res.Value = value
		res.size = len(str)
		res.h = byte(res.Key[0])
		res.Isbnd = true
		self.root.Next = res
		return
	}
	self.root.Next.insert(&str, value, len(str), self)
}

// InsertR is a routine with same functionality as `Insert` but returns
// a pointer to the last/insert node. It starts from root node.
func (self *Radix) InsertR(str string, value interface{}) *RDXNode {
	var ret **RDXNode = new(*RDXNode)
	*ret = nil
	if self.root.Next == nil {
		self.Insert(str, value)
		*ret = self.root.Next
		return *ret
	}
	var ct int
	self.root.Next.retrinsert(&str, value, len(str), self, ret, &ct)
	// fmt.Println("ct ", ct)
	// fmt.Printf("(InsertR)N I.\n\tret(%p)\n\tself.root.next(%v)\n", ret, self.root.Next)
	return *ret
}

// InsertRFrom is a routine similar to `InsertR` but starts from `node` given in
// function arguments.
func (self *Radix) InsertRFrom(node *RDXNode, str string, value interface{}) (*RDXNode, int) {
	var ret **RDXNode = new(*RDXNode)
	*ret = nil
	if self.root.Next == nil {
		self.Insert(str, value)
		*ret = self.root.Next
		return *ret, len(str)
	}
	// node.Link = node.retrinsert(&str, value, len(str), self, ret)
	var ct int
	node.retrinsert(&str, value, len(str), self, ret, &ct)
	// fmt.Println("ct is :", ct)
	// fmt.Printf("(InsertRFrom)info\n\tret(%p)\n\t*ret(%v)\n\tn(%v)\n\tnode(%v)\n", ret, *ret, n, node)
	// if node.Link == nil {
	//   fmt.Println("NODELINK IS NIL")
	//   node.Link = *ret
	// }
	// if node != nil {
	// 	fmt.Printf("(InsertRFrom)node is not nil.\n\tnode.Next(%v)\n\tnode.Link(%v)\n\tnode.Key(%v)\n\tnode(%v)\n", node.Next, node.Link, node.Key, node)
	// }
	return *ret, ct
}

// Find tries to find a particular chain in the tree ( containing input word )
// and returns head, tail and their result as concatenated string. String result
// must be checked against the original word ( or alternatively checking boundaries using tail node )
func (self *Radix) Find(key string) (*RDXResult, *RDXNode, string) {
	var (
		node   *RDXResult = NewRDXResult()
		tail   *RDXNode
		result bytes.Buffer
	)
	// self.root.find(&key, len(key), node)
	self.root.rfind(&key, len(key), node)
	result.WriteString("")
	for curr := node; curr != nil; curr = curr.Next {
		if (*curr.Link) != nil {
			result.WriteString((*curr.Link).Key)
			tail = (*curr.Link)
		}
	}
	return node, tail, result.String()
}

// FindFrom is a method that works exactly as Find except it starts its search
// from a given node `node`.
func (self *Radix) FindFrom(root *RDXNode, key string) (*RDXResult, *RDXResult, string) {
	var (
		node   *RDXResult = NewRDXResult()
		tail   *RDXResult
		result bytes.Buffer
	)
	root.rfind(&key, len(key), node)
	result.WriteString("")
	for curr := node; curr != nil; curr = curr.Next {
		if (*curr.Link) != nil {
			result.WriteString((*curr.Link).Key)
			// NOTE
			//  this returns the parent node
			//  unlinke `Find` method
			tail = curr
		}
	}
	return node, tail, result.String()
}

// Remove is the parent of remove routine and returns wether the operation went successfully
func (self *Radix) Remove(key string) bool {
	var sc bool = false
	self.root = self.root.remove(&key, len(key), &sc, self)

	return sc
}

// Traverse recursively finds all boundary nodes.
func (self *Radix) Traverse(node *RDXNode) (subs []string, err error) {
	var stack []*RDXNode
	node.rtraverse(&subs, &stack)
	return subs, nil
}

// SetOpts sets the option byte `Opts`.
func (self *RDXNode) SetOpts(opts byte) {
	self.Opts = opts
}

// SetOptsO sets the option byte `Opts` and return its old value.
func (self *RDXNode) SetOptsO(opts byte) byte {
	oldval := self.Opts
	self.Opts = opts
	return oldval
}

// SetValue sets the internal `Value` interface.
func (self *RDXNode) SetValue(value interface{}) {
	self.Value = nil
	self.Value = value
}

// GetValue returns the internal `Value` interface.
func (self *RDXNode) GetValue() interface{} {
	return self.Value
}

// GetOpts returns the option byte `Opts`.
func (self *RDXNode) GetOpts() byte {
	return self.Opts
}

// IsLeaf returns true if current is a leaf node.
func (self *RDXNode) IsLeaf() bool {
	return self.Next != nil
}

// HasParent returns true if current node has a parent.
func (self *RDXNode) HasParent() bool {
	return self.Link != nil
}

// Release adds the given node into the memory pool
// for recycling. It is cruical to ensure that
// this node is not being used by any other part
// of the program, otherwise it can have unperdictable
// side effects or cause fatal failures.
func (self *Radix) Release(node *RDXNode) {
	if self.pool.index != -1 || self.pool.index != self.pool.capacity {
		for i := 0; i < self.pool.index; i++ {
			if self.pool.nodes[i] == nil {
				self.pool.nodes[i] = node
				return
			}
		}
	}
}

// Split handles a special case in insertion
func (self *RDXNode) split(k int, l int, radix *Radix) {
	n := radix.newNode()
	n.Key = self.Key[k:]
	n.h = byte(n.Key[0])
	n.size = l
	n.Link = self.Link
	n.Isbnd = self.Isbnd
	n.NoMrg = self.NoMrg
	n.Value = self.Value
	self.Link = n
	self.Key = self.Key[:k]
	self.size = k
	self.h = byte(self.Key[0])
	self.Isbnd = false
	self.NoMrg = false
	self.Value = nil
}

// SplitV is the same as `Split` but for when `k==self.size`. It also returns
// the new branch.
func (self *RDXNode) splitV(k int, l int, radix *Radix, value interface{}) *RDXNode {
	n := radix.newNode()
	n.Key = self.Key[k:]
	n.h = byte(n.Key[0])
	n.size = l
	n.Link = self.Link
	n.Isbnd = self.Isbnd
	n.NoMrg = self.NoMrg
	n.Value = self.Value
	self.Link = n
	self.Key = self.Key[:k]
	self.size = k
	self.h = byte(self.Key[0])
	self.Isbnd = true
	self.NoMrg = true
	self.Value = value

	return n
}

// RTraverse recursively finds all boundary nodes.
func (self *RDXNode) rtraverse(subs *[]string, stack *[]*RDXNode) {
	var (
		stacklen int  = len(*stack)
		wrote    bool = false
	)
	if self == nil {
		return
	}
	goto HEAD

WRITER:
	if stacklen > 0 {
		var (
			buf bytes.Buffer
			s   string
		)
		for _, sn := range *stack {
			buf.WriteString(sn.Key)
		}
		if buf.Len() > 0 {
			s = buf.String()
			(*subs) = append((*subs), s)
		}
	}
	goto MAIN

HEAD:
	(*stack) = append((*stack), self)
	stacklen = len(*stack)
	if self.Isbnd && wrote == false {
		wrote = true
		goto WRITER
	}

MAIN:
	if self.Isbnd && self.Link == nil {
		if stacklen > 1 {
			(*stack) = (*stack)[:stacklen-1]
		} else {
			(*stack) = (*stack)[:0]
		}
		stacklen = len(*stack)
	}
	self.Link.rtraverse(subs, stack)
	self.Next.rtraverse(subs, stack)
	if stacklen >= 1 {
		(*stack) = (*stack)[:stacklen-1]
	}

	return
}

// Rinsert recursively inserts the new node ( or replace it )
func (self *RDXNode) rinsert(key *string, value interface{}, n int, radix *Radix) *RDXNode {
	if self == nil {
		var res *RDXNode = radix.newNode()
		res.Key = (*key)[:n]
		res.Value = value
		res.size = n
		res.h = byte(res.Key[0])
		res.Isbnd = true
		return res
	}
	var k int = self.size
	if (*key)[0] != self.h {
		k = 0
		self.Next = self.Next.rinsert(key, value, n, radix)
		return self
	} else {
		for i := 0; i < self.size; i++ {
			if i == n || ((*key)[i]^self.Key[i]) != 0 {
				k = i
				break
			}
		}
		if k < n {
			if k < self.size {
				self.split(k, self.size-k, radix)
			}
			var nn string = (*key)[k:]
			self.Link = self.Link.rinsert(&nn, value, n-k, radix)
		}
		return self
	}
}

// NOTE: this is under development.
// retrinsert is similar to `rinsert`. It updates a double pointer to
// the last inserted/replaced node.
func (self *RDXNode) retrinsert(key *string, value interface{}, n int, radix *Radix, ret **RDXNode, ct *int) *RDXNode {
	// fmt.Printf("(retrinsert)ret info.\nret(%v)\n", *ret)
	if self == nil {
		var res *RDXNode = radix.newNode()
		res.Key = (*key)[:n]
		res.Value = value
		res.size = n
		res.h = byte(res.Key[0])
		res.Isbnd = true
		*ct = n
		*ret = nil
		*ret = res
		// fmt.Printf("(retrinsert) setting ret in case 1.\n\tret(%v)\n", ret)
		return res
	}
	var k int = self.size
	for i := 0; i < self.size; i++ {
		if i == n || ((*key)[i]^self.Key[i]) != 0 {
			k = i
			break
		}
	}
	// fmt.Println("(retrinsert) K FOR *key and k. (*key, k, n, self.Key)", *key, k, n, self.Key, self.size)
	if k == n && (len(*key) == len(self.Key) && *key == self.Key) {
		// fmt.Println("(retrinsert) INSIDE THIS IMPORTANT SITUATION.", *key, self.Key)
		*ret = nil
		*ret = self
		*ct = n
		return self
	}
	if k == 0 {
		self.Next = self.Next.retrinsert(key, value, n, radix, ret, ct)
	} else if k <= n {
		if k < self.size {
			// fmt.Println("(retrinsert) SPLLITING.", *key, self.Key, ct)
			if k == n {
				// fmt.Println("SPLITTINV", *key, self.Key)
				// nbr := self.splitV(k, self.size-k, radix, value)
				self.splitV(k, self.size-k, radix, value)
				*ret = nil
				// *ret = nbr
				*ret = self
				*ct = n
				// fmt.Println("SPLITTING RET IS ", *ret)
			} else {
				self.split(k, self.size-k, radix)
			}
		} else if k == self.size {
			self.Isbnd = true
		}
		if k != n || k < self.size {
			var nn string = (*key)[k:]
			// fmt.Printf("(retrinsert) setting ret in case 3.\n\tret(%v)\n", ret)
			self.Link = self.Link.retrinsert(&nn, value, n-k, radix, ret, ct)
		}
	}
	// fmt.Printf("(retrinsert)before final return.\n\tself.Next(%v)\n\tself.Link(%v)\n\tself(%v)\n\tret(%v)\n", self.Next, self.Link, self, *ret)

	return self
}

// Insert is non-recursive routine ( parent of rinsert ) used to reduce function call overhead.
// It is very efficient and fast when cache lines are warmed up
func (self *RDXNode) insert(key *string, value interface{}, n int, radix *Radix) {
	var (
		curr *RDXNode = self
		ln   *RDXNode = curr
		k    int
	)
	for {
		if curr == nil {
			var res *RDXNode = radix.newNode()
			res.Key = (*key)[:n]
			res.Value = value
			res.size = n
			res.h = byte(res.Key[0])
			res.Next = ln.Next
			res.Isbnd = true
			ln.Next = res
			return
		} else if byte((*key)[0]) != curr.h {
			curr = curr.Next
			continue
		}
		k = curr.size
		for i := 0; i < curr.size; i++ {
			if i == n || ((*key)[i]^curr.Key[i]) != 0 {
				k = i
				break
			}
		}
		if k < n {
			if k < curr.size {
				curr.split(k, curr.size-k, radix)
			}
			var nn string = (*key)[k:]
			var sz int = n - k
			if curr.Link == nil {
				var res *RDXNode = radix.newNode()
				res.Key = nn
				res.Value = value
				res.size = sz
				res.h = byte(res.Key[0])
				curr.Link = res
				res.Isbnd = true
				return
			}
			curr.Link = curr.Link.rinsert(&nn, value, sz, radix)
			return
		}
		return
	}
}

// find recursively locates the particular node. It returns a pointer
// to `RDXResult` struct.
func (self *RDXNode) rfind(key *string, n int, res *RDXResult) *RDXNode {
	if self == nil {
		return nil
	}
	// fmt.Println("(rfind): ", *key, self.Key, self.Next, self.Link)
	if byte((*key)[0]) != self.h {
		// fmt.Println("(rfind):key[0]!=self.h. ", self.Next, self.Link)
		return self.Next.rfind(key, n, res)
	} else {
		var k int = self.size
		for i := 0; i < self.size; i++ {
			if i == n || ((*key)[i]^self.Key[i]) != 0 {
				k = i
				break
			}
		}
		if k == n {
			*res.Link = self
			res.Next = NewRDXResult()
			res = res.Next
			return self
		}
		if k == self.size {
			*res.Link = self
			res.Next = NewRDXResult()
			res = res.Next
			nn := (*key)[k:]
			return self.Link.rfind(&nn, n-k, res)
		}
	}
	return nil
}

// join merges two nodes ( in deleting scenario )
func (self *RDXNode) join() {
	var (
		buff          = bytes.Buffer{}
		ln   *RDXNode = self.Link
		blen int
	)
	buff.WriteString(self.Key)
	buff.WriteString(ln.Key)
	blen = buff.Len()
	self.Key = buff.String()
	self.size = blen
	self.Link = ln.Link
	self.Isbnd = ln.Isbnd
}

// remove recursively locates the node and removes the chain from tree
func (self *RDXNode) remove(key *string, n int, sc *bool, radix *Radix) *RDXNode {
	if self == nil {
		return nil
	}
	var k int = self.size
	for i := 0; i < self.size; i++ {
		if i == n || ((*key)[i]^self.Key[i]) != 0 {
			k = i
			break
		}
	}
	if k == n {
		*sc = true
		if radix.onRemove != nil {
			radix.onRemove(self)
		}
		radix.Release(self)
		return self.Next
	}
	if k == 0 {
		self.Next = self.Next.remove(key, n, sc, radix)
	} else if k == self.size {
		nn := (*key)[k:]
		self.Link = self.Link.remove(&nn, n-k, sc, radix)
		if self.Link != nil && self.Link.Next == nil {
			// NOTE: *IMPORTANT*
			//  NoMrg is for non mergable nodes
			//  side effects are not clear
			if self.NoMrg == false {
				self.join()
			}
			// END
		}
	}
	return self
}
