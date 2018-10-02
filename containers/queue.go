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

// Queue is implemented using `interface` slice.
type Queue []interface{}

// NewQueue alloactes and initializes a new
// 'Queue' and returns a pointer to it.
func NewQueue() *Queue {
	var q *Queue = &Queue{}

	return q
}

// Push pushes `item` into back of Queue.
func (q *Queue) Push(item interface{}) bool {
	ls := *q
	ls = append(ls, item)
	*q = ls
	return true
}

// Pop removes front most item from the Queue.
func (q *Queue) Pop() (item interface{}) {
	ls := *q
	var l int = len(ls)
	if l == 0 {
		return nil
	}
	item = ls[0]
	ls = append(ls[:0], ls[1:]...)
	*q = ls
	return item
}

// Size returns current number of elements in Queue.
func (q *Queue) Size() int {
	ls := *q
	return len(ls)
}

// Head returns bottom of the Queue without removing it.
func (q *Queue) Head() (item interface{}) {
	if len(*q) == 0 {
		return nil
	}
	item = (*q)[0]
	return item
}

// Peek returns top of the Queue without removing it.
func (q *Queue) Peek() (item interface{}) {
	var l int = len(*q)
	if l == 0 {
		return nil
	}
	item = (*q)[l-1]
	return item
}

// Empty returns true if Queue is empty.
func (q *Queue) Empty() bool {
	ls := *q
	return ls.Size() == 0
}
