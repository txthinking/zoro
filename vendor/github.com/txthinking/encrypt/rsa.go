package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// key is private key, PKCS#8.
func RSASignWithSHA256(data []byte, key []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	b := h.Sum(nil)
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	b, err = rsa.SignPKCS1v15(rand.Reader, k.(*rsa.PrivateKey), crypto.SHA256, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// key is public key, nil means passed.
func RSAVerifyWithSHA256(data, sign, key []byte) error {
	h := sha256.New()
	h.Write(data)
	b := h.Sum(nil)
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("invalid key")
	}
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	if err := rsa.VerifyPKCS1v15(k.(*rsa.PublicKey), crypto.SHA256, b, sign); err != nil {
		return err
	}
	return nil
}

// key is private key, PKCS#1.
func RSASignWithSHA256PKCS1(data []byte, key []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	b := h.Sum(nil)
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid key")
	}
	k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	b, err = rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA256, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
