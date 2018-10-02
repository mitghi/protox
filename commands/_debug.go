package commands

// NOTE
// . this is snippet is from old parseLine implementation.
// for shouldContinue {
// 	shouldContinue = l != 0
// 	value, rem, err := parseString(remaining)
// 	result = append(result, value)
// 	if err != nil {
// 		goto ERROR
// 	}
// 	if rem > 0 {
// 		remaining = remaining[rem:]
// 		l -= rem
// 	} else if rem == 0 {
// 		goto ERROR
// 	}
// 	if remaining[0] == TOK_CR || remaining[0] == TOK_LF || remaining[0] != TOK_LEN {
// 		goto ERROR
// 	}
// 	if rem == -1 {
// 		remaining = remaining[:0]
// 		if l != 0 {
// 			hasRemaining = false
// 		}
// 	}
// }
// NOTE
// . when remaining == 0, it means that all lines in the original
// input are consumed  and there is not any remaining left to be
// parsed.

// parseLine is a function that only parsed a sentence ( textual protocol line ) and returns its values along a boolean `hasRemaining` to indicate whether there is still remaining bytes to be parsed. It returns an error in case of unsuccesfull invokation.
// func parseLine(arg []byte) (result []interface{}, hasRemaining bool, err error) {
// 	var (
// 		// stack     *deque         // stack used to build ast (abstract syntax tree)
// 		shouldContinue bool           = true
// 		l              int            = len(arg) // l is the length of byte slice
// 		remaining      []byte                    // remaining pushes slice pointer into remaining bytes
// 		remindex       int                       // remindex is the remaining index used to push pointer
// 		hinfo          *arrHeaderInfo            // hinfo contains array header info
// 	)
// 	// NOTE
// 	// . by default prove that no bytes left to be parsed unless
// 	//   proved otherwise.
// 	hasRemaining = false
// 	// NOTE
// 	// . check against preconditions for early error detection
// 	if l == 0 || l < 6 {
// 		err = EVIOLEN
// 		goto ERROR
// 	}
// 	hinfo, err = parseArrayHeader(arg)
// 	if err != nil {
// 		goto ERROR
// 	}
// 	// set remindex to array header boundary
// 	remindex = hinfo.end
// 	// NOTE
// 	// . (remindex+1) pushes the pointer forward to pass the array header
// 	// boundary.
// 	if (remindex + 1) > (l - 1) {
// 		goto ERROR
// 	}
// 	// remaining points beyond array header bounadry.
// 	remaining = arg[remindex+1:]
// 	// NOTE
// 	// . in case of an correct index, last character in the byte slice
// 	//   denoted by `hinfo.end` should point to line-feed (TOK_LF),
// 	//   otherwise its invalid.
// 	if arg[hinfo.end] != TOK_LF {
// 		goto ERROR
// 	}
// 	shouldContinue = l != 0
// 	for shouldContinue {
// 		value, rem, err := parseString(remaining)
// 		result = append(result, value)
// 		if err != nil {
// 			goto ERROR
// 		}
// 		switch {
// 		case rem == 0:
// 			// NOTE
// 			// . unknown error, invalid length or format
// 			err = EINLEN
// 			goto ERROR
// 		case rem == -1:
// 			// NOTE
// 			// . flip continue flag to terminate, set remaining flag and
// 			//   zero out length.
// 			shouldContinue = false
// 			if l != 0 {
// 				hasRemaining = false
// 			}
// 			l = 0
// 			break
// 		case rem > 0:
// 			remaining = remaining[rem:]
// 			l -= rem
// 			shouldContinue = l != 0
// 		default:
// 			err = EINVAL
// 			goto ERROR
// 		}
// 		if remaining[0] == TOK_CR || remaining[0] == TOK_LF || remaining[0] != TOK_LEN {
// 			err = EVIOGEN
// 			goto ERROR
// 		}

// 	}
// 	// NOTE
// 	// . when l == 0, it means that all lines in the original
// 	// input are consumed  and there is not any remaining left to be
// 	// parsed.
// 	if l != 0 {
// 		goto ERROR
// 	}

// 	return result, hasRemaining, nil
// ERROR:
// 	return nil, false, err
// }

// - MARK: Rest section.

