package util

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// SignSHA256WithKey SHA-256签名
func SignSHA256WithKey(str, key string) string {
	strSign := str + key
	// 计算SHA-256哈希值
	hash := sha256.Sum256([]byte(strSign))
	// 哈希转16进制小写字符串
	signValue := hex.EncodeToString(hash[:])
	return signValue
}

func SignMD5Decode(data []byte, key string) (string, error) {
	keyByte, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write(data)
	hash.Write(keyByte)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func EcbEncrypt(data []byte, key string) ([]byte, error) {
	// base64字符串转Byte
	keyByte, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	// 填充为block的整数倍
	data = PKCS7Padding(data, block.BlockSize())
	// 按块block加密
	ciphertext := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(ciphertext[bs:be], data[bs:be])
	}
	return ciphertext, nil
}

func EcbDecrypt(data []byte, key string) ([]byte, error) {
	// base64字符串转Byte
	keyByte, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	// 按块block解密
	ciphertext := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(ciphertext[bs:be], data[bs:be])
	}
	return PKCS7UnPadding(ciphertext), nil
}
