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

// TODO
// . write comments
// . finish test suits

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	RSADefaultC = 4096
)

var (
	RSALenError    error = errors.New("crypto: invalid rsa len.")
	RSAPow2Error   error = errors.New("crypto: not power of 2.")
	RSAInvalidFile error = errors.New("crypto: no public/private key in file.")
	RSAUnsupported error = errors.New("crypto: unsupported block type.")
)

func makeRSA(c int) (privk *rsa.PrivateKey, pubk *rsa.PublicKey, err error) {
	privk, err = rsa.GenerateKey(rand.Reader, c)
	if err != nil {
		return nil, nil, err
	}
	pubk = &privk.PublicKey
	return privk, pubk, nil
}

func EncryptWith(pubk *rsa.PublicKey, text []byte) ([]byte, error) {
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, pubk, text)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func DecryptWith(privk *rsa.PrivateKey, ct []byte) ([]byte, error) {
	text, err := rsa.DecryptPKCS1v15(rand.Reader, privk, ct)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func SignWith(privk *rsa.PrivateKey, text []byte) ([]byte, error) {
	h := sha256.Sum256(text)
	s, err := rsa.SignPKCS1v15(rand.Reader, privk, crypto.SHA256, h[:])
	if err != nil {
		return nil, err
	}
	return s, nil
}

func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	fc, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(fc)
	if block == nil {
		return nil, RSAInvalidFile
	}
	if block.Type != "RSA PUBLIC KEY" {
		return nil, RSAUnsupported
	}
	ret, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubk, ok := ret.(*rsa.PublicKey)
	if !ok {
		return nil, RSAInvalidFile
	}
	return pubk, nil
}

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	fc, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(fc)
	if block == nil {
		return nil, RSAInvalidFile
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, RSAUnsupported
	}
	privk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privk, nil
}

func FromFile(f string) (*RSA, error) {
	fc, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(fc)
	if block == nil {
		return nil, RSAInvalidFile
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, RSAUnsupported
	}
	privk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubk := &privk.PublicKey
	r := NewRSA(privk, pubk)
	return r, nil
}

func VerifySignatureWith(pubk *rsa.PublicKey, text []byte, sig []byte) error {
	h := sha256.Sum256(text)
	return rsa.VerifyPKCS1v15(pubk, crypto.SHA256, h[:], sig)
}

func NewRSA(privk *rsa.PrivateKey, pubk *rsa.PublicKey) *RSA {
	return &RSA{privk, pubk}
}

func MakeRSA(c int) (ret *RSA, err error) {
	p := (c != 0) && ((c & (c - 1)) == 0)
	if c < 1024 {
		return nil, RSALenError
	} else if !p {
		return nil, RSAPow2Error
	}
	privk, pubk, err := makeRSA(c)
	if err != nil {
		return nil, err
	}
	return &RSA{privk, pubk}, nil
}

func (r *RSA) Encrypt(text []byte) ([]byte, error) {
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, r.pubk, text)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (r *RSA) Decrypt(ct []byte) ([]byte, error) {
	text, err := rsa.DecryptPKCS1v15(rand.Reader, r.privk, ct)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func (r *RSA) Sign(text []byte) ([]byte, error) {
	h := sha256.Sum256(text)
	s, err := rsa.SignPKCS1v15(rand.Reader, r.privk, crypto.SHA256, h[:])
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *RSA) VerifySignature(text []byte, sig []byte) error {
	h := sha256.Sum256(text)
	return rsa.VerifyPKCS1v15(r.pubk, crypto.SHA256, h[:], sig)
}

func (r *RSA) GetPublicKey() *rsa.PublicKey {
	return r.pubk
}

func (r *RSA) GetPrivateKey() *rsa.PrivateKey {
	return r.privk
}

func (r *RSA) MarshalPublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(r.pubk)
}

func (r *RSA) MarhsalPrivateKey() ([]byte, error) {
	PRIVASN1 := x509.MarshalPKCS1PrivateKey(r.privk)
	return PRIVASN1, nil
}

func (r *RSA) Save() error {
	PRIVASN1 := x509.MarshalPKCS1PrivateKey(r.privk)
	PUBASN1, err := x509.MarshalPKIXPublicKey(r.pubk)
	if err != nil {
		return err
	}
	privkey := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: PRIVASN1,
		},
	)
	pubkey := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: PUBASN1,
		},
	)
	err = ioutil.WriteFile("private.key", privkey, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("public.key", pubkey, 0644)
	if err != nil {
		return err
	}

	return nil
}