// TODO
// . finish this function and add a generic parsing capability.
// parse is a function that creates an abstract syntax tree from
// a textual protocol line. It has the capability to parse nested structures from the textual representation.
// func parse(arg []byte) (result interface{}, remaining []byte, etype ParseNodeType, err error) {
// 	var (
// 		b        []byte = arg
// 		l        int    = len(arg)
// 		remindex int
// 		v        byte
// 	)
// 	if l == 0 {
// 		goto ERROR
// 	}
// 	v = b[0]
// 	switch v {
// 	case TOK_ART:
// 		h, err := parseArrayHeader(b)
// 		if err != nil {
// 			goto ERROR
// 		}
// 		remindex = h.end
// 		remaining = arg[remindex:]
// 		res := []interface{}{h}
// 		for i := 0; i < h.es; i++ {
// 			r, rem, err := parseString(remaining)
// 			if err != nil {
// 				goto ERROR
// 			}
// 			res = append(res, r)
// 			if rem == 0 {
// 				break
// 			}
// 		}
// 		result = res
// 		goto OK
// 	case TOK_LEN:
// 		result, remindex, err = parseString(b)
// 		if err != nil {
// 			goto ERROR
// 		}
// 		remaining = arg[remindex:]
// 		goto OK
// 	default:
// 		err = EINVAL
// 		goto ERROR
// 	} // end switch ( v )

// ERROR:
// 	return nil, nil, PNDNone, err
// OK:
// 	return result, remaining, PNDNone, nil
// }

// // tstRecParse is a recursive parser that is used for testing purposes. It will be merged into `parse(....)`, and then removed once the desired threshold for stability and correctness is achieved.
// func tstRecParse(arg []byte) (result []interface{}, err error) {
// 	var (
// 		shouldContinue bool           = true
// 		remaining      []byte         // remaining pushes slice pointer into remaining bytes
// 		remindex       int            // remindex is the remaining index used to push pointer
// 		hinfo          *arrHeaderInfo // hinfo contains array header info
// 	)
// 	// NOTE
// 	// . check against preconditions for early error detection
// 	if len(arg) == 0 || len(arg) < 6 {
// 		err = EVIOLEN
// 		goto ERROR
// 	}
// 	hinfo, err = parseArrayHeader(arg)
// 	if err != nil {
// 		goto ERROR
// 	}
// 	// set remindex to array header boundary
// 	remindex = hinfo.end
// 	// NOTE
// 	// . (remindex+1) pushes the pointer forward to pass the array header
// 	// boundary.
// 	if remindex+1 > len(arg)-1 {
// 		goto ERROR
// 	}
// 	// remaining points beyond array header bounadry.
// 	remaining = arg[remindex+1:]
// 	if arg[hinfo.end] != TOK_LF {
// 		goto ERROR
// 	}
// 	shouldContinue = len(remaining) != 0
// 	for shouldContinue {
// 		value, rem, err := parseString(remaining)
// 		result = append(result, value)
// 		if err != nil {
// 			goto ERROR
// 		}
// 		switch {
// 		case rem == 0:
// 			// NOTE
// 			// . NOP
// 			// shouldContinue = false
// 			err = EINLEN
// 			goto ERROR
// 		case rem == -1:
// 			shouldContinue = false
// 			break
// 		case rem > 0:
// 			remaining = remaining[rem:]
// 			shouldContinue = len(remaining) != 0
// 		default:
// 			err = EINVAL
// 			goto ERROR
// 		}
// 		if remaining[0] == TOK_CR || remaining[0] == TOK_LF || remaining[0] != TOK_LEN {
// 			err = EVIOGEN
// 			goto ERROR
// 		}

// 	}
// 	// NOTE
// 	// . when remaining == 0, it means that all lines in the original
// 	// input are consumed  and there is not any remaining left to be
// 	// parsed.
// 	if len(remaining) != 0 {
// 		goto ERROR
// 	}

// 	return result, nil
// ERROR:
// 	return nil, err
// }

