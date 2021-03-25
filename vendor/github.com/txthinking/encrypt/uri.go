package encrypt

import (
	"net/url"
	"strings"
)

func URIEscape(s string) string {
	s = url.QueryEscape(s)
	s = strings.Replace(s, "+", "%20", -1)
	return s
}

func URIUnescape(s string) (string, error) {
	return url.QueryUnescape(s)
}
