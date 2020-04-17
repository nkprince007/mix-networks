package mixes

import (
	"testing"
)

func TestEncryption(t *testing.T) {
	mixPublicKey := ReadPublicKey("./test_public.pem")
	mixPrivateKey := ReadPrivateKey("./test_private.pem")
	message := Message{"This is a test message", ":8000"}

	firstLayerEncMsg := EncryptWithPublicKey(&message, mixPublicKey)
	tmpMsg := firstLayerEncMsg.Wrap("hello")
	secondLayerEncMsg := EncryptWithPublicKey(&tmpMsg, mixPublicKey)

	decryptedMessage1 := DecryptWithPrivateKey(&secondLayerEncMsg, mixPrivateKey)
	firstLayerEncMsg = decryptedMessage1.Unwrap()
	decryptedMessage := DecryptWithPrivateKey(&firstLayerEncMsg, mixPrivateKey)

	if message != decryptedMessage {
		t.Error("Encrypted and decrypted messages are not the same:", message, decryptedMessage)
	}
}
