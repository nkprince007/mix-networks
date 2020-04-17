package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/nkprince007/mix-networks/mixes"
)

const privateKeyPath = "./proxy/proxy-privkey.pem"

func usage() {
	programName := os.Args[0]
	fmt.Printf("Usage: %s <port>\n", programName)
	fmt.Println("port:\tThe port number to start TCP listener on")
	fmt.Println()
}

func parseArguments(args []string) (port int, err error) {
	if len(args) != 1 {
		usage()
		err = errors.New("Invalid number of arguments")
		return
	}

	port, err = strconv.Atoi(args[0])
	return
}

func handleRequest(conn net.Conn, privKey *rsa.PrivateKey) {
	encryptedMessage := &mixes.EncryptedMessage{}
	json.NewDecoder(conn).Decode(encryptedMessage)

	msg := mixes.DecryptWithPrivateKey(encryptedMessage, privKey)
	unwrappedMessage := msg.Unwrap()
	fmt.Println(unwrappedMessage)

	// TODO: Add it to mix and forward when mix conditions are met
}

func main() {
	port, err := parseArguments(os.Args[1:])
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := "127.0.0.1:" + string(port)
	fmt.Printf("Starting proxy using private key: %s at %s\n", privateKeyPath, addr)
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
