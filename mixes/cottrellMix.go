package mixes

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type CottrellMix struct {
	MinimumPoolSize       int
	Threshold             int
	Fraction              float32
	mu                    sync.Mutex
	inputMsgs             []EncryptedMessage
	readyToForwardChannel chan MessageBatch
	TimeBufferMillis      time.Duration
	ticker                *time.Ticker
}

func (m *CottrellMix) Init() {
	if m.inputMsgs == nil {
		m.inputMsgs = []EncryptedMessage{}
	}

	if m.readyToForwardChannel == nil {
		m.readyToForwardChannel = make(chan MessageBatch, 200) //TODO: Arbitrary size, should we make it configurable?
	}
	if m.ticker == nil {
		go initCTicker(m) //TODO: start ticker on first message or on initalization?
	}
}

func initCTicker(m *CottrellMix) {
	m.ticker = time.NewTicker(m.TimeBufferMillis)
	go func() {
		for {
			select {
			case <-m.ticker.C:
				fmt.Println("Tick. The buffer size is " + strconv.Itoa(len(m.inputMsgs)))
				if len(m.inputMsgs) >= m.Threshold+m.MinimumPoolSize {
					m.Forward()
				}
			}
		}
	}()
}

func (m *CottrellMix) ReadyToForwardChannel() chan MessageBatch {
	return m.readyToForwardChannel
}

func (m *CottrellMix) Forward() {
	m.mu.Lock()
	fmt.Println("Forwarding messages")
	n := int(m.Fraction * float32(len(m.inputMsgs)))
	fmt.Println("n is " + strconv.Itoa(n))
	m.inputMsgs = shuffle(m.inputMsgs)
	m.readyToForwardChannel <- MessageBatch{Messages: m.inputMsgs[len(m.inputMsgs)-n:]}
	m.inputMsgs = m.inputMsgs[:len(m.inputMsgs)-n]
	fmt.Println("Forwarded messages, buffer is now " + strconv.Itoa(len(m.inputMsgs)))
	m.mu.Unlock()
}

func (m *CottrellMix) AddMessage(msg EncryptedMessage) {
	m.mu.Lock()
	m.inputMsgs = append(m.inputMsgs, msg)
	m.mu.Unlock()
}
