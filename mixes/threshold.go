package mixes

import (
	"sync"
)

type ThresholdMix struct {
	Size int

	mu                    sync.Mutex
	inputMsgs             []EncryptedMessage
	readyToForwardChannel chan MessageBatch
}

func (m *ThresholdMix) Init() {
	if m.inputMsgs == nil {
		m.inputMsgs = []EncryptedMessage{}
	}

	if m.readyToForwardChannel == nil {
		m.readyToForwardChannel = make(chan MessageBatch, 200) //TODO: Arbitrary size, should we make it configurable?
	}
}

func (m *ThresholdMix) ReadyToForwardChannel() chan MessageBatch {
	return m.readyToForwardChannel
}

func (m *ThresholdMix) Forward() {
	m.readyToForwardChannel <- MessageBatch{Messages: shuffle(m.inputMsgs)}
	m.inputMsgs = []EncryptedMessage{}
}

func (m *ThresholdMix) AddMessage(msg EncryptedMessage) {
	m.mu.Lock()
	m.inputMsgs = append(m.inputMsgs, msg)
	if len(m.inputMsgs) == m.Size {
		m.Forward()
	}
	m.mu.Unlock()
}
