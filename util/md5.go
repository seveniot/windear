package util

import (
	"encoding/hex"
	"crypto/md5"
)

func MD5(raw string)string{
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(raw))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
