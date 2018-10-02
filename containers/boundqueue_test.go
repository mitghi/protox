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

import "testing"

func TestQueue(t *testing.T) {
	var (
		s *BoundedQueue = NewBoundedQueue(8, true)
	)
	for i := 0; i < 10; i++ {
		s.Push(i)
	}
	if Size := s.Size(); Size != 10 {
		t.Fatalf("invalid stack Size, expected Size==10, got %d.", Size)
	}
	if Peek := s.Peek(); Peek == nil || (Peek != nil && Peek.(int) != 9) {
		t.Fatalf("invalid value, expected Peek==10, got %v", Peek)
	}
	if Head := s.Head(); Head == nil && (Head != nil && Head.(int) != 0) {
		t.Fatalf("invalid value, expected Head==9, got %v", Head)
	}
	if item := s.Pop(); item == nil || (item != nil && item.(int) != 0) {
		t.Fatalf("invalid value, expected item==9, got %v", item)
	}
	for i := 0; i < 8; i++ {
		_ = s.Pop()
	}
	if Size := s.Size(); Size != 1 {
		t.Fatalf("expected stack to be of Size 1, got %d", Size)
	}
	_ = s.Pop()
	if Size := s.Size(); Size != 0 {
		t.Fatalf("expected empty stack to be of Size 0, got %d", Size)
	}
	if s.Pop() != nil {
		t.Fatal("inconsistent state")
	}
	if s.Head() != nil {
		t.Fatal("inconsistent state")
	}
	if s.Peek() != nil {
		t.Fatal("inconsistent state")
	}
	for i := 0; i < 10; i++ {
		_ = s.Pop()
	}
}
