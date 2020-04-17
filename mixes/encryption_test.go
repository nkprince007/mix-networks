package mixes

import (
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	mixPublicKey, mixPrivateKey := getKeyPair()
	message := getMessage()
	fmt.Println("Message before encryption : " + string(message))
	encryptedMessage := EncryptWithPublicKey(message, mixPublicKey)
	decryptedMessage := DecryptWithPrivateKey(encryptedMessage, mixPrivateKey)
	fmt.Println("Decrypted message : " + string(decryptedMessage))
}

func getKeyPair() (*rsa.PublicKey, *rsa.PrivateKey) {
	return fileToPublicKey("../keys/mix/public.pem"), fileToPrivateKey("../keys/mix/private.pem")
}

func getMessage() []byte {
	return []byte("This is a test message")
}
