package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(raw string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(raw))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
