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

import "errors"

// OrderedSet Error Messages
var (
	EOSINVAL error = errors.New("OrderedSet: item not found.")
	EOSIGC   error = errors.New("OrderedSet: error in comparision of generic value types.")
)

// SetItemInterface is the interface which must be implemented
// by any type willing to get stored in the container. Note that
// container values must be homogenious.
type SetItemInterface interface {
	More(SetItemInterface) bool
	Less(SetItemInterface) bool
	Equal(SetItemInterface) bool
	Value() interface{}
	Compare(SetItemInterface) (int, error)
}

// OrderedSet is the implementation for ordered sets. It has
// time complexity of O(log(N)+ k) where k is number of elements
// in left and right wing when items are identical but carry
// different values.
type OrderedSet struct {
	list []SetItemInterface
}

// NewOrderedSet returns a pointer to a new allocated
// OrderedSet.
func NewOrderedSet() *OrderedSet {
	return &OrderedSet{}
}

// insertAt is a receiver that inserts a value in an specific
// slice slot.
func (s *OrderedSet) insertAt(index int, value SetItemInterface) {
	var (
		sl int                = len(s.list)
		ns []SetItemInterface = s.list
	)
	if sl == cap(ns) {
		ns = make([]SetItemInterface, sl+1, cap(ns)+8)
		copy(ns, s.list[:index])
	} else {
		ns = append(ns, nil)
	}
	copy(ns[index+1:], s.list[index:])
	ns[index] = value
	s.list = ns
}

// Insert is a receiver method which takes a argument of
// compatible `SetItemInterface` and stores the value in
// the internal container. It returns true when insertion
// is successfull and returns false in case `item` already
// exists in the container.
func (s *OrderedSet) Insert(item SetItemInterface) bool {
	if s.Exists(item) {
		return false
	}
	for i := 0; i < s.Count()-1; i++ {
		if s.list[i].More(item) {
			s.insertAt(i, item)
			return true
		}
	}
	s.list = append(s.list, item)
	return true
}

// Remove is a receiver method which takes a argument of
// compatible type `SetItemInterface` and removes it from
// the container. It returns false in case of non existing
// item and returns true when item is succesfully removed
// from the container.
func (s *OrderedSet) Remove(item SetItemInterface) bool {
	var (
		index int
		err   error
	)
	if index, err = s.Index(item); index < 0 || err != nil {
		// TODO
		// . handlethis error appropirately.
		return false
	}
	copy(s.list[index:], s.list[index+1:])
	s.list[len(s.list)-1] = nil
	s.list = s.list[:len(s.list)-1]
	return true
}

// Pop removes and returns the oldest item. It returns
// `nil` when unsuccessfull/empty.
func (s *OrderedSet) Pop() (item SetItemInterface) {
	if s.Count() == 0 {
		return nil
	}
	item = s.list[0]
	copy(s.list[0:], s.list[1:])
	s.list[len(s.list)-1] = nil
	s.list = s.list[:len(s.list)-1]

	return item
}

// Exists is a receiver method which takes a argument of
// compatible type `SetItemInterface` and returns whether
// the given item exists in the set or not.
func (s *OrderedSet) Exists(item SetItemInterface) bool {
	_, err := s.Index(item)
	if err != nil {
		return false
	}
	return true
}

// Count returns the current size of the container.
func (s *OrderedSet) Count() int {
	return len(s.list)
}

// Index is a receiver function that takes a item of compatible
// type `SetItemInterface` as argument and finds its index in the
// underlaying container. It returns an error when unable to locate
// the value or a non homogenious value, otherwise it reutns the
// index with error code set to null.
func (s *OrderedSet) Index(item SetItemInterface) (int, error) {
	var (
		lb  int = 0
		rb  int = s.Count() - 1
		mid int
	)
	for lb <= rb {
		mid = lb + ((rb - lb) / 2)
		comp, err := s.list[mid].Compare(item)
		if err != nil {
			return -1, EOSIGC
		}
		if comp == 0 {
			return mid, nil
		} else if comp > 0 {
			rb = mid - 1
		} else if comp < 0 {
			lb = mid + 1
		} else {
			// search left and right sides
			// to find correct values in case
			// structs are used as values.
			for i := mid; i < s.Count()-1; i++ {
				if s.list[i+1].Equal(item) {
					return (i + 1), nil
				} else if s.list[i].Less(s.list[i+1]) {
					break
				}
			}
			for i := mid; i > 0; i-- {
				if s.list[i-1].Equal(item) {
					return (i - 1), nil
				} else if s.list[i].More(s.list[i-1]) {
					break
				}
			}
			return -1, EOSINVAL
		}
	}
	return -1, EOSINVAL
}

// OrsItem is the base structure for Ordered Set elements. It
// conforms to `SetItemInterface` interface and can be embedded
// into other stuctures.
type OrsItem struct {
	value interface{}
}

// Value returns the actual value from the struct.
func (osi *OrsItem) Value() interface{} {
	return osi.value
}

// Compare is a receiver method which takes a argument of
// compatible `SetItemInterface` type and compares the
// given argument with itself. It returns -1, 0, 1 to indicate
// lessThan, equal, moreThan respectively and returns -2 with
// an error to indicate any problem related to the values. Note that
// interface methods must be reimplemented in the struct willing to
// conform to `SetItemInterface` interface.
func (osi *OrsItem) Compare(item SetItemInterface) (int, error) {
	// nv := item.Value()
	// if osi.value == nv {
	// 	return 0, nil
	// } else if osi.value > nv {
	// 	return 1, nil
	// } else if osi.value < nv {
	// 	return -1, nil
	// }
	return -2, EOSIGC
}

// More is a receiver method that checks whether the given
// argument is bigger than the own value.
func (osi *OrsItem) More(item SetItemInterface) bool {
	isMore, _ := osi.Compare(item)
	if isMore == 1 {
		return true
	}
	return false
}

// Less is a receiver method that checks whether the given
// argument is less than the own value.
func (osi *OrsItem) Less(item SetItemInterface) bool {
	isLess, _ := osi.Compare(item)
	if isLess == -1 {
		return true
	}
	return false
}

// Equal is a receiver method that checks whether the given
// argument is equal to the own value.
func (osi *OrsItem) Equal(item SetItemInterface) bool {
	isEqual, _ := osi.Compare(item)
	if isEqual == 0 {
		return true
	}
	return false
}
