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

// func InsertPath(self *Radix, path []byte, value interface{}, sep string, callback func(**RDXNode)) **RDXNode {
// 	var (
// 		buff  bytes.Buffer
// 		paths [][]byte
// 		plen  int
// 		curr  **RDXNode
// 	)
// 	// TODO
// 	//  check the error code
// 	paths, _ = messages.TopicComponents(path)
// 	plen = len(paths) - 1
// 	for i, v := range paths {
// 		if i == 0 {
// 			buff.Write(v)
// 			a, _ := self.InsertRFrom(self.GetRoot(), buff.String(), value)
// 			curr = &a
// 			buff.Reset()
// 			buff.WriteString(a.Key)
// 			buff.WriteByte('/')
// 			b, _ := self.InsertRFrom(*curr, buff.String(), value)
// 			curr = &b
// 			buff.Reset()
// 			continue
// 		}
// 		buff.WriteByte('/')
// 		buff.Write(v)
// 		a, _ := self.InsertRFrom(*curr, buff.String(), value)
// 		if a != nil {
// 			curr = &a
// 			buff.Reset()
// 			buff.WriteString(a.Key)
// 		} else {
// 			log.Println("(InsertPath)a==nil. TODO: FATAL SITUATION, CHECK.")
// 		}
// 		if buff.Len() > 0 {
// 			lch := buff.Bytes()[buff.Len()-1]
// 			if strings.ContainsAny(string(lch), sep) == true {
// 				n := *curr
// 				if n != nil && n.Key == "*" {
// 					n.SetOpts(0x1)
// 				}
// 			}
// 		}
// 		if i != plen {
// 			buff.WriteByte('/')
// 			b, _ := self.InsertRFrom(*curr, buff.String(), value)
// 			curr = &b
// 		}
// 		buff.Reset()
// 	}
// 	return curr
// }
