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

import "strings"

// FindPrefix finds longest common prefix between two strings
// and writes it to int pointer `pf`.
func FindPrefix(a *string, b *string, pf *int) {
	var (
		la int = len(*a)
		lb int = len(*b)
	)
	for i := 0; i < la; i++ {
		if i == lb || ((*b)[i] != '*' || (*a)[i] != (*b)[i]) {
			*pf = i
			return
		}
	}
	*pf = la
	return
}

// LCPWithSeparator is similar to `FindPrefix` but returns
// longest common prefix with evaluating wildcard positions.
func LCPWithSeparator(a string, b string, spl byte) int {
	var (
		la int = len(a)
		lb int = len(b)
	)
	for i := 0; i < la; i++ {
		if i == lb || (a[i] != b[i]) {
			if i < lb && i <= la {
				var l int = i
				for j := i; j < lb; j++ {
					if b[j] == spl {
						l = j
						break
					}
				}
				return i + LCPWithSeparator(a[i+1:], b[l:], spl)
			}
			return i
		}
	}
	return la
}

// Match returns wether two strings seperated by delimiter
// `spl` are exactly a match. First string can contain a wildcard
// `wildcard` which forces the algrithm to count it as a match.
// For example: `test/*/a` and `test/a/a` are exact matches
// because `*` matches any string at its level.
func Match(a string, b string, spl string, wildcard string) bool {
	var (
		nas []string
		bas []string
	)
	// precheck
	if a[0] != b[0] || (len(a) == 0 || len(b) == 0) {
		return false
	}
	nas, bas = strings.Split(a, spl), strings.Split(b, spl)
	if len(nas) > len(bas) {
		return false
	} else if (len(nas) < len(bas)) && nas[len(nas)-1] != wildcard {
		return false
	}
	for i := 0; i < len(nas); i++ {
		if nas[i] == wildcard {
			if len(bas) > 1 {
				continue
			} else {
				// if i == len(nas) {
				// 	return true
				// }
			}
			return false
		}
		if nas[i] != bas[i] {
			return false
		}
	}
	return true
}

// MatchSplits returns wether two strings seperated by delimiter
// `spl` are exactly a match. First string can contain a wildcard
// `wildcard` which forces the algrithm to count it as a match.
// For example: `/test/*/a` and `/test/a/a` are exact matches
// because `*` matches any string at its level. This function is
// not implemented recursively because of performance issues (
// non recursive version is about 500ms faster than a recursive
// impl.).
func MatchSplits(a string, b string, spl string, wildcard string) bool {
	var (
		na  string
		ba  string
		nas []string
		bas []string
		pf  int
	)
	FindPrefix(&a, &b, &pf)
	na, ba = a[pf:], b[pf:]
	nas, bas = strings.Split(na, spl), strings.Split(ba, spl)
	// if len(nas) != len(bas) {
	// 	return false
	// }
	for i := 0; i < len(nas); i++ {
		if nas[i] == wildcard {
			if len(bas) > 1 {
				continue
			} else {
				// NOTE: NEW:
				// This is when len(a)<len(b)
				// and a's tail is a wildcard.
				if i == len(nas) {
					return true
				}
				// END
				return false
			}
		}
		if nas[i] != bas[i] {
			return false
		}
	}
	return true
}
