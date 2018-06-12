package util

import (
	"fmt"
	"crypto/sha256"
	"crypto/sha1"
)

func Sha256(s string) string{
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha1(s string) string{
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
