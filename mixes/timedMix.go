package mixes

import (
	"fmt"
	"sync"
	"time"
)

type TimedMix struct {
	TimeBufferMillis      time.Duration
	mu                    sync.Mutex
	inputMsgs             []EncryptedMessage
	ticker                *time.Ticker
	readyToForwardChannel chan MessageBatch
}

func (m *TimedMix) Init() {
	if m.inputMsgs == nil {
		m.inputMsgs = []EncryptedMessage{}
	}

	if m.readyToForwardChannel == nil {
		m.readyToForwardChannel = make(chan MessageBatch, 200) //TODO: Arbitrary size, should we make it configurable?
	}
	if m.ticker == nil {
		go initTicker(m) //TODO: start ticker on first message or on initalization?
	}
}

func initTicker(m *TimedMix) {
	m.ticker = time.NewTicker(m.TimeBufferMillis)
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.Forward()
			}
		}
	}()
}

func (m *TimedMix) ReadyToForwardChannel() chan MessageBatch {
	return m.readyToForwardChannel
}

func (m *TimedMix) Forward() {
	fmt.Println("Tick")
	m.mu.Lock()
	m.readyToForwardChannel <- MessageBatch{Messages: shuffle(m.inputMsgs)}
	m.inputMsgs = []EncryptedMessage{}
	m.mu.Unlock()
}

func (m *TimedMix) AddMessage(msg EncryptedMessage) {
	m.mu.Lock()
	m.inputMsgs = append(m.inputMsgs, msg)
	m.mu.Unlock()
}
