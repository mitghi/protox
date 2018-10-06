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

import (
	"bytes"
)

// Traverse TODO: NEEDS COMMENT INFO
func RTraverse(n *RDXNode, subs *[]string, stack *[]*RDXNode) {
	var (
		stacklen int  = len(*stack)
		wrote    bool = false
	)
	if n == nil {
		return
	}
	goto HEAD

WRITER:
	if stacklen > 0 {
		var (
			buf bytes.Buffer
			s   string
		)
		for _, sn := range *stack {
			buf.WriteString(sn.Key)
		}
		if buf.Len() > 0 {
			s = buf.String()
			(*subs) = append((*subs), s)
		}
	}
	goto MAIN

HEAD:
	(*stack) = append((*stack), n)
	stacklen = len(*stack)
	if n.Isbnd && wrote == false {
		wrote = true
		goto WRITER
	}

MAIN:
	if n.Isbnd && n.Link == nil {
		if stacklen > 1 {
			(*stack) = (*stack)[:stacklen-1]
		} else {
			(*stack) = (*stack)[:0]
		}
		stacklen = len(*stack)
	}
	RTraverse(n.Link, subs, stack)
	RTraverse(n.Next, subs, stack)
	if stacklen >= 1 {
		(*stack) = (*stack)[:stacklen-1]
	}

	return
}
