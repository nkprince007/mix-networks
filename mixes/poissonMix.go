package mixes

import (
	"crypto/rsa"
	"time"
)

type PoissonMix struct {
	readyToForwardChannel chan MessageBatch
	privKey               *rsa.PrivateKey
}

func (m *PoissonMix) Init() {
	if m.readyToForwardChannel == nil {
		m.readyToForwardChannel = make(chan MessageBatch, 200) //TODO: Arbitrary size, should we make it configurable?
	}

	if m.privKey == nil {
		m.privKey = ReadPrivateKey("../.keys/proxy/private.pem")
	}
}

func (m *PoissonMix) ReadyToForwardChannel() chan MessageBatch {
	return m.readyToForwardChannel
}

func (m *PoissonMix) processMessage(msg EncryptedMessage) {
	decryptedMsg := DecryptWithPrivateKey(&msg, m.privKey)
	time.Sleep(time.Duration(decryptedMsg.Delay) * time.Millisecond)
	m.readyToForwardChannel <- MessageBatch{[]EncryptedMessage{msg}}
}

func (m *PoissonMix) AddMessage(msg EncryptedMessage) {
	go m.processMessage(msg)
}
