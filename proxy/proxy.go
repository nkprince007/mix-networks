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
	"time"

	"github.com/nkprince007/mix-networks/mixes"
)

const privateKeyPath = "../keys/proxy/private.pem"

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

type Proxy struct {
	mix     mixes.Mix
	privKey *rsa.PrivateKey
	addr    string
}

func (p *Proxy) run() {

	ln, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	go p.handleReqsReadyToForward(p.mix.ReadyToForwardChannel())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go p.handleRequest(conn)
	}
}

func (p *Proxy) handleRequest(conn net.Conn) {
	fmt.Println("Request recieved")
	encryptedMessage := &mixes.EncryptedMessage{}
	json.NewDecoder(conn).Decode(encryptedMessage)
	p.mix.AddMessage(*encryptedMessage)
}

func (p *Proxy) forwardMessage(encryptedMessage mixes.EncryptedMessage) {
	decryptedMsg := mixes.DecryptWithPrivateKey(&encryptedMessage, p.privKey)
	recipientAddr := decryptedMsg.Addr
	recipientMessage := decryptedMsg.Unwrap()
	mixes.SendMessage(&recipientMessage, recipientAddr)
}

func (p *Proxy) handleReqsReadyToForward(readyToForwardChannel chan mixes.MessageBatch) {
	for msgBatch := range readyToForwardChannel {
		for _, msg := range msgBatch.Messages {
			p.forwardMessage(msg)
		}
	}
}

func main() {
	port, err := parseArguments(os.Args[1:])
	//TODO: choose mix strategy based on input argument
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := "127.0.0.1:" + strconv.Itoa(port)
	fmt.Printf("Starting proxy using private key: %s at %s\n", privateKeyPath, addr)
	privKey := mixes.ReadPrivateKey(privateKeyPath)
	mix := getCottrellMix()
	proxy := Proxy{mix, privKey, addr}
	proxy.run()
}
