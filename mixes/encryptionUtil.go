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

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		log.Println("is encrypted pem block")
		b, _ = x509.DecryptPEMBlock(block, nil)
	}
	key, _ := x509.ParsePKCS1PrivateKey(b)
	return key
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileToPublicKey(filePath string) *rsa.PublicKey {
	dat, err := ioutil.ReadFile(filePath)
	check(err)
	return BytesToPublicKey(dat)
}

func fileToPrivateKey(filePath string) *rsa.PrivateKey {
	dat, err := ioutil.ReadFile(filePath)
	check(err)
	return BytesToPrivateKey(dat)
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		log.Println("is encrypted pem block")
		b, _ = x509.DecryptPEMBlock(block, nil)
	}
	ifc, _ := x509.ParsePKIXPublicKey(b)
	key, _ := ifc.(*rsa.PublicKey)
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
