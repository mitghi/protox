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
	"golang.org/x/crypto/scrypt"
)

// NewCrypto returns a new struct pointer of type `Crypto`. It uses default settings for n, r, p, keyLen respectively 16384, 8, 1, 32.
func NewCrypto() *Crypto {
	var result *Crypto = &Crypto{
		salt:   []byte{},
		n:      16384,
		r:      8,
		p:      1,
		keyLen: 32,
	}
	return result
}

// SetSalt assigns internal salt value to a given `salt` byte array. It must be called on each new `Crypto` instance.
func (self *Crypto) SetSalt(salt *[]byte) {
	self.salt = make([]byte, len((*salt)))
	copy(self.salt, (*salt))
}

// NewCryptoFromArgs returns a new struct pointer of type `Crypto` and sets the internal according to given arguments. It is not neccessary to use this function directly as the default settings used by `NewCrypto()` is sufficient.
func NewCryptoFromArgs(salt *[]byte, n int, r int, p int, keyLen int) *Crypto {
	var result *Crypto = NewCrypto()
	result.n = n
	result.r = r
	result.p = p
	result.keyLen = keyLen
	result.salt = (*salt)

	return result
}

// Encrypt returns the encrypted data. It may fail, error should be explicitely checked to assure correctness.
func (self *Crypto) Encrypt(input *[]byte) (enc string, err error) {
	buff, err := scrypt.Key((*input), self.salt, self.n, self.r, self.p, self.keyLen)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

// Decrypt returns the decrypted data. Simlilar to `Encrypt`, error must be checked explicitely.
func (self *Crypto) Decrypt(input *[]byte) (dec string, err error) {
	return self.Encrypt(input)
}
