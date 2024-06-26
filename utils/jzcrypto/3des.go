package jzcrypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
)

// TripleDesEncrypt 加密
func TripleDesEncrypt(original, key, iv string) string {
	keyByte := zeroPadding([]byte(key), 24)
	dataEnc := &dataEnc{
		data: []byte(original),
		key:  keyByte,
		iv:   []byte(iv),
	}
	src := tripleDecCbcPkcs7Encrypt(dataEnc)
	data64 := base64.StdEncoding.EncodeToString(src)
	return data64
}

// TripleDesDecrypt 解密
func TripleDesDecrypt(src, key, iv string) (string, error) {
	src64, _ := base64.StdEncoding.DecodeString(src)
	keyByte := zeroPadding([]byte(key), 24)
	dataEnc := &dataEnc{
		data: src64,
		key:  keyByte,
		iv:   []byte(iv),
	}

	var data []byte
	var err error
	func() {
		defer func() {
			if pac := recover(); pac != nil {
				err = errors.New("解析异常")
			}
		}()
		data, err = tripleDesCbcPkcs7Decrypt(dataEnc)
	}()
	return string(data), err
}

type TripleDesConf struct {
	Key string
	Iv  string
}

// TripleDesDecrypt 解密 数组配置
func TripleDesDecryptConf(src string, confList []TripleDesConf) (s string, err error) {
	for _, v := range confList {
		s, err = TripleDesDecrypt(src, v.Key, v.Iv)
		if err != nil {
			s, err = TripleDesDecrypt(src, v.Key, v.Iv)
		} else {
			return s, err
		}
	}
	return s, err
}

type dataEnc struct {
	data []byte
	key  []byte
	iv   []byte
}

// tripleDecCbcPkcs7Encrypt
func tripleDecCbcPkcs7Encrypt(data *dataEnc) []byte {
	//指定算法，返回一个三重DES算法的Block的接口对象
	block, err := des.NewTripleDESCipher(data.key)
	if err != nil {
		panic(err)
	}
	//对明文进行填充
	plainText := padding(data.data, block.BlockSize())

	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, data.iv)
	//对连续的数据块进行加密
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	return cipherText
}

// tripleDesCbcPkcs7Decrypt
func tripleDesCbcPkcs7Decrypt(enc *dataEnc) ([]byte, error) {
	//指定解密算法，返回一个Block接口对象
	block, err := des.NewTripleDESCipher(enc.key)
	if err != nil {
		return nil, err
	}
	iv := enc.iv
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(enc.data))
	blockMode.CryptBlocks(plainText, enc.data)
	//删除填充
	plainText, err = unPaddingNew(plainText)
	return plainText, err
}

// zeroPadding 对key进行填充
func zeroPadding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext) < blockSize {
		padding := blockSize - len(ciphertext)%blockSize
		padText := bytes.Repeat([]byte{0}, padding)
		return append(ciphertext, padText...)
	} else {
		return ciphertext
	}
}

// padding 对明文进行填充
func padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// unPadding  使用PKCS7
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

// unPadding  使用PKCS7 - 函padding 长度检验
func unPaddingNew(input []byte) ([]byte, error) {
	if input == nil || len(input) == 0 {
		return nil, errors.New("invalid padding - empty")
	}

	pc := input[len(input)-1]
	pl := int(pc)
	err := checkPaddingIsValid(input, pl)
	if err != nil {
		return nil, err
	}

	return input[:len(input)-pl], nil
}

// padding 长度检验
func checkPaddingIsValid(input []byte, paddingLength int) error {
	if len(input) < paddingLength {
		return errors.New("invalid padding - 1")
	}

	p := input[len(input)-(paddingLength):]
	for _, pc := range p {
		if uint(pc) != uint(len(p)) {
			return errors.New("invalid padding - 2")
		}
	}

	return nil
}
