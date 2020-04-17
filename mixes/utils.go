package mixes

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"time"
)

// Wrap encapsulates encrypted message in another message layer
func (msg *EncryptedMessage) Wrap(addr string) Message {
	e, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return Message{Content: string(e), Addr: addr}
}

// Unwrap returns the wrapped encrypted message
func (msg *Message) Unwrap() EncryptedMessage {
	encMsg := EncryptedMessage{}
	if err := json.Unmarshal([]byte(msg.Content), &encMsg); err != nil {
		log.Fatal(err)
	}
	return encMsg
}

func SendMessage(msg *EncryptedMessage, addr string) {
	// sending request to proxy
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Connecting to %s through tcp failed\n", addr)
	}
	defer conn.Close()
	err = json.NewEncoder(conn).Encode(&msg)
}

func shuffle(arr []EncryptedMessage) []EncryptedMessage {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(arr),
		func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	return arr
}
