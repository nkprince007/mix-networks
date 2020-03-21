package mixes

import (
	"fmt"
	"sync"
	"time"
)

type TimedMix struct {
	TimeBufferMillis        time.Duration
	mu                      sync.Mutex
	msgsPool                []Message //All pooled messages which are yet to be picked
	ticker                  *time.Ticker
	msgsRecieveReadyChannel chan []Message //Channel with message chunks which are ready to be picked
}

func (m *TimedMix) init() {
	if m.msgsPool == nil {
		m.msgsPool = make([]Message, 0)
	}
	if m.msgsRecieveReadyChannel == nil {
		//TODO: How many messages will the buffer hold before panicking? Can this be made configurable?
		m.msgsRecieveReadyChannel = make(chan []Message, 30)
	}
	if m.ticker == nil {
		go initTicker(m)
	}
}

func initTicker(m *TimedMix) {
	fmt.Println("Starting ticker at:-")
	fmt.Println(time.Now())
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

func (m *TimedMix) Forward() {
	//[TimeBufferMillis] time has passed between the last batch
	// Moving messages from pool to recieve channel
	m.mu.Lock()
	if len(m.msgsPool) > 0 {
		fmt.Println("Moving " + messagesToString(m.msgsPool) + " from pool to channel at:-")
		fmt.Println(time.Now())
		m.msgsRecieveReadyChannel <- m.msgsPool
		m.msgsPool = make([]Message, 0)
	} else {
		fmt.Println("Empty pool, ignoring forward")
	}
	m.mu.Unlock()
}

func (m *TimedMix) AddMessage(msg Message) {
	m.init() //TODO: should init be common
	m.mu.Lock()
	m.msgsPool = append(m.msgsPool, msg)
	m.mu.Unlock()
}

func (m *TimedMix) CleanUp() {
	fmt.Println("Cleanup complete...")
	m.Forward()
}

func (m *TimedMix) GetMessages() (msgs []Message) {
	m.init() //TODO: should init be common
	msgs = shuffle(<-m.msgsRecieveReadyChannel)
	return
}
