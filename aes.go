/**
 * @Author: DPY
 * @Description: AES 高级加密标准（Advanced Encryption Standard，AES）
 * @File:  sym_aes
 * @Version: 1.0.0
 * @Date: 2021/11/19 14:57
 */

package wecom

/*
- 这个标准用来替代原先的DES（Data Encryption Standard）
- AES的区块长度固定为128位，密钥长度则可以是128，192或256位
- AES的处理单位是字节
- AES为分组密码，分组密码也就是把明文分成一组一组的，每组长度相等，每次加密一组数据，直到加密完整个明文。
- 在AES标准规范中，分组长度只能是128位，也就是说，每个分组为16个字节（每个字节8位）。
- 密钥的长度可以使用128位、192位或256位。密钥的长度不同，推荐加密轮数也不同
*/

// 原理介绍可参考 https://blog.csdn.net/qq_28205153/article/details/55798628

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func AesEncryptCBC(msg, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	msg = pkcs7Padding(msg, block.BlockSize())
	crypt, err := cbcEncrypt(block, msg, iv...)
	if err != nil {
		return nil, err
	}
	return crypt, nil
}

func AesDecryptCBC(crypt, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypt, err := cbcDecrypt(block, crypt, iv...)
	if err != nil {
		return nil, err
	}
	return decrypt, nil
}

func AesEncryptCFB(msg, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	msg = pkcs7Padding(msg, block.BlockSize())
	crypt, err := cfbEncrypt(block, msg, iv...)
	if err != nil {
		return nil, err
	}
	return crypt, nil
}

func AesDecryptCFB(crypt, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypt, err := cfbDecrypt(block, crypt, iv...)
	if err != nil {
		return nil, err
	}
	decrypt = pkcs57Trimming(decrypt)
	return decrypt, nil
}

func AesEncryptGCM(msg, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	msg = pkcs7Padding(msg, block.BlockSize())
	crypt, err := gcmEncrypt(block, msg, nonce)
	if err != nil {
		return nil, err
	}
	return crypt, nil
}

func AesDecryptGCM(crypt, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypt, err := gcmDecrypt(block, crypt, nonce)
	if err != nil {
		return nil, err
	}
	decrypt = pkcs57Trimming(decrypt)
	return decrypt, nil
}

var quickKey = []byte("daipengyuan-1987")

// QuickEncrypt 使用aes-gcm与内置默认密钥快速加密
func QuickEncrypt(msg string) string {
	salt := []byte(RandString(12))
	crypted, err := AesEncryptGCM([]byte(msg), quickKey, salt)
	if err != nil {
		return ""
	}
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	crypted64Len := enc.EncodedLen(len(crypted))
	crypted64 := make([]byte, crypted64Len)
	enc.Encode(crypted64, crypted)
	return fmt.Sprintf("$%s$%s", salt, crypted64)
}

// QuickDecrypt 使用aes-gcm与内置默认密钥快速解密
func QuickDecrypt(crypt string) string {
	cryptedl := strings.Split(crypt, "$")
	if len(cryptedl) != 3 {
		return ""
	}
	salt := []byte(cryptedl[1])
	crypted64 := []byte(cryptedl[2])
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	crypted := make([]byte, enc.DecodedLen(len(crypted64)))
	enc.Decode(crypted, crypted64)
	msg, err := AesDecryptGCM(crypted, quickKey, salt)
	if err != nil {
		return ""
	}
	return string(msg)
}

func RandString(len int) string {
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rd.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

/*
加密共有5种模式
- 电码本模式（Electronic Codebook Book (ECB) (已废弃)
- 密码分组链接模式（Cipher Block Chaining (CBC)）
- 计算器模式（Counter (CTR)）
- 密码反馈模式（Cipher FeedBack (CFB)）
- GCM加密模式(Galois/Counter Mode)
*/

func cbcEncrypt(block cipher.Block, plainText []byte, iv ...[]byte) ([]byte, error) {
	ivb := make([]byte, block.BlockSize())
	if len(iv) > 0 {
		ivb = iv[0]
	}
	// 原始数据补码
	bPlainText := pkcs7Padding(plainText, block.BlockSize())
	// 构建加密输出结果
	var cipherText = make([]byte, len(bPlainText))
	// 创建cbc加密模式
	mode := cipher.NewCBCEncrypter(block, ivb)
	mode.CryptBlocks(cipherText, bPlainText)

	return cipherText, nil
}

func cbcDecrypt(block cipher.Block, cipherText []byte, iv ...[]byte) ([]byte, error) {
	ivb := make([]byte, block.BlockSize())
	if len(iv) > 0 {
		ivb = iv[0]
	}
	if (len(cipherText) % block.BlockSize()) != 0 {
		return nil, errors.New("error crypt size")
	}
	plainText := make([]byte, len(cipherText))
	cbc := cipher.NewCBCDecrypter(block, ivb)
	cbc.CryptBlocks(plainText, cipherText)
	plainText = pkcs57Trimming(plainText)
	return plainText, nil
}

func cfbEncrypt(block cipher.Block, plainText []byte, iv ...[]byte) ([]byte, error) {
	ivb := make([]byte, block.BlockSize())
	if len(iv) > 0 {
		ivb = iv[0]
	}
	var cipherText = make([]byte, len(plainText))
	mode := cipher.NewCFBEncrypter(block, ivb)
	mode.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

func cfbDecrypt(block cipher.Block, cipherText []byte, iv ...[]byte) ([]byte, error) {
	ivb := make([]byte, block.BlockSize())
	if len(iv) > 0 {
		ivb = iv[0]
	}
	var plainText = make([]byte, len(cipherText))
	mode := cipher.NewCFBDecrypter(block, ivb)
	mode.XORKeyStream(plainText, cipherText)
	return plainText, nil
}

func gcmEncrypt(block cipher.Block, plainText []byte, nonce []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	dst := make([]byte, gcm.NonceSize())

	cipherText := gcm.Seal(dst, nonce, plainText, nil)
	return cipherText, nil
}

func gcmDecrypt(block cipher.Block, cipherText []byte, nonce []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	cipherText = cipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func pkcs57Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func pkcs7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5Padding(cipherText []byte) []byte {
	return pkcs7Padding(cipherText, 8)
}
