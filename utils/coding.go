package utils

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"log"
	"time"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"hash/crc32"
	crand "crypto/rand"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/md5"
)

// md5
func Md5Sum(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// FileMD5
func FileMD5(filePath string) (string, error) {
	var ret string
	file, err := os.Open(filePath)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ret, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	ret = hex.EncodeToString(hashInBytes)
	return ret, nil
}

// 加解密函数 根据dz的Authcode改写的go版本
// params[0] 加密or解密 bool true：加密 false：解密 默认false
// params[1] 秘钥
// params[2] 加密：过期时间
// params[3] 动态秘钥长度 默认：4位 不能大于32位
func Authcode(text string, params ...interface{}) string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("authcode error:%#v", err)
		}
	}()

	l := len(params)

	isEncode := false
	key := "DH-Framework"
	expiry := 0
	cKeyLen := 10

	if l > 0 {
		isEncode = params[0].(bool)
	}

	if l > 1 {
		key = params[1].(string)
		if key=="" {
			key = "DH-Framework"
		}
	}

	if l > 2 {
		expiry = params[2].(int)
		if expiry < 0 {
			expiry = 0
		}
	}

	if l > 3 {
		cKeyLen = params[3].(int)
		if cKeyLen < 0 {
			cKeyLen = 0
		}
	}
	if cKeyLen > 32 {
		cKeyLen = 32
	}

	timestamp := time.Now().Unix()

	// md5加密key
	mKey := Md5Sum(key)

	// 参与加密的
	keyA := Md5Sum(mKey[0:16])
	// 用于验证数据有效性的
	keyB := Md5Sum(mKey[16:])
	// 动态部分
	var keyC string
	if cKeyLen > 0 {
		if isEncode {
			// 加密的时候，动态获取一个秘钥
			keyC = Md5Sum(fmt.Sprint(timestamp))[32 - cKeyLen:]
		} else {
			// 解密的时候从头部获取动态秘钥部分
			keyC = text[0:cKeyLen]
		}
	}

	// 加入了动态的秘钥
	cryptKey := keyA + Md5Sum(keyA + keyC)
	// 秘钥长度
	keyLen := len(cryptKey)
	if isEncode {
		// 加密 前10位是过期验证字符串 10-26位字符串验证
		var d int64
		if expiry > 0 {
			d = timestamp + int64(expiry)
		}
		text = fmt.Sprintf("%010d%s%s", d, Md5Sum(text + keyB)[0:16], text)
	} else {
		// 解密
		text = string(Base64Decode(text[cKeyLen:]))
	}

	// 字符串长度
	textLen := len(text)
	if textLen <= 0 {
		panic(fmt.Sprintf("auth[%s]textLen<=0", text))
	}

	// 密匙簿
	box := Range(0, 256)
	// 对称算法
	var rndKey []int
	cryptKeyB := []byte(cryptKey)
	for i := 0; i < 256; i++ {
		pos := i % keyLen
		rndKey = append(rndKey, int(cryptKeyB[pos]))
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndKey[i]) % 256
		box[i], box[j] = box[j], box[i]
	}

	textB := []byte(text)
	a := 0
	j = 0
	var result []byte
	for i := 0; i < textLen; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		box[a], box[j] = box[j], box[a]
		result = append(result, byte(int(textB[i]) ^ (box[(box[a] + box[j]) % 256])))
	}

	if isEncode {
		return keyC + strings.Replace(Base64Encode(result), "=", "", -1)
	}

	// 获取前10位，判断过期时间
	d := Atoi64(string(result[0:10]), 0)
	if (d == 0 || d - timestamp > 0) && string(result[10:26]) == Md5Sum(string(result[26:]) + keyB)[0:16] {
		return string(result[26:])
	}

	panic(fmt.Sprintf("auth[%s]", text))

	return ""
}

// AuthcodeUrl 处理Authcode函数的加密解密结果以便url传输
func AuthcodeUrl(text string, params ...interface{}) string {
	isEncode := false
	if len(params) > 0 {
		isEncode = params[0].(bool)
	}
	if isEncode { //加密
		return strings.Replace(strings.Replace(Authcode(text, params...), "+", ",", -1), "/", "-", -1)
	} else {
		return Authcode(strings.Replace(strings.Replace(text, ",", "+", -1), "-", "/", -1), params...)
	}
}

// JsonEncodeWithError 编码JSON
func JsonEncodeWithError(m interface{}) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}


// JsonDecodeWithError 解码JSON
func JsonDecodeWithError(str string, v ...interface{}) (interface{}, error) {
	var m interface{}
	if len(v) > 0 {
		m = v[0]
	} else {
		m = make(map[string]interface{})
	}

	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// JsonEncode 编码JSON
func JsonEncode(m interface{}) string {
	s, err := JsonEncodeWithError(m)
	if err != nil {
		log.Printf("Json Encode[%#v] Error:%s", m, err.Error())
		return ""
	}
	return s
}

// JsonDecode 解码JSON
func JsonDecode(str string, v ...interface{}) interface{} {
	i, err := JsonDecodeWithError(str, v...)
	if err != nil {
		log.Printf("Json Decode[%s] Error:%s", str, err.Error())
		return nil
	}
	return i
}

func Crc32(text string) string {
	h := crc32.NewIEEE()
	io.WriteString(h, text)
	return fmt.Sprintf("%d", h.Sum32())
}

// RsaEncode rsa加密
func RsaEncode(b, rsaKey []byte) ([]byte, error) {
	block, _ := pem.Decode(rsaKey)
	if block == nil {
		return b, errors.New("key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return b, err
	}
	return rsa.EncryptPKCS1v15(crand.Reader, pub.(*rsa.PublicKey), b)
}

// RsaDecode rsa解密
func RsaDecode(b, rsaKey []byte) ([]byte, error) {
	block, _ := pem.Decode(rsaKey)
	if block == nil {
		return b, errors.New("key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return b, err
	}
	return rsa.DecryptPKCS1v15(crand.Reader, priv, b)
}

func HashHmac(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func HashHmacRaw(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return fmt.Sprintf("%s", mac.Sum(nil))
}

// Base64Encode Base64编码
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Base64Decode Base64解码
func Base64Decode(str string) []byte {
	var b []byte
	var err error
	x := len(str) * 3 % 4
	switch {
	case x == 2:
		str += "=="
	case x == 1:
		str += "="
	}
	if b, err = base64.StdEncoding.DecodeString(str); err != nil {
		return b
	}

	return b
}

// UrlEncode 编码
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// UrlDecode 解码
func UrlDecode(str string) string {
	ret, _ := url.QueryUnescape(str)
	return ret
}

/*func Urlencode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func Urldecode(str string) string {
	b, e := base64.URLEncoding.DecodeString(str)
	if e != nil {
		log.Printf("urldecode error:%s", e.Error())
		return ""
	}

	return string(b)
}*/
