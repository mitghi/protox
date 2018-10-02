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

package crypto

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestRSA(t *testing.T) {
	var text []byte = []byte("a simple test message")
	r, err := MakeRSA(RSADefaultC)
	if err != nil {
		t.Fatal("err!=nil, expected nil", err)
	}
	ct, err := r.Encrypt(text)
	if err != nil {
		t.Fatal("err!=nil, expected nil", err)
	}
	fmt.Println("cipher text: ", strconv.Quote(string(ct)))
	txt, err := r.Decrypt(ct)
	if err != nil {
		t.Fatal("err!=nil, expected nil", err)
	}
	err = r.Save()
	if err != nil {
		t.Fatal("err!=nil, expected nil", err)
	}
	if _, err := os.Stat("./public.key"); err != nil {
		t.Fatal("public.key is not saved on disk")
	}
	if _, err := os.Stat("./private.key"); err != nil {
		t.Fatal("private.key is not saved on disk")
	}
	nr, err := FromFile("./private.key")
	if err != nil {
		t.Fatal("err!=nil, expected nil. Cannot load private file from disk.", err)
	}
	nrpubk, err := nr.MarshalPublicKey()
	if err != nil {
		t.Fatal("err!=nil, expected nil. Cannot get marshalled public key.", err)
	}
	rpubk, err := r.MarshalPublicKey()
	if err != nil {
		t.Fatal("err!=nil, expected nil. Cannot get marshalled public key.", err)
	}
	if c := bytes.Compare(nrpubk, rpubk); c != 0 {
		t.Fatal("nrpubk!=rpubk, expected c==0", c)
	}
	fmt.Println("unciphered text:", strconv.Quote(string(txt)))
	if _, err = LoadPrivateKey("./private.key"); err != nil {
		t.Fatal("err!=nil, expected nil. Cannot load PrivateKey from file.", err)
	}
	if _, err = LoadPublicKey("./public.key"); err != nil {
		t.Fatal("err!=nil, expected nil. Cannot load PublicKey from file.", err)
	}
	signature, err := r.Sign(text)
	if err != nil {
		t.Fatal("err!=nil, expected nil. Cannot sign.", err)
	}
	fmt.Println("signature: ", strconv.Quote(string(signature)))

	err = r.VerifySignature([]byte("a simple test message"), signature)
	if err != nil {
		t.Fatal("err!=nil, expected nil. Cannot verify signature.", err)
	}
	err = r.VerifySignature([]byte("a simple invalid message"), signature)
	if err == nil {
		t.Fatal("err==nil, expected !nil. Verified using wrong message.", err)
	}
	// test invalid length
	_, err = MakeRSA(1023)
	if err == nil {
		t.Fatal("err==nil, expected !nil")
	}
}
