package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// Sign 签名
func Sign(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write(String2Bytes(str))
	cipherStr := md5Ctx.Sum(nil)
	md5str := hex.EncodeToString(cipherStr)
	return md5str
	//return strings.ToUpper(md5str)
}
