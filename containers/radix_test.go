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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

// RALLC is preallocation rate
const RALLC = 1000

var words []string = []string{"romane", "romanus", "romulus", "rubens", "ruber", "rubicon", "rubicundus"}
var words2 []string = []string{"monochroic", "monochroics", "monochromasies", "monochromasy", "monochromat", "monochromate", "monochromates", "monochromatic", "monochromatically", "monochromaticities", "monochromaticity", "monochromatics", "monochromatism", "monochromatisms", "monochromator", "monochromators", "monochromats", "monochrome", "monochromes", "monochromic", "monochromical", "monochromies", "monochromist", "monochromists"}
var wordSlice []string

func init() {
	data, err := ioutil.ReadFile("./misc/list.json")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	err = json.Unmarshal(data, &wordSlice)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Printf("+ parsed list.json ( length: %d ).\n", len(wordSlice))
}

func TestRadixNormal(t *testing.T) {
	radix := NewRadix(RALLC)
	for index, word := range words {
		radix.Insert(word, index)
	}
	radix.Print()
}

func TestRadix(t *testing.T) {
	radix := NewRadix(RALLC)
	radix.Insert(words[0], 0)
	if radix.root.Next == nil {
		t.Fatalf("Invalid radix tree, expected %d edges in root.", 1)
	}

	if radix.root.Next.Key != words[0] {
		t.Fatalf("Invalid radix tree, expected root edge to be %s, got %s.", words[0], radix.root.Key)
	}

	radix.Insert(words[1], 0)
	if radix.root.Next.Key != "roman" {
		t.Fatal("Invalid head")
	}
	for _, word := range words[2:] {
		radix.Insert(word, 0)
	}
	if radix.root.Next.Key != "r" {
		t.Fatal("Invalid head")
	}
	for edge := radix.root.Next.Link; edge != nil; edge = edge.Next {
		switch edge.Key {
		case "ub":
			continue
		case "om":
			continue
		default:
			t.Fatalf("Invalid children after traversing head ( %s )", edge.Key)
		}
	}
	radix.Insert("new", 1)
	radix.Insert("newanother", 1)
	radix.Insert("newanoather", 1)
	radix.Insert("romanee", 1)

	radix.Print()
}

func TestFindStringPath(t *testing.T) {
	radix := NewRadix(RALLC)
	errmsg := "Invalid string, expected %s, got %s."
	for _, word := range words[:3] {
		radix.Insert(word, 1)
	}
	word := words[2]
	_, _, joinedResult := radix.Find(word)
	if joinedResult != word {
		t.Fatalf(errmsg, word, joinedResult)
	}
	word = "romanee"
	radix.Insert(word, 1)
	_, _, joinedResult = radix.Find(word)
	if joinedResult != word {
		t.Fatalf(errmsg, word, joinedResult)
	}
	word = "romaneey"
	_, _, joinedResult = radix.Find(word)
	if joinedResult == word {
		t.Fatalf("Found string that was not inserted into radix ( %s ) ", joinedResult)
	}
	word = "romay"
	_, _, joinedResult = radix.Find(word)
	if joinedResult == word {
		t.Fatalf("Found string that was not inserted into radix ( %s )", joinedResult)
	}
}

func TestBoundaries(t *testing.T) {
	radix := NewRadix(RALLC)
	for index, word := range words {
		radix.Insert(word, index)
	}
	radix.Print()

	var truthTable = []struct {
		Key            string
		ParentBoundary bool
		ChildBoundary  bool
	}{
		{"romane", false, true},
		{"romanus", false, true},
		{"rubic", false, false},
		{"rubens", false, true},
		{"ruber", false, true},
	}

	for _, expect := range truthTable {
		_, tail, joinedResult := radix.Find(expect.Key)
		if joinedResult != expect.Key {
			t.Fatalf("Incorrect result from radix, expected %s got %s", expect.Key, joinedResult)
		}
		expectedChildBound := expect.ChildBoundary
		if expectedChildBound != tail.Isbnd {
			t.Fatal("Inconsistent boundaries in radix")
		}
	}
	radix.Insert("roman", 1)
	radix.Insert("roan", 1)
	radix.Print()
	// TODO: check boundaries
}

func TestRemove(t *testing.T) {
	radix := NewRadix(RALLC)
	// TODO
	//  the commented lines below needs to be verified as the underlaying structure has changed
	//
	// if status := radix.Remove(""); status != false {
	//   t.Fatal("Removed empty radix")
	// }
	if status := radix.Remove("test"); status != false {
		t.Fatal("Removed non existing string")
	}
	for _, word := range words[:3] {
		radix.Insert(word, 1)
	}

	toremove := words[2]
	status := radix.Remove(toremove)
	if status != true {
		t.Fatalf("Cannot remove %s from radix. %+v", toremove, radix.root.Link)
	}
	for _, word := range words[3:] {
		radix.Insert(word, 1)
	}

	status = radix.Remove("roman")
	if status != true {
		t.Fatal("Cannot remove roman from radix")
	}

	radix.Insert("royal", 1)
	status = radix.Remove("r")

	for _, word := range words {
		radix.Insert(word, 1)
	}

	status = radix.Remove("embrace")
	if status == true {
		t.Error("Tried to remove non existing string")
	}

	radix.Print()
}

func BenchmarkFindInRadix270000(t *testing.B) {
	radix := NewRadix(RALLC)
	for _, word := range wordSlice {
		radix.Insert(word, 0)
	}
	t.ResetTimer()

	samples := [...]string{"monochromate", "mistranscriptions"}

	for _, sample := range samples {
		_, _, joinedResult := radix.Find(sample)
		if joinedResult != sample {
			t.Fatalf("Invalid string in benchmark, expected %s(len: %d) got %s(len: %d)", sample, len(joinedResult), sample, len(sample))
		}
	}
}

func BenchmarkInsertIntoRadix270000(t *testing.B) {
	length := 270000
	radix := NewRadix(RALLC)

	for i := 0; i < length; i++ {
		radix.Insert(wordSlice[i], i)
	}
}

func BenchmarkInsertAndFindInRadix270000(t *testing.B) {
	length := 270000
	radix := NewRadix(RALLC)

	for i := 0; i < length; i++ {
		radix.Insert(wordSlice[i], i)
	}

	var inp string = "monochromate"
	head, tail, out := radix.Find(inp)
	fmt.Printf("Node(%+v)<head: %p, tail: %p>, Input(\"%s\"), Output(\"%s\")\n", head, head, tail, inp, out)
}

func TestRadix2(t *testing.T) {
	radix := NewRadix(RALLC)

	for index, word := range words2 {
		radix.Insert(word, index)
	}

	radix.Print()
}
