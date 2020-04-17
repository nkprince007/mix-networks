package main

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/nkprince007/mix-networks/mixes"
// )

// type MessageBatch struct {
// 	messages []mixes.Message
// }

// type MixNew interface {
// 	AddMessage(mixes.Message)
// 	readyToForwardChannel() chan []MessageBatch
// }

// type ReqSender interface {
// 	AddRequest(mixes.Message)
// }

// type ReqReciever interface {
// 	ProcessRequest(mixes.Message)
// }

// type ConsoleReader struct {
// }

// func (c *ConsoleReader) init() {
// 	for {
// 		var input string
// 		fmt.Scanln(&input)
// 		if input == "quit" {
// 			break
// 		}
// 		var message mixes.Message
// 		message.Content = input
// 		c.AddRequest(message)
// 	}
// }

// func (c *ConsoleReader) AddRequest(message mixes.Message) {
// 	fmt.Println("Added request " + message.Content)
// }

// func TestMix(t *testing.T) {
// 	// reader := ConsoleReader{}
// 	// reader.init()
// 	fmt.Println("This is a test")
// 	var input string
// 	fmt.Scanln(&input)
// 	fmt.Println(input)
// }
