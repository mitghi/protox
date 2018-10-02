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

package messages

import (
	"fmt"
	"testing"
)

func TestSubscribeBasic(t *testing.T) {
	var (
		s      *Subscribe = NewSub()
		cl     string     = "client1"
		topics [][]byte   = [][]byte{
			[]byte("a/simple/topic"),
			[]byte("a/simple/*"),
		}
	)
	for _, topic := range topics {
		if err := s.rinsert(topic, 0x1, cl); err != nil {
			t.Fatal("unable to insert.")
		}
		if err := s.rinsert(topic, 0x1, "client2"); err != nil {
			t.Fatal("unable to insert.")
		}
	}
	st := &substorage{}
	if err := s.rfind(topics[1], 0x1, st); err != nil {
		t.Fatal("unable to find subscribers")
	}
	fmt.Println(st)

	if err := s.rremove(topics[0], "client1"); err != nil {
		t.Fatal("unable to remove.", err)
	}
	s.debugPrint()

}
