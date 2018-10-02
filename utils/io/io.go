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

package io

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

// FileExists returns wether the path/file exists. It returns an error
// if corresponding file/path is non-existent.
func FileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

// GetPipedStdin reads from standard input if data is piped ( using pipe `|` operator ).
// It returns an error and an empty string otherwise.
func GetPipedStdin() (string, error) {
	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		str := string(bytes)
		return str, nil
	}
	return "", errors.New("no piped stdin")
}

// GetTermStdin reads from standard input if data is redirected ( using io directors ).
// It returns an erro and an empty string otherwise.
func GetTermStdin() (string, error) {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return input, nil
	}
	return "", errors.New("no stdin from term")
}

// SelectStdin is a convenience routine and checks both `GetPipedStdin` and `GetTermStdin`
// for data. It returns first of two with available data. Error must be checked explicitely.
func SelectStdin() (string, error) {
	pipeout, err := GetPipedStdin()
	if err == nil {
		return pipeout, nil
	}
	termout, err := GetTermStdin()
	if err == nil {
		return termout, nil
	}
	return "", errors.New("no data in stdin")
}

// Exec executes a shell commands and returns the output.
func Exec(name string, args ...string) (output []byte, err error) {
	return exec.Command(name, args...).Output()
}
