package encrypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// MD5 对字符串进行MD5加密
func MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

func MD5to16(str string) string {
	return MD5(str)[8:24]
}

// Sha1 对字符串进行Sha1加密
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

// Hmac 对值进行Hmac加密
func Hmac(key, value string) string {
	hash := hmac.New(md5.New, []byte(key))
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

// EncryptPassword 对密码进行加密处理
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword 对密码进行比对
func ComparePassword(encryptedPassword string, plaintextPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(plaintextPassword)); err != nil {
		return false, err
	}

	return true, nil
}
