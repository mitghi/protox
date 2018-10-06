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

// Stack is implemented using `interface` slice.
type Stack []interface{}

// NewQueue alloactes and initializes a new
// 'Deque' and returns a pointer to it.
func NewStack() *Stack {
	var s *Stack = &Stack{}

	return s
}

// Push pushes an `item` into the Stack.
func (s *Stack) Push(item interface{}) {
	ls := *s
	ls = append(ls, item)
	*s = ls
}

// Pop removes and returns the most recent item
// from the Stack.
func (s *Stack) Pop() (item interface{}) {
	ls := *s
	var l int = len(ls)
	if l == 0 {
		return nil
	}
	item = ls[l-1]
	ls[l-1] = nil
	ls = append(ls[:0], ls[:l-1]...)
	*s = ls
	return item
}

// Size returns current number of elements in Stack.
func (s *Stack) Size() int {
	ls := *s
	return len(ls)
}

// Head returns bottom of the Stack without removing it.
func (s *Stack) Head() (item interface{}) {
	if len(*s) == 0 {
		return nil
	}
	item = (*s)[0]
	return item
}

// Peek returns top of Stack without removing it.
func (s *Stack) Peek() (item interface{}) {
	var l int = len(*s)
	if l == 0 {
		return nil
	}
	item = (*s)[l-1]
	return item
}

// Empty returns true if Stack is empty.
func (s *Stack) Empty() bool {
	ls := *s
	return ls.Size() == 0
}
