package rsa

import (
	"io/ioutil"
)

const (
	privateKeyPath = "config/rsa/db.rsa"
	publicKeyPath  = "config/rsa/db.rsa.pub"
)

// InitPublicKey return *rsa.PublicKey
func InitPublicKey() (string, error) {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	return string(verifyBytes), nil
}

// InitPrivateKey return *rsa.PrivateKey
func InitPrivateKey() (string, error) {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	return string(signBytes), nil
}
