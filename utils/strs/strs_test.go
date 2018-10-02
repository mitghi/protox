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

package strs

import "testing"

func TestMatch(t *testing.T) {
	a, b := "a/*/path", "a/simple/path"
	if !Match(a, b, "/", "*") {
		t.Fatal("expected true, cannot match a-b.", a, b)
	}
	a, b = "a/*/path", "a/simple/location"
	if Match(a, b, "/", "*") {
		t.Fatal("expected false, wrong match between a-b.", a, b)
	}
	a, b = "a/*/path/a", "a/simple/location"
	if Match(a, b, "/", "*") {
		t.Fatal("expected false, wrong match between a-b.", a, b)
	}
	a, b = "a/*/location", "a/simple/location/a"
	if Match(a, b, "/", "*") {
		t.Fatal("expected false, wrong match between a-b.", a, b)
	}
}

func TestMatch2(t *testing.T) {
	a, b := "a/*/path", "a/simple/path"
	if !Match(a, b, "/", "*") {
		t.Fatal("expected true, cannot match a-b.", a, b)
	}
	a, b = "a/*", "a/simple/path"
	if !Match(a, b, "/", "*") {
		t.Fatal("expected true, cannot match a-b.", a, b)
	}
	a, b = "a/*/path/thing", "a/simple/path"
	if Match(a, b, "/", "*") {
		t.Fatal("expected false, invalid match between  a-b.", a, b)
	}
	a, b = "a/*/simple/location", "a/another/simple/thing"
	if Match(a, b, "/", "*") {
		t.Fatal("expected false, invalid match between  a-b.", a, b)
	}
}

func TestRecFind(t *testing.T) {
	spl := byte('/')
	a, b := "a/*/path", "a/simple/path"
	if v := LCPWithSeparator(a, b, spl); v != 7 {
		t.Fatal("expected 7, got:", v)
	}
	a, b = "a/*", "a/simple/path"
	if v := LCPWithSeparator(a, b, spl); v != 2 {
		t.Fatal("expected 2, got:", v)
	}
	a, b = "a/*/path/thing", "a/simple/path"
	if v := LCPWithSeparator(a, b, spl); v != 7 {
		t.Fatal("expected 7, got:", v)
	}
	a, b = "a/*/simple/location", "a/another/simple/thing"
	if v := LCPWithSeparator(a, b, spl); v != 11 {
		t.Fatal("expected 11, got:", v)
	}
}
