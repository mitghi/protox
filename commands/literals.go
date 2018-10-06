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

package commands

// ArrayNode is a struct that represents an parsed array from textual protocol line. It is used to represent an array into the abstract syntax tree.
type ArrayNode struct {
	// next points to the next node in the chain.
	next ParseNodeInterface
	// elems is a slice of elements ( array descendants ).
	elems []ParseElementInterface
	// info contains all necessary array information such as
	// capacity, length and its elemnt type.
	info *parseNodeInfo
}

// NewArrayNode allocates and initializes a new `ArrayNode`
// and returns a pointer to it.
func NewArrayNode() *ArrayNode {
	var result *ArrayNode = &ArrayNode{
		next:  nil,
		elems: nil,
		info: &parseNodeInfo{
			cap:    0,
			length: 0,
			etype:  PNDArrayLit,
		},
	}
	return result
}

// Value is a receiver method that returns the current value
// of array ( slice of element pointers ).
func (an *ArrayNode) Value() interface{} {
	return an.elems
}

// Next is a receiver method that returns the next node in the
// chain.
func (an *ArrayNode) Next() ParseNodeInterface {
	return an.next
}

// AddValue is a receiver method that pushes an argument of
// type `parseElementInterface` into the internal element slice.
func (an *ArrayNode) AddValue(value ParseElementInterface) {
	an.elems = append(an.elems, value)
}

// SetNext is a receiver method that sets the `next` pointer to the
// argument of type `ParseNodeInterface`.
func (an *ArrayNode) SetNext(node ParseNodeInterface) {
	an.next = node
}
