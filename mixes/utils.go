package mixes

import (
	"math/rand"
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
