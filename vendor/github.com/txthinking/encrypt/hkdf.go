package encrypt

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

func HkdfSha256RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error) {
	hash := sha256.New

	salt = make([]byte, sl)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return
	}

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

func HkdfSha256WithSalt(secret, salt, info []byte) (key []byte, err error) {
	hash := sha256.New

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

func HkdfSha1RandomSalt(secret, info []byte, sl int) (key []byte, salt []byte, err error) {
	hash := sha1.New

	salt = make([]byte, sl)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return
	}

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}

func HkdfSha1WithSalt(secret, salt, info []byte) (key []byte, err error) {
	hash := sha1.New

	hkdf := hkdf.New(hash, secret, salt, info)

	key = make([]byte, hash().Size())
	if _, err = io.ReadFull(hkdf, key); err != nil {
		return
	}
	return
}
