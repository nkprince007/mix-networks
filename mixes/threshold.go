package mixes

import (
	"fmt"
	"sync"
)

type ThresholdMix struct {
	Size int

	mu   sync.Mutex
	msgs chan Message
	done chan interface{}
}

func (m *ThresholdMix) init() {
	if m.done == nil {
		m.done = make(chan interface{})
	}

	if m.msgs == nil {
		m.msgs = make(chan Message, m.Size)
	}
}

func (m *ThresholdMix) Forward() {
	close(m.msgs)
	fmt.Println("Channel is full")
	m.done <- nil
}

func (m *ThresholdMix) AddMessage(msg Message) {
	m.init()

	m.mu.Lock()
	fmt.Println("received msg: ", msg)

	select {
	case m.msgs <- msg:
		if len(m.msgs) == cap(m.msgs) {
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
	close(m.msgs)
	close(m.done)
}

func (m *ThresholdMix) GetMessages() (msgs []Message) {
	m.init()

	select {
	case <-m.done:
		defer m.mu.Unlock()
		for msg := range m.msgs {
			msgs = append(msgs, msg)
		}

		m.msgs = make(chan Message, m.Size)
		msgs = shuffle(msgs)
	}
	return
}
