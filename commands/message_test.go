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

import (
	"fmt"
	"testing"
)

func TestMessage(t *testing.T) {
	var (
		f    string = addArray("SET", "a", "b")
		fexp string = "*3:5\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\nb\r\n"
	)
	if f != fexp {
		t.Fatal("invalid value")
	}
	fmt.Println(f)
}

func TestArrayHeader(t *testing.T) {
	var (
		table []struct {
			header    []byte
			elemsize  int
			totalsize int
			valid     bool
		} = []struct {
			header    []byte
			elemsize  int
			totalsize int
			valid     bool
		}{
			{[]byte("*4:18\r\n"), 4, 18, true},
			{[]byte("*4:\r\n"), 0, 0, false},
			{[]byte("*4:\r"), 0, 0, false},
			{[]byte("*4:"), 0, 0, false},
			{[]byte("*4"), 0, 0, false},
		}
	)
	for i, e := range table {
		if e.valid {
			res, err := parseArrayHeader(e.header)
			if err != nil {
				fmt.Println(string(e.header))
				t.Fatalf("[ case %d ] inconsistent state, expected err==nil. (%+v),err:%s", i, e, err)
			}
			if res.es != e.elemsize || res.ts != e.totalsize {
				t.Fatalf("[ case %d ] incorrect values, expected es==%d got %d, expected ts==%d got %d.", i, e.elemsize, res.es, e.totalsize, res.ts)
			}
		} else {
			_, err := parseArrayHeader(e.header)
			if err == nil {
				t.Fatalf("expected error!=nil. (%+v)", e)
			}
		}
	}
}
