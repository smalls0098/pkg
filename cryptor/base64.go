package cryptor

import (
	"encoding/base64"
)

// Base64StdEncodeStr encode string with base64 encoding
func Base64StdEncodeStr(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64StdDecodeStr decode a base64 encoded string
func Base64StdDecodeStr(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

// Base64StdEncode encode string with base64 encoding
func Base64StdEncode(s []byte) []byte {
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(s)))
	enc.Encode(buf, s)
	return buf
}

// Base64StdDecode decode a base64 encoded string
func Base64StdDecode(s []byte) []byte {
	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, s)
	if err != nil {
		return make([]byte, 0)
	}
	return dbuf[:n]
}
