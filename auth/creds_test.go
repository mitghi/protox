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

package auth

import "testing"

func TestCredsMatch(t *testing.T) {
	creds := []*Creds{
		{Username: "test", Password: "passwd", ClientId: "clid"},
		{Username: "test", Password: "passwd", ClientId: "clid"},
		{Username: "test2", Password: "passwd2", ClientId: "clid2"},
		{Username: "test 2", Password: "passwd 2", ClientId: "clid 2"},
		{Username: "test2 ", Password: "passwd2 ", ClientId: "clid2 "},
	}
	if !creds[0].Match(creds[1]) {
		t.Fatal("invalid result")
	}
	if creds[2].Match(creds[3]) {
		t.Fatal("invalid result")
	}
	creds[2].cleanInput(creds[2])
	creds[4].cleanInput(creds[4])
	au, ap, ac := creds[2].GetCredentials()
	bu, bp, bc := creds[4].GetCredentials()

	if valid := ((au == bu) && (ap == bp)) && (ac == bc); !valid {
		t.Fatal("invalid result")
	}
}
