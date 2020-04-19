package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/nkprince007/mix-networks/mixes"
)

const privateKeyPath = "../.keys/recipient/private.pem"

func handleRequest(conn net.Conn, privKey *rsa.PrivateKey) {
	encryptedMessage := &mixes.EncryptedMessage{}
	json.NewDecoder(conn).Decode(encryptedMessage)

	msg := mixes.DecryptWithPrivateKey(encryptedMessage, privKey)
	fmt.Println("Received response: ", msg.Content)
}

func main() {
	addr := ":8001"
	fmt.Printf("Starting recipient using private key: %s at %s\n", privateKeyPath, addr)
	privKey := mixes.ReadPrivateKey(privateKeyPath)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleRequest(conn, privKey)
	}
}
