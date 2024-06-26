package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// AesEncrypt AES-CBC加密, 参数及返回值皆使用字符串
// result[0] base64编码后的密文
// result[1] 错误信息, 无错误返回 ok
func AesEncryptStr(content, key string) (string, string) {
	ret, err := AesEncrypt([]byte(content), []byte(padKey(key)))
	if err != nil {
		return "", err.Error()
	}
	return Base64Encode(ret), "ok"
}

// AesEncrypt AES-CBC解密, 参数及返回值皆使用字符串
// param: ciphertext 需base64编码
// result[0] 解密结果
// result[1] 错误信息, 无错误返回 ok
func AesDecryptStr(ciphertext, key string) (string, string) {
	ret, err := AesDecrypt(Base64Decode(ciphertext), []byte(padKey(key)))
	if err != nil {
		return "", err.Error()
	}
	return string(ret), "ok"
}

// AesEncrypt AES-CBC加密
func AesEncrypt(content, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	content = PKCS7Padding(content, blockSize)
	// 向量 (key[:blockSize]) 是密钥的前 blockSize (16) 个字节
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	ciphertext := make([]byte, len(content))
	blockMode.CryptBlocks(ciphertext, content)
	return ciphertext, nil
}

// AesEncrypt AES-CBC解密
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	oriData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(oriData, ciphertext)
	oriData = PKCS7UnPadding(oriData)
	return oriData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(oriData []byte) []byte {
	length := len(oriData)
	unpadding := int(oriData[length-1])
	if length < unpadding {
		return []byte{}
	}
	return oriData[:(length - unpadding)]
}

func padKey(key string) string {
	kl := len(key)
	if kl < 16 {
		key = fmt.Sprintf("%016s", key)
	} else if 16 < kl && kl < 24 {
		key = fmt.Sprintf("%024s", key)
	} else if 24 < kl && kl < 32 {
		key = fmt.Sprintf("%032s", key)
	} else if 32 < kl {
		key = Substr(key, 0, 32)
	}
	return key
}
