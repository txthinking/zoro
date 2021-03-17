package encrypt

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
)

func HmacSha256(message, key []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(message); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func CheckHmacSha256(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(message); err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC), nil
}

func HmacSha1(message, key []byte) ([]byte, error) {
	mac := hmac.New(sha1.New, key)
	if _, err := mac.Write(message); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func CheckHmacSha1(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha1.New, key)
	if _, err := mac.Write(message); err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC), nil
}
