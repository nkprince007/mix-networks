package main

import (
	"fmt"
	"testing"
)

type EncryptedMessage struct {
	Content  string `json:"content"`
	Password string `json:"password"`
}

type MessageBatch struct {
	messages []EncryptedMessage
}

type MixNew interface {
	AddMessage(EncryptedMessage)
	readyToForwardChannel() chan []MessageBatch
}

type ReqSender interface {
	AddRequest(EncryptedMessage)
}

type ReqReciever interface {
	ProcessRequest(EncryptedMessage)
}

type ConsoleReader struct {
}

func (c *ConsoleReader) init() {
	for {
		var input string
		fmt.Scanln(&input)
		if input == "quit" {
			break
		}
		var message EncryptedMessage
		message.Content = input
		c.AddRequest(message)
	}
}

func (c *ConsoleReader) AddRequest(message EncryptedMessage) {
	fmt.Println("Added request " + message.Content)
}

func TestMix(t *testing.T) {
	// reader := ConsoleReader{}
	// reader.init()
	fmt.Println("This is a test")
	var input string
	fmt.Scanln(&input)
	fmt.Println(input)
}
