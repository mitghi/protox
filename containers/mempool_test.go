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
	"fmt"
	"testing"
	"time"
)

func TestBuffPool(t *testing.T) {
	bp := NewBufferPool(time.Duration(time.Millisecond*500), 0, 100)
	_, _, err := bp.Run()
	if err != nil {
		t.Fatal("error is not nil(", err, ")")
	}
	pool := bp.CreateNewContext(5)
	buff := pool.Get()
	buff = buff[0:7]
	tt := []byte("testing")
	copy(buff, tt[:])
	fmt.Println("buff : ", string(buff), len(buff), cap(buff))
	pool.Release(buff)
	for i := 0; i < 28; i++ {
		c := pool.Get()
		c = c[0:7]
		a := pool.Get()
		pool.Get()
		pool.Release(a)
	}
	if rs := pool.RefreshStat(); rs != 1 {
		t.Fatal("invalid stat")
	}
	pool.Flush()
}
