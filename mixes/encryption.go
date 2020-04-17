package mixes

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// ReadPublicKey gets RSA public key from given path
func ReadPublicKey(filePath string) *rsa.PublicKey {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Reading public key failed: %s", err.Error())
	}
	block, _ := pem.Decode(data)
	b := block.Bytes
	if enc := x509.IsEncryptedPEMBlock(block); enc {
		b, _ = x509.DecryptPEMBlock(block, nil)
	}
	ifc, _ := x509.ParsePKIXPublicKey(b)
	key, _ := ifc.(*rsa.PublicKey)
	return key
}

// ReadPrivateKey gets RSA public key from given path
func ReadPrivateKey(filePath string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Reading private key failed: %s", err.Error())
	}
	block, _ := pem.Decode(data)
	b := block.Bytes
	if enc := x509.IsEncryptedPEMBlock(block); enc {
		b, _ = x509.DecryptPEMBlock(block, nil)
	}
	key, _ := x509.ParsePKCS1PrivateKey(b)
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, _ := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, _ := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	return plaintext
}
