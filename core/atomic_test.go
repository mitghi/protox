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

package core

import (
	"fmt"
	"testing"
)

func TestAtomicCounter(t *testing.T) {
	cb := func(val int64) {
		fmt.Println("counter is triggered and invoked this function. Value ", val)
	}
	ac, err := NewAtomicCounter(0, 100, cb)
	if err != nil {
		t.Fatal("cannot get a new counter with valid arguments.")
	}
	var ch chan struct{} = make(chan struct{}, 5)

	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				_ = ac()
			}
			ch <- struct{}{}
		}()
	}
	for i := 0; i < 5; i++ {
		<-ch
	}
	close(ch)
	if ct := ac(); ct >= 100 && ct <= 0 {
		t.Fatal("invalid value.")
	}

}

func TestExpiringAtomicCounter(t *testing.T) {
	cb := func(val int64) {
		fmt.Println("counter is triggered and invoked this function. Value ", val)
	}
	ecb := func(val int64) {
		fmt.Println("expired counter has triggered this callback.")
	}
	ac, err := NewExpiringAtomicCounter(0, 100, cb, ecb)
	if err != nil {
		t.Fatal("cannot get a new counter with valid arguments.")
	}
	var ch chan struct{} = make(chan struct{}, 5)

	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 20; i++ {
				_, status := ac()
				if status == ACExpired {
					t.Fatal("invalid value")
				}
			}
			ch <- struct{}{}
		}()
	}
	for i := 0; i < 5; i++ {
		<-ch
	}
	close(ch)
	for i := 0; i < 5; i++ {
		if _, status := ac(); status != ACExpired {
			t.Fatal("invalid status")
		}
	}
}
