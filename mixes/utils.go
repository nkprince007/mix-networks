package mixes

import (
	"math/rand"
	"strconv"
	"time"
)

func shuffle(arr []Message) []Message {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(arr),
		func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	return arr
}

func contains(arr []Message, key Message) bool {
	for _, val := range arr {
		if val == key {
			return true
		}
	}
	return false
}

func compareContents(a, b []Message) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !contains(b, v) {
			return false
		}
	}
	return true
}

func makeRange(min, max int) []Message {
	a := make([]Message, max-min+1)
	for i := range a {
		a[i] = Message(min + i)
	}
	return a
}

func getIntervals(timeBufferInMillis time.Duration, percentages []int) []time.Duration {
	var intervals []time.Duration
	for _, p := range percentages {
		intervals = append(intervals, (timeBufferInMillis/100)*time.Duration(p))
	}
	return intervals
}

func sum(units []time.Duration) time.Duration {
	var sum time.Duration
	for _, unit := range units {
		sum += unit
	}
	return sum
}

func messagesToString(pool []Message) string {
	stmt := ""
	for _, msg := range pool {
		stmt += (" " + strconv.Itoa(int(msg)))
	}
	return stmt
}