// // parseLine is a function that only parsed a sentence ( textual protocol line ) and returns its values along a boolean `hasRemaining` to indicate whether there is still remaining bytes to be parsed. It returns an error in case of unsuccesfull invokation.
// func parseLine(arg []byte) (result []interface{}, hasRemaining bool, err error) {
// 	var (
// 		// stack     *deque         // stack used to build ast (abstract syntax tree)
// 		l         int            = len(arg) // l is the length of byte slice
// 		remaining []byte                    // remaining pushes slice pointer into remaining bytes
// 		remindex  int                       // remindex is the remaining index used to push pointer
// 		hinfo     *arrHeaderInfo            // hinfo contains array header info
// 	)
// 	// NOTE
// 	// . by default prove that no bytes left to be parsed unless
// 	//   proved otherwise.
// 	hasRemaining = false
// 	// NOTE
// 	// . check against preconditions for early error detection
// 	if l == 0 || l < 6 {
// 		err = EVIOLEN
// 		return nil, false, err
// 	}
// 	hinfo, err = parseArrayHeader(arg)
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	// set remindex to array header boundary
// 	remindex = hinfo.end
// 	// NOTE
// 	// . (remindex+1) pushes the pointer forward to pass the array header
// 	// boundary.
// 	if (remindex + 1) > (l - 1) {
// 		err = EVIOLEN
// 		return nil, false, err
// 	}
// 	// remaining points beyond array header bounadry.
// 	remaining = arg[remindex+1:]
// 	// NOTE
// 	// . in case of an correct index, last character in the byte slice
// 	//   denoted by `hinfo.end` should point to line-feed (TOK_LF),
// 	//   otherwise its invalid.
// 	if arg[hinfo.end] != TOK_LF {
// 		err = EVIOTOK
// 		return nil, false, err
// 	}
// 	value, rem, err := parseString(remaining)
// 	result = append(result, value)
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	switch {
// 	case rem == 0:
// 		// NOTE
// 		// . unknown error, invalid length or format
// 		err = EINLEN
// 		return nil, false, err
// 	case rem == -1:
// 		// NOTE
// 		// . flip continue flag to terminate, set remaining flag and
// 		//   zero out length.
// 		if l != 0 {
// 			hasRemaining = false
// 		}
// 		l = 0
// 		break
// 	case rem > 0:
// 		remaining = remaining[rem:]
// 		l -= rem
// 	default:
// 		err = EINVAL
// 		return nil, false, err
// 	}
// 	if remaining[0] == TOK_CR || remaining[0] == TOK_LF || remaining[0] != TOK_LEN {
// 		err = EVIOGEN
// 		return nil, false, err
// 	}

// 	return result, hasRemaining, nil
// }

// ---------------------------------------------

// func TestNormalParse(t *testing.T) {
// 	v := []byte("*2:4\r\n$1\r\na\r\n$1\r\nb\r\n")
// 	result, remaining, etype, err := parse(v)
// 	fmt.Println(result, remaining, etype, err)
// }

// // TestRecursiveParse is a test suite function that
// // will be removed once recursive parser stabilizes.
// func TestRecursiveParse(t *testing.T) {
// 	var (
// 		arg       []byte = []byte("*2:4\r\n$1\r\na\r\n$1\r\nb\r\n")
// 		remaining []byte
// 		// result    []interface{}
// 		remindex int
// 		es       int // elemnt size
// 		ts       int // total element size ( indivual )
// 	)
// 	hinfo, err := parseArrayHeader(arg)
// 	if err != nil {
// 		t.Fatalf("expected error to be nil, got %+v instead.", err)
// 	}
// 	fmt.Printf("%+v\n", hinfo)
// 	remindex = hinfo.end
// 	if remindex+1 > len(arg)-1 {
// 		t.Fatal("inconsistent state, remaining length exceeds length of the given input.")
// 	}
// 	remaining = arg[remindex+1:]
// 	es = hinfo.es
// 	ts = hinfo.ts
// 	if arg[hinfo.end] != TOK_LF {
// 		t.Fatalf("expected the last argument to be LINE FEED character, got %s instead", string(arg[hinfo.end]))
// 	}
// 	fmt.Printf("remaining (%d), result(%+v), element size (%d), total size (%d)\n", remaining, hinfo, es, ts)

// 	for len(remaining) != 0 {
// 		value, rem, err := parseString(remaining)
// 		if err != nil {
// 			t.Fatalf("expected error to be nil, instead got %+v", err)
// 		}
// 		fmt.Printf("local info (value, remaining bytes index): %s, %d\n", value, rem)
// 		if rem > 0 {
// 			remaining = remaining[rem:]
// 		} else if rem == 0 {

// 			t.Fatal("inconsistent state, unknown condition is reached during parsing remaining bytes length.")
// 		}
// 		if remaining[0] == TOK_CR || remaining[0] == TOK_LF {
// 			t.Fatalf("expected remaining first byte to not be either CARRIAGE RETURN or LINE FEED, instead got value %s", string(remaining[0]))
// 		}
// 		if remaining[0] != '$' {
// 			t.Fatalf("inconsistent state, expected first character in the remaining bytes to be '$', got %s instead.", string(remaining[0]))
// 		}

// 		if rem == -1 {
// 			remaining = remaining[:0]
// 		}
// 	}
// 	if len(remaining) != 0 {
// 		t.Fatalf("inconsistent state, expected length of the remaining bytes in the buffer to be 0, got %d .", len(remaining))
// 	}
// }

// func TestRecursiveParseDummy(t *testing.T) {
// 	var arg []byte = []byte("*2:4\r\n$1\r\na\r\n$1\r\nb\r\n")
// 	_, err := tstRecParse(arg)
// 	if err != nil {
// 		t.Fatalf("expected error==nil, instad got : %+v", err)
// 	}
// }

// func TestParseLine(t *testing.T) {
// 	const _fname string = "TestParseLine"
// 	var arg []byte = []byte("*2:2\r\n$1\r\na\r\n$1\r\nb\r\n")
// 	line, hasRemaining, err := parseLine(arg)
// 	if err != nil {
// 		t.Fatalf("expected err==nil, got %+v.", err)
// 	}
// 	fmt.Printf("[%s] %+v, %t ", _fname, line, hasRemaining)
// 	for _, l := range line {
// 		fmt.Println("line from parseLine function is : ", string(l.([]byte)))
// 	}
// }
