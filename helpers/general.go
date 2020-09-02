package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func Encrypt(plaintext, passphrase string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	block, _ := aes.NewCipher([]byte(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	textByte := gcm.Seal(
		nonce,
		nonce,
		[]byte(plaintext),
		nil)
	return base64.StdEncoding.EncodeToString(textByte), nil
}

func Decrypt(cipherText, key string) (string, error) {
	if cipherText == "" {
		return "", nil
	}

	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()

	textByte, _ := base64.StdEncoding.DecodeString(cipherText)
	nonce, cipherTextByteClean := textByte[:nonceSize], textByte[nonceSize:]
	plaintextByte, err := gcm.Open(
		nil,
		nonce,
		cipherTextByteClean,
		nil)
	if err != nil {
		return "", err
	}

	return string(plaintextByte), nil
}

func EncryptFile(filePath, passPhrase string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return Encrypt(string(content), passPhrase)
}

func DecryptFile(chiperText, passPhrase, output string) error {
	text, err := Decrypt(chiperText, passPhrase)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(output, []byte(text), os.FileMode(777))
}
