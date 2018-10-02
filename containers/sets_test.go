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

import (
	"errors"
	"fmt"
	"testing"
)

type avc struct {
	OrsItem
}

func newavc(value interface{}) *avc {
	return &avc{OrsItem: OrsItem{value: value}}
}

func (a *avc) Compare(item SetItemInterface) (int, error) {
	nv, ok := item.Value().(int)
	if !ok {
		return -2, errors.New("OrderedSet: error in comparing generic interface{} values, [ requires homogenious types ].")
	}
	if a.value.(int) == nv {
		return 0, nil
	} else if a.value.(int) > nv {
		return 1, nil
	} else if a.value.(int) < nv {
		return -1, nil
	}
	// TODO
	// . handle this is the unknown case
	return -2, errors.New("OrderedSet: error in comparing generic interface{} values.")
}

func (a *avc) More(item SetItemInterface) bool {
	isMore, _ := a.Compare(item)
	if isMore == 1 {
		return true
	}
	return false
}

func (a *avc) Less(item SetItemInterface) bool {
	isLess, _ := a.Compare(item)
	if isLess == -1 {
		return true
	}
	return false
}

func (a *avc) Equal(item SetItemInterface) bool {
	isEqual, _ := a.Compare(item)
	if isEqual == 0 {
		return true
	}
	return false
}

func TestOrderedSet(t *testing.T) {
	s := NewOrderedSet()
	for i := 2; i < 10; i++ {
		s.Insert(&avc{OrsItem: OrsItem{value: i}})
	}
	if s.list[0].Value().(int) != 2 {
		t.Fatal("Invalid value")
	}
	s.Insert(&avc{OrsItem: OrsItem{value: 0}})
	for _, val := range s.list {
		fmt.Println(val)
	}
	if s.list[0].Value().(int) != 0 {
		t.Fatal("Invalid value")
	}
	for _, val := range s.list {
		fmt.Println(val)
	}
}

func TestOrderSetRemove(t *testing.T) {
	s := NewOrderedSet()
	for i := 0; i < 10; i++ {
		s.Insert(&avc{OrsItem: OrsItem{value: i}})
	}
	if !s.Remove(&avc{OrsItem: OrsItem{value: 5}}) {
		t.Fatal("cannot remove from OrderedSet.")
	}
	if s.Remove(&avc{OrsItem: OrsItem{value: 1345123512344}}) {
		t.Fatal("inconsistent state, remove non existence value.")
	}
	for _, val := range s.list {
		fmt.Println(val)
	}
	if !s.Insert(&avc{OrsItem: OrsItem{value: 5}}) {
		t.Fatal("unable to insert a value into the orderedset.")
	}
	if s.Insert(&avc{OrsItem: OrsItem{value: 5}}) {
		t.Fatal("inserted an already existing value again into the orderedset.")
	}
	fmt.Println("----")
	for _, val := range s.list {
		fmt.Println(val)
	}
}

func TestOrderSetPop(t *testing.T) {
	s := NewOrderedSet()
	for i := 0; i < 10; i++ {
		s.Insert(&avc{OrsItem: OrsItem{value: i}})
	}
	value := s.Pop()
	if value == nil {
		t.Fatal("expected value!=nil")
	}
	// readd popped value
	s.Insert(value)
	// it must be equal to
	// old value.
	value = s.Pop()
	if value == nil {
		t.Fatal("expected value!=nil")
	}
	fmt.Println(value)
	fmt.Println("----")
	for _, val := range s.list {
		fmt.Println(val)
	}
}
