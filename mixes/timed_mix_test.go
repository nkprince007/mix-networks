package mixes

import (
	"fmt"
	"testing"
	"time"
)

func TestTimedMix(t *testing.T) {

	//Setup a time pool size
	mixTimeBufferSize := 500 * time.Millisecond
	mix := TimedMix{TimeBufferMillis: mixTimeBufferSize}

	//Setting up intervals between message, using intervals less than and greater than the pool size
	intervals := getIntervals(mixTimeBufferSize, []int{100, 70, 20, 10, 10, 10, 10})
	input := makeRange(1, len(intervals))

	//The maximum time that the test will wait for the mix to return all messages
	maxWaitTime := mixTimeBufferSize * time.Duration(len(intervals))
	fmt.Println(mixTimeBufferSize)
	fmt.Println(maxWaitTime)

	//Feeding messages at the given intervals into the mix
	for i := 0; i < len(input); i++ {
		func(msg Message, minTimeInterval time.Duration) {
			time.Sleep(minTimeInterval)
			mix.AddMessage(msg)
		}(input[i], intervals[i])
	}

	// Retrieve messages when available and compare against input
	var allRecievedMessages []Message
	waitTimer := time.NewTimer(4000 * time.Millisecond)

	for {
		select {
		case <-waitTimer.C:
			fmt.Println("Timeout, breaking")
			if !compareContents(input, allRecievedMessages) {
				t.Errorf("expected: %v, got: %v", input, allRecievedMessages)
			}
			return
		default:
			output := mix.GetMessages()
			fmt.Println("Recieved :- " + messagesToString(output))
			allRecievedMessages = append(allRecievedMessages, output...)
		}
	}
}
