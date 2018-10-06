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
	"bytes"
	"fmt"
	"testing"
)

func TestParseString(t *testing.T) {
	var (
		table []struct {
			input     []byte
			expect    []byte
			remaining int
			valid     bool
		} = []struct {
			input     []byte
			expect    []byte
			remaining int
			valid     bool
		}{
			// test successfull parse
			{[]byte("$1\r\na\r\n"), []byte("a"), -1, true},
			{[]byte("$4\r\nabcd\r\n$1\r\na\r\n"), []byte("abcd"), 10, true},

			// test inconsistent length
			{[]byte("$4\r\na\r\n"), []byte(""), -1, false},

			// test invalid format
			{[]byte("$$1\r\na\r\n"), []byte(""), -1, false},
			{[]byte("$1\na\n"), []byte(""), -1, false},

			// test length violation
			{[]byte("$1\r\npayload\r\n$2\r\nab\r\n"), []byte(""), -1, false},
		}
	)
	for i, e := range table {
		res, remaining, err := parseString(e.input)
		if err != nil {
			if !e.valid {
				if e.remaining != remaining {
					t.Fatalf("[case %d] expected invalid case to have identical remaining value ( %d != %d )", i, e.remaining, remaining)
				}
				continue
			}
			t.Fatalf("[case %d] expected err==nil, got %s.", i, err)
		}
		if bytes.Compare(res, e.expect) != 0 {
			t.Fatalf("[case %d] expected bytes cmp == 0. (%s != %s)", i, string(res), string(e.input))
		}
		if e.remaining != remaining {
			t.Fatalf("[case %d] invalid remaining value between ( %d != %d )", i, e.remaining, remaining)
		}
	}
}

func TestReadArray(t *testing.T) {
	var values [][]byte = [][]byte{
		[]byte("*2:2\r\n$1\r\na\r\n$1\r\nb\r\n"),
		[]byte("*2:2\r\n*2:2\r\n$1\r\na\r\n$1\r\nb\r\n*2:2\r\n$1\r\nc\r\n$1\r\nd\r\n"),
	}
	for _, v := range values {
		result, etype, err := parseArray(&v)
		if err != nil {
			t.Fatalf("expected err==nil, got %+v", err)
		}
		fmt.Println("(TestReadArray) result is :", result, etype, err)
		fmt.Println(result)
	}
}

func TestParse(t *testing.T) {
	v := []byte("*2:2\r\n*2:2\r\n$1\r\na\r\n$1\r\nb\r\n*2:2\r\n$1\r\nc\r\n$1\r\nd\r\n")
	result, err := Parse(&v)
	if err != nil {
		t.Fatalf("expected err==nil, got %+v", err)
	}
	fmt.Println(result)
}
