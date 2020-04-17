package main

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/nkprince007/mix-networks/mixes"
)

const (
	messageInterval = 500 * time.Millisecond
	recipientAddr   = ":8001"
	proxyAddr       = ":8000"
)

func sendMessage(message string, proxyKey, recipientKey *rsa.PublicKey) {
	// encrypt recipient message
	msg := mixes.Message{Content: message, Addr: ""}
	recipientEncMsg := mixes.EncryptWithPublicKey(&msg, recipientKey)

	// wrap recipient message in one more encryption layer
	msg = recipientEncMsg.Wrap(recipientAddr)
	proxyEncMsg := mixes.EncryptWithPublicKey(&msg, proxyKey)

	// sending request to proxy
	mixes.SendMessage(&proxyEncMsg, proxyAddr)
	fmt.Println("Sent request: ", string(msg.Content))
}

func main() {
	tick := time.Tick(messageInterval)
	counter := 0
	recipientPublicKey := mixes.ReadPublicKey("./sender/recipient-pubkey.pem")
	proxyPublicKey := mixes.ReadPublicKey("./sender/proxy-pubkey.pem")

	for {
		select {
		case <-tick:
			counter++
			go sendMessage(strconv.Itoa(counter), proxyPublicKey, recipientPublicKey)
		}
	}
}
