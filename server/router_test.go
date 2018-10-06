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

package server

/**
*
* NOTE:
* This is under development and will change, therefore no comments will be added
* till all methods are finished.
*
**/

import (
	"fmt"
	"testing"

	"github.com/mitghi/protox/utils/strs"
)

func TestAddSub(t *testing.T) {
	r := NewRouter()
	r.AddSub("client", "testing/a/simple/string", 1)
	r.AddSub("client2", "testing/*", 1)
	r.AddSub("client3", "testing/a/simulation", 1)
	r.AddSub("client4", "test/a", 1)
	r.AddSub("SPC", "testing/*", 1)
	r.subs.PrintV()
	fmt.Println("-------------------------------")
}

func TestAddSub2(t *testing.T) {
	r := NewRouter()
	r.AddSub("client", "testing/a/simple/string", 1)
	r.AddSub("client2", "testing/*", 1)
	r.AddSub("client3", "testing/a/simulation", 1)
	r.AddSub("client4", "test/a", 1)
	r.AddSub("SPC", "testing/*", 1)
	r.subs.PrintV()
	fmt.Println("-------------------------------")
	r.subs.Remove("testing/*")
	r.subs.Remove("test/a")
	r.subs.PrintV()
	fmt.Println("-------------------------------")

}

func TestAddSub3(t *testing.T) {
	r := NewRouter()
	r.AddSub("client", "a/simple/path", 1)
	r.AddSub("client", "a/simple/*/thing", 1)
	r.AddSub("client", "a/another/simple/thing", 1)
	r.AddSub("client", "a/another/simulating/thing", 1)
	r.AddSub("client", "aa/branch", 1)
	r.AddSub("client", "a/another/simul", 1)
	r.AddSub("client", "a/another/sim", 1)
	r.AddSub("client", "a/another/sima", 1)
	r.subs.PrintV()
	fmt.Println("-------------------------------")
	m, err := r.FindSub("a/simple/path")
	if err != nil {
		t.Fatal("err!=nil", err)
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func TestAddSub4(t *testing.T) {
	r := NewRouter()
	r.AddSub("client1", "a/simple/path", 1)
	r.AddSub("client2", "a/*", 1)
	r.AddSub("client3", "a/another/simple/thing", 1)
	r.AddSub("client4", "a/*/location", 1)
	r.subs.PrintV()
	fmt.Println("-------------------------------")
	m, err := r.FindRawSubscribers("a/another/simple/thing")
	if err != nil {
		t.Fatal("err!=nil, expected to be nil", err)
	}
	for k, v := range m {
		for _, s := range v {
			fmt.Printf("%s, %+v\n", k, s)
			fmt.Println(strs.Match(s.topic, "a/another/simple/thing", "/", "*"))
		}
	}
}

func TestFindSub(t *testing.T) {
	r := NewRouter()
	r.AddSub("client1", "a/simple/path", 1)
	r.AddSub("client2", "a/*", 1)
	r.AddSub("client3", "a/another/simple/thing", 1)
	r.AddSub("client4", "a/*/location", 1)
	r.subs.PrintV()
	fmt.Println("-------------------------------")
	m, err := r.FindSubC("a/another/simple/thing")
	if err != nil {
		t.Fatal("err!=nil, expected to be nil", err)
	}
	for k, v := range m {
		for _, s := range v {
			fmt.Printf("%s, %+v\n", k, s)
			fmt.Println(strs.Match(s.topic, "a/another/simple/thing", "/", "*"))
		}
	}

}

func TestCache(t *testing.T) {
	r := NewRouter()
	r.AddSub("client1", "a/awesome/topic", 1)
	r.AddSub("client2", "a/*/topic", 2)
	r.AddSub("client3", "a/*", 0)
	r.AddSub("client4", "b/topic", 0)
	r.AddSub("client2", "a/awesome/topic", 1)
	m, err := r.FindSubC("a/awesome/topic")
	if err != nil {
		t.Fatal("err!=nil, expected to be nil", err)
	}
	fmt.Println(m)
	for k, v := range m {
		for _, s := range v {
			fmt.Printf("%s, %+v\n", k, s)
			fmt.Println(strs.Match(s.topic, "a/another/simple/thing", "/", "*"))
		}
	}
	r.subs.cache.CachePrintV()
	r.subs.PrintV()
	fmt.Println("-------------------------------")
	r.RemoveSub("client1", "a/awesome/topic")
	r.subs.cache.CachePrintV()
	r.subs.PrintV()
}
