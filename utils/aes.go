package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"github.com/astaxie/beego"
	"net/url"
	"strings"
)

var aes_Key = "123456789"

func init()  {
	if str := beego.AppConfig.String("aes"+ "::key");str != ""{
		aes_Key = str
	}
}


//加密
func EncryptWithAESUrl(message string) (string,bool) {

	hash := md5.New()
	hash.Write([]byte(aes_Key))
	keyData := hash.Sum(nil)

	block, err := aes.NewCipher(keyData)
	if err != nil {
		return err.Error(),false
	}

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	enc := cipher.NewCBCEncrypter(block, iv)
	content := pKCS5Padding([]byte(message), block.BlockSize())
	crypted := make([]byte, len(content))
	enc.CryptBlocks(crypted, content)
	str := base64.StdEncoding.EncodeToString(crypted)

	return encodeURIComponent(str),true
}

func EncryptWithAES(message string) (string,bool) {

	hash := md5.New()
	hash.Write([]byte(aes_Key))
	keyData := hash.Sum(nil)

	block, err := aes.NewCipher(keyData)
	if err != nil {
		return err.Error(),false
	}

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	enc := cipher.NewCBCEncrypter(block, iv)
	content := pKCS5Padding([]byte(message), block.BlockSize())
	crypted := make([]byte, len(content))
	enc.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted),true
}

//解密
func DecryptWithAES(message string) (string,bool) {


	hash := md5.New()
	hash.Write([]byte(aes_Key))
	keyData := hash.Sum(nil)

	block, err := aes.NewCipher(keyData)
	if err != nil {
		return err.Error(),false
	}

	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

	messageData, _ := base64.StdEncoding.DecodeString(message)
	dec := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(messageData))
	dec.CryptBlocks(decrypted, messageData)
	return string(pKCS5Unpadding(decrypted)),true
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5Unpadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}