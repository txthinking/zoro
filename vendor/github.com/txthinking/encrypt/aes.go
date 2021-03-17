package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AES256KeyLength is the length of key for AES 256 crypt
const AES256KeyLength = 32

// AESMake256Key cut or append empty data on the key
// and make sure the key lenth equal 32
func AESMake256Key(k []byte) []byte {
	if len(k) < AES256KeyLength {
		a := make([]byte, AES256KeyLength-len(k))
		return append(k, a...)
	}
	if len(k) > AES256KeyLength {
		return k[:AES256KeyLength]
	}
	return k
}

// AESCFBEncrypt encrypt s with given k.
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits.
// First 16 bytes of cipher data is the IV.
func AESCFBEncrypt(s, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cb[aes.BlockSize:], s)
	return cb, nil
}

// AESDecrypt decrypt c with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func AESCFBDecrypt(c, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	if len(c) < aes.BlockSize {
		err := errors.New("crypt data is too short")
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cb, cb)
	return cb, nil
}

// AESCBCEncrypt encrypt s with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func AESCBCEncrypt(s, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	if len(s)%aes.BlockSize != 0 {
		return nil, errors.New("invalid length of s")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	cb := make([]byte, aes.BlockSize+len(s))
	iv := cb[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cb[aes.BlockSize:], s)
	return cb, nil
}

// AESCBCDecrypt decrypt c with given k
// k should be 128/256 bits, otherwise it will append empty data or cut until 256 bits
// First 16 bytes of cipher data is the IV.
func AESCBCDecrypt(c, k []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	if len(c) < aes.BlockSize {
		return nil, errors.New("c too short")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	iv := c[:aes.BlockSize]
	cb := c[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cb, cb)
	return cb, nil
}

// PKCS5Padding can append data of PKCS5
// Common blockSize is aes.BlockSize
func PKCS5Padding(c []byte, blockSize int) []byte {
	pl := blockSize - len(c)%blockSize
	p := bytes.Repeat([]byte{byte(pl)}, pl)
	return append(c, p...)
}

// PKCS5UnPadding can unappend data of PKCS5
func PKCS5UnPadding(s []byte) ([]byte, error) {
	l := len(s)
	if l == 0 {
		return nil, errors.New("s too short")
	}
	pl := int(s[l-1])
	if l < pl {
		return nil, errors.New("s too short")
	}
	return s[:(l - pl)], nil
}

// AESGCMEncrypt encrypt s use k and nonce
func AESGCMEncrypt(s, k, n []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	c := g.Seal(nil, n, s, nil)
	return c, nil
}

// AESGCMDecrypt decrypt s use k and nonce
func AESGCMDecrypt(c, k, n []byte) ([]byte, error) {
	if len(k) != 16 && len(k) != 32 {
		k = AESMake256Key(k)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	s, err := g.Open(nil, n, c, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}
