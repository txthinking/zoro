package encrypt

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// SHA1 encrypt s according to sha1 algorithm
func SHA1(s string) (r string) {
	var h hash.Hash
	h = sha1.New()
	io.WriteString(h, s)
	r = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA256 encrypt s according to sha256 algorithm
func SHA256(s string) (r string) {
	var h hash.Hash
	h = sha256.New()
	io.WriteString(h, s)
	r = hex.EncodeToString(h.Sum(nil))
	return
}

// SHA256 encrypt s according to sha256 algorithm
func SHA256Bytes(s []byte) ([]byte, error) {
	var h hash.Hash
	h = sha256.New()
	n, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	if n != len(s) {
		return nil, errors.New("Write length error")
	}
	r := h.Sum(nil)
	return r, nil
}
