package mixes

import (
	"fmt"
	"sync"
)

type ThresholdMix struct {
	Size int

	mu          sync.Mutex
	msgsChannel chan Message
	done        chan interface{}
}

func (m *ThresholdMix) init() {
	if m.done == nil {
		m.done = make(chan interface{})
	}

	if m.msgsChannel == nil {
		m.msgsChannel = make(chan Message, m.Size)
	}
}

func (m *ThresholdMix) Forward() {
	close(m.msgsChannel)
	fmt.Println("Channel is full")
	m.done <- nil
}

func (m *ThresholdMix) AddMessage(msg Message) {
	m.init()

	m.mu.Lock()
	fmt.Println("received msg: ", msg)

	select {
	case m.msgsChannel <- msg:
		if len(m.msgsChannel) == cap(m.msgsChannel) {
			m.Forward()
		} else {
			m.mu.Unlock()
		}
	default:
		m.Forward()
	}
}

func (m *ThresholdMix) CleanUp() {
	fmt.Println("Cleanup complete...")
	close(m.msgsChannel)
	close(m.done)
}

func (m *ThresholdMix) GetMessages() (msgs []Message) {
	m.init()

	select {
	case <-m.done:
		defer m.mu.Unlock()
		for msg := range m.msgsChannel {
			msgs = append(msgs, msg)
		}

		m.msgsChannel = make(chan Message, m.Size)
		msgs = shuffle(msgs)
	}
	return
}
