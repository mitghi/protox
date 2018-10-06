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
* TODO:
* . write documentation
* . unify interface
*
* type QueueInterface interface {
* 	Push(interface{}) bool
* 	Pop() interface{}
* 	Size() int
* 	Head() interface{}
* 	Peek() interface{}
* 	Empty() bool
* }
**/

type BoundedQueue struct {
	items      []interface{}
	head       int64
	tail       int64
	capacity   int
	size       int
	autoresize bool
}

func NewBoundedQueue(capacity int, autoresize bool) *BoundedQueue {
	var bq *BoundedQueue = &BoundedQueue{
		items:      make([]interface{}, capacity, capacity),
		capacity:   capacity,
		autoresize: autoresize,
	}
	return bq
}

func (bq *BoundedQueue) Push(item interface{}) bool {
	if bq.size == bq.capacity {
		if !bq.autoresize {
			return false
		}
		bq.resize(bq.capacity * 2)
	}
	bq.items[bq.tail%int64(bq.capacity)] = item
	bq.tail++
	bq.size++
	return true
}

func (bq *BoundedQueue) Pop() (item interface{}) {
	if bq.size == 0 {
		return nil
	}
	if bq.autoresize && (bq.size < bq.capacity/8 && bq.size > 8) {
		bq.resize(bq.capacity / 2)
	}
	item = bq.items[bq.head%int64(bq.capacity)]
	bq.head++
	bq.size--
	return item
}

func (bq *BoundedQueue) Head() interface{} {
	if bq.size > 0 {
		return bq.items[bq.head%int64(bq.capacity)]
	}
	return nil
}

func (bq *BoundedQueue) Peek() interface{} {
	if bq.size > 0 {
		return bq.items[(bq.tail-1)%int64(bq.capacity)]
	}
	return nil
}

func (bq *BoundedQueue) Empty() bool {
	return bq.size == 0
}

func (bq *BoundedQueue) Size() int {
	return bq.size
}

func (bq *BoundedQueue) resize(nsize int) {
	var (
		ns   []interface{} = make([]interface{}, nsize, nsize)
		qcap int64         = int64(bq.capacity)
	)
	for i := bq.head; i < bq.tail; i++ {
		ns[i-bq.head] = bq.items[i%qcap]
	}
	bq.capacity = nsize
	bq.head = 0
	bq.tail = int64(bq.size)
	bq.items = ns
}
