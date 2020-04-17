package mixes

import (
	"encoding/json"
	"testing"
)

func TestEncryption(t *testing.T) {
	mixPublicKey := ReadPublicKey("./test_public.pem")
	mixPrivateKey := ReadPrivateKey("./test_private.pem")
	message := []byte("This is a test message")

	encryptedMessage := EncryptWithPublicKey(message, mixPublicKey)
	e, err := json.Marshal(&encryptedMessage)
	if err != nil {
		t.Error(err)
	}
	encryptedMessage1 := EncryptWithPublicKey(e, mixPublicKey)

	var decryptedMessage1 EncryptedMessage
	e = DecryptWithPrivateKey(&encryptedMessage1, mixPrivateKey)
	err = json.Unmarshal(e, &decryptedMessage1)
	if err != nil {
		t.Error(err)
	}

	decryptedMessage := DecryptWithPrivateKey(&decryptedMessage1, mixPrivateKey)

	if string(message) != string(decryptedMessage) {
		t.Error("Encrypted and decrypted messages are not the same:", message, decryptedMessage)
	}
}
