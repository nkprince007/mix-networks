package mixes

import (
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	mixPublicKey := ReadPublicKey("./test_public.pem")
	mixPrivateKey := ReadPrivateKey("./test_private.pem")
	message := []byte("This is a test message")

	fmt.Println("Message before encryption : " + string(message))
	encryptedMessage := EncryptWithPublicKey(message, mixPublicKey)
	decryptedMessage := DecryptWithPrivateKey(encryptedMessage, mixPrivateKey)
	fmt.Println("Decrypted message : " + string(decryptedMessage))
}
