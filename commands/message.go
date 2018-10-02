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

package commands

import (
	"bytes"
	"fmt"
	"strconv"
)

// addString is a function that takes a single argument of type `string`
// and creates an string representation.
func addString(arg string) string {
	return fmt.Sprintf(STRFormat, len(arg), arg)
}

// addArray is a function that takes variadic arguments of type
// `string` and creates an array and returns an string represntation.
func addArray(args ...string) string {
	var (
		tl   int = 0
		buff bytes.Buffer
	)
	for _, v := range args {
		tl += len(v)
		buff.WriteString(addString(v))
	}
	return fmt.Sprintf(ARRHFormatApnd, len(args), tl, buff.String())
}

// addArrayHeader is a function that takes two arguments,
// number of elements `nval` and total length `tl` respectively
// and returns an array header.
func addArrayHeader(nval int, tl int) string {
	return fmt.Sprintf(ARRHFormat, nval, tl)
}

// parseArrayHeader is a function that takes an slice starting from
// astrix sign up until CRLF ( e.g. "*4:18\r\n" ) and parses
// it values. It returns an error to indicate a failure. In case of
// faulty inputs, it tries to recover from a panic.
func parseArrayHeader(arg []byte) (hinfo *arrHeaderInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	var (
		l        int    = len(arg)
		state    Token  = OP_NONE
		eeindex  int    = -1 // eeindex is the element end index
		tseindex int    = -1 // tseindex is totalsize end index
		sepindex int    = -1 // sepindex is the index of ':' seperator
		el       []byte      // first field ( before : )
		tl       []byte      // second field ( after : )
		v        byte
	)
	// TODO
	// . prevent memory allocations in this function
	// . add global state (preallocated) ?
	err = EINVARRH
	if l == 0 {
		return nil, err
	} else if arg[0] != TOK_ART {
		return nil, EVIOTOK
	}
	hinfo = &arrHeaderInfo{}
ML:
	for i := 0; i < l; i++ {
		v = arg[i]
		switch v {
		case TOK_ART:
			if state != OP_NONE {
				err = EINVARRH
				break ML
			}
			state = OP_ARR_ES
			break
		case TOK_SEP:
			if state != OP_ARR_ES {
				err = EVIOTOK
				break ML
			}
			state = OP_ARR_TS
			sepindex = i
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch state {
			case OP_ARR_ES:
				eeindex = i
				break
			case OP_ARR_TS:
				tseindex = i
				break
			default:
				err = EVIOTOK
				break ML
			}
			break
		case TOK_CR:
			if i-1 <= 0 || arg[i+1] != TOK_LF {
				err = EINVARRH
				break ML
			}
			hinfo.end = i + 1
			if (sepindex == -1) || (eeindex == -1) || (tseindex == -1) {
				err = EINVARRH
				break ML
			}
			el = arg[1:sepindex]              // ( & -> : )
			tl = arg[sepindex+1 : tseindex+1] // [ : -> \r )
			if len(el) == 0 || len(tl) == 0 {
				err = EINVARRH
				break ML
			}
			hinfo.es, err = strconv.Atoi(string(el))
			if err != nil {
				break ML
			}
			hinfo.ts, err = strconv.Atoi(string(tl))
			if err != nil {
				break ML
			}
			// TODO
			// . remove appending to byte slice
			hinfo.c = append(hinfo.c, el)
			hinfo.c = append(hinfo.c, tl)
			goto OK
		}
	}

	return nil, err
OK:
	return hinfo, nil
}

