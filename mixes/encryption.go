package mixes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
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
func EncryptWithPublicKey(msg *Message, pub *rsa.PublicKey) EncryptedMessage {
	plaintext, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("JSON Encoding error:", plaintext, err)
	}

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Fatal(err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err.Error())
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, key, nil)
	if err != nil {
		log.Fatalln("encryption error:", err.Error())
	}

	return EncryptedMessage{
		Data:     ciphertext,
		Nonce:    nonce,
		Password: encrypted,
	}
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(msg *EncryptedMessage, priv *rsa.PrivateKey) Message {
	key := msg.Password
	ciphertext := msg.Data
	nonce := msg.Nonce

	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, key, nil)
	if err != nil {
		log.Fatalln("decryption error:", err.Error())
	}

	block, err := aes.NewCipher(decrypted)
	if err != nil {
		log.Fatal(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err.Error())
	}

	deciphered, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	var message Message
	err = json.Unmarshal(deciphered, &message)
	if err != nil {
		log.Fatal("JSON decoding error:", err)
	}
	return message
}
