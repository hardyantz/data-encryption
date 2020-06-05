package config

import "github.com/hardyantz/data-encryption/config/rsa"

type Config struct {
	dbPublicKey string
	dbPrivateKey string
}

type ConfImplementation interface {
	PublicKey() string
	PrivateKey() string
}

func NewConfImpl() *Config{
	return &Config{}
}

func (c *Config) PublicKey() string {
	publicKey, _ := rsa.InitPublicKey()
	return publicKey
}

func (c *Config) PrivateKey() string {
	privateKey, _ := rsa.InitPrivateKey()
	return privateKey
}