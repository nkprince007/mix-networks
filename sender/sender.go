package main

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/nkprince007/mix-networks/mixes"
)

const (
	messageInterval = 500 * time.Millisecond
	recipientAddr   = ":8001"
	proxyAddr       = ":8000"
)

var isPoisson = false
var lambda float64

func init() {
	rand.Seed(time.Now().UnixNano())
	lambda = rand.Float64()
}

func getDelay(lambda float64) uint64 {
	return uint64(rand.ExpFloat64() / lambda)
}

func sendMessage(message string, proxyKey, recipientKey *rsa.PublicKey) {
	// encrypt recipient message
	msg := mixes.Message{Content: message, Addr: ""}
	recipientEncMsg := mixes.EncryptWithPublicKey(&msg, recipientKey)

	// wrap recipient message in one more encryption layer
	msg = recipientEncMsg.Wrap(recipientAddr)
	if isPoisson {
		msg.Delay = getDelay(lambda) * 1000 // converting to milliseconds
	}
	proxyEncMsg := mixes.EncryptWithPublicKey(&msg, proxyKey)

	// sending request to proxy
	mixes.SendMessage(&proxyEncMsg, proxyAddr)
	fmt.Println("Sent request: " + message)
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "poisson" {
		isPoisson = true
	}

	tick := time.Tick(messageInterval)
	counter := 0
	recipientPublicKey := mixes.ReadPublicKey("../.keys/recipient/public.pem")
	proxyPublicKey := mixes.ReadPublicKey("../.keys/proxy/public.pem")

	for {
		select {
		case <-tick:
			counter++
			message := strconv.Itoa(counter)
			if len(os.Args) > 2 {
				message = os.Args[2] + " " + message
			}
			go sendMessage(message, proxyPublicKey, recipientPublicKey)
			if counter == 5 {
				break
			}
		}
	}
}
