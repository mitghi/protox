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

// Deque is implemented using `interface` slice.
type Deque []interface{}

// NewDeque alloactes and initializes a new
// 'Deque' and returns a pointer to it.
func NewDeque() *Deque {
	var d *Deque = &Deque{}

	return d
}

// PushFront pushes `item` into front of Deque.
func (d *Deque) PushFront(item interface{}) {
	ls := *d
	ls = append([]interface{}{item}, ls...)
	*d = ls
}

// PushBack pushes `item` into back of Deque.
func (d *Deque) PushBack(item interface{}) {
	ls := *d
	ls = append(ls, item)
	*d = ls
}

// PopFront removes front most item from the Deque.
func (d *Deque) PopFront() (item interface{}) {
	ls := *d
	var l int = len(ls)
	if l == 0 {
		return nil
	}
	item = ls[0]
	ls = append(ls[:0], ls[1:])
	*d = ls

	return item
}

// PopBack removes top most item from the Deque.
func (d *Deque) PopBack() (item interface{}) {
	ls := *d
	var l int = len(ls)
	if l == 0 {
		return nil
	}
	item = ls[l-1]
	ls[l-1] = nil
	ls = append(ls[:0], ls[:l-1]...)
	*d = ls

	return item
}

// Size returns current number of elements in Deque.
func (d *Deque) Size() int {
	ls := *d
	return len(ls)
}

// Head returns bottom of the Deque without removing it.
func (d *Deque) Head() (item interface{}) {
	if len(*d) == 0 {
		return nil
	}
	item = (*d)[0]

	return item
}

// Peek returns top of the Deque without removing it.
func (d *Deque) Peek() (item interface{}) {
	var l int = len(*d)
	if l == 0 {
		return nil
	}
	item = (*d)[l-1]

	return item
}

// Empty returns true if Deque is empty.
func (d *Deque) Empty() bool {
	ls := *d

	return ls.Size() == 0
}
