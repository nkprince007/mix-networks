package mixes

import (
	"sync"
	"testing"
)

func TestThresholdMix(t *testing.T) {
	mix := ThresholdMix{Size: 2}
	input := []Message{1, 2, 3, 4, 5, 6}
	mux := &sync.Mutex{}
	sentMessages := []Message{}

	for _, msg := range input {
		go func(mesg Message) {
			mux.Lock()
			defer mux.Unlock()
			mix.AddMessage(mesg)
			sentMessages = append(sentMessages, mesg)
		}(msg)
	}

	mixRatio := len(input) / mix.Size
	for i := 0; i < mixRatio; i++ {
		output := mix.GetMessages()
		if !compareContents(sentMessages, output) {
			t.Errorf("expected: %v, got: %v", sentMessages, output)
		}
		sentMessages = []Message{}
	}
}
