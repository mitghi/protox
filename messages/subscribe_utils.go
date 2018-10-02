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

import "fmt"

// debugPrint is a wrapper for `dprint`.
func (self *Subscribe) debugPrint() {
	self.dprint("", 0)
}

// dprint traverses the chain and print its data recursively.
func (self *Subscribe) dprint(lmsg string, ln int) {
	if self == nil {
		return
	}
	fmt.Printf("topic:|\x1b[36m%-40s\033[0m %-v|\n", lmsg, self.topic)
	fmt.Printf("subs:|\x1b[36m%-40s\033[0m %-v|\n", lmsg, self.subs)
	fmt.Printf("qos :|\x1b[36m%-40s\033[0m %-v|\n", lmsg, self.qos)
	fmt.Printf("next:|\x1b[36m%-40s\033[0m %-v|\n", lmsg, self.next)
	fmt.Println("-------------------------------------------------")
	for k, _ := range self.next {
		v := self.next[k]
		v.dprint(lmsg+"....", ln+1)
	}
}
