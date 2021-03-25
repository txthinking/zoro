package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

// MD5 encrypt s according to md5 algorithm
func MD5(s string) (r string) {
	var h hash.Hash
	h = md5.New()
	io.WriteString(h, s)
	r = hex.EncodeToString(h.Sum(nil))
	return
}