// parseString is a function that parses an string from
// input. It returns an error and empty slice in case of
// failure as well as an int which refers to beginning of
// next line or returns -1 in when it ceased to exist. It
// is important to note that this function only parses
// very first serialized line.
func parseString(arg []byte) (value []byte, remaining int, err error) {
	// NOTE
	// . remaining bytes assumed to be non existing unless
	//   proved otherwise.
	remaining = -1
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	var (
		v      byte
		state  Token = OP_NONE
		l      int   = len(arg)
		ls     int   = -1 // ls is start index of length field
		le     int   = -1 // le is end index of length field
		rindex int   = -1 // index is position of first '\r' occurance
		pl     int   = -1 // pl is payload length ( converted to integer )
		ps     int   = -1 // ps is payload start position ( rindex+2 )
		pe     int   = -1 // pe is payload end position ( ps + pl )
	)
	// NOTE
	// . inputs with length < 6 are invalid because least feasible
	//   value ( string with 0 length ) is represented as "$0\r\n\r\n",
	//   thus length < 6 is a invalid input.
	if l == 0 {
		err = EINVAL
		goto ERROR
	} else if l < 6 {
		err = EINLEN
		goto ERROR
	}
	switch arg[0] {
	case TOK_ART:
		break
	case TOK_LEN:
		break
	default:
		err = EVIOTOK
		goto ERROR
	}
	for i := 0; i < l; i++ {
		v = arg[i]
		switch v {
		case TOK_LEN:
			if state != OP_NONE {
				err = fmt.Errorf(cErrf, "invalid state.")
				goto ERROR
			}
			state = OP_LEN
			ls = i + 1
			break
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if state != OP_LEN {
				err = fmt.Errorf(cErrf, "invalid state.")
				goto ERROR
			}
			if (i+1 < l-1) && (arg[i+1] == TOK_CR) {
				le = i + 1
			}
			break
		case TOK_CR:
			switch state {
			case OP_LEN:
				if le == -1 {
					err = fmt.Errorf(cErrf, "invalid length end index")
					goto ERROR
				}
				rindex = i
				ps = rindex + 2
				pl, err = strconv.Atoi(string(arg[ls:le]))
				if err != nil {
					goto ERROR
				}
				pe = ps + pl
				if pe >= l {
					err = fmt.Errorf(cErrf, "length is out of bound.")
					goto ERROR
				} else if arg[pe] != TOK_CR {
					err = EVIOLEN
					goto ERROR
				}
				value = arg[ps:pe]
				if (pe + 2) < (l - 1) {
					// NOTE:
					// . adding 2 by passes line feed character
					//   and pushes the pointer to beggining of
					//   next line.
					remaining = pe + 2
				}
				goto OK
			default:
				err = EINVAL
				goto ERROR
			} // end switch ( state )
		default:
			err = EINVAL
			goto ERROR
		} // end switch ( v )
	} // end for

ERROR:
	return nil, remaining, err

OK:
	return value, remaining, nil
}

// parse is a function that creates an abstract syntax tree from
// a textual protocol line. It has the capability to parse nested structures from the textual representation.
func parseArray(arg *[]byte) (result interface{}, etype ParseNodeType, err error) {
	var (
		b        []byte = *arg
		current  interface{}
		l        int = len((*arg))
		remindex int
		v        byte
	)
	if l == 0 {
		err = EINLEN
		goto ERROR
	}
	v = b[0]
	switch v {
	case TOK_ART:
		h, err := parseArrayHeader(b)
		if err != nil {
			goto ERROR
		}
		etype = PNDArrayLit
		current = &arrinfo{}
		remindex = h.end
		(*arg) = (*arg)[remindex+1:]
		current.(*arrinfo).expected = h.es
		var counter int = 0
		for i := 0; i < h.es; i++ {
			r, _, err := parseArray(arg)
			if err != nil {
				goto ERROR
			}
			if len((*arg)) == 0 {
				break
			}
			current.(*arrinfo).values = append(current.(*arrinfo).values, r)
			counter += 1
		}
		if counter != h.es {
			err = EVIOLEN
			goto ERROR
		}
		current.(*arrinfo).curr = counter
		result = current
		goto OK
	case TOK_LEN:
		res, remindex, err := parseString(b)
		if err != nil {
			goto ERROR
		}
		result = string(res)
		if remindex < len((*arg))-1 && remindex >= 0 {
			(*arg) = (*arg)[remindex:]
		}
		goto OK
	default:
		err = EINVAL
		goto ERROR
	} // end switch ( v )

ERROR:
	return nil, PNDNone, err
OK:
	return result, etype, nil
}

// Parse is a wrapper function around `parseArray`.
func Parse(arg *[]byte) (result interface{}, err error) {
	if arg == nil {
		err = EINVAL
		goto ERROR
	} else if len(*arg) == 0 {
		err = EINVAL
		goto ERROR
	}
	switch (*arg)[0] {
	case TOK_ART:
		break
	case TOK_LEN:
		break
	default:
		err = EVIOTOK
	}
	result, _, err = parseArray(arg)
	if err != nil {
		goto ERROR
	}

	return result, nil
ERROR:
	return nil, err
}
