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

const privateKeyPath = "../.keys/proxy/private.pem"
const dropRequestsArg = "drop"

func usage() {
	programName := os.Args[0]
	fmt.Printf("Usage: %s <mix>\n", programName)
	fmt.Println("mix:\tThe mixing strategy to be used")
	fmt.Println("\t\tshould be one of threshold, timed, cottrell or RGB")
	fmt.Println()
}

func parseArguments(args []string) (mix mixes.Mix, err error, dropRequests bool) {
	if len(args) == 0 || len(args) > 2 {
		usage()
		err = errors.New("Invalid number of arguments")
		return
	}

	switch args[0] {
	case "threshold":
		mix = getThresholdMix()
	case "timed":
		mix = getTimedMix()
	case "cottrell":
		mix = getCottrellMix()
	case "RGB":
		mix = getRGBMix()
	case "poisson":
		mix = getPoissonMix()
	default:
		usage()
		err = errors.New("Unrecognized mix, ")
	}

	dropRequests = false
	if len(args) == 2 {
		dropRequests = (args[1] == dropRequestsArg)
	}
	return
}

type proxy struct {
	mix      mixes.Mix
	privKey  *rsa.PrivateKey
	addr     string
	dropReqs bool
	reqCount int
}

func (p *proxy) run() {
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
		if p.dropReqs {
			go p.handleRequestWithDrop(conn)
		} else {
			go p.handleRequest(conn)
		}
	}
}

func (p *proxy) handleRequest(conn net.Conn) {
	fmt.Println("Request recieved")
	encryptedMessage := &mixes.EncryptedMessage{}
	json.NewDecoder(conn).Decode(encryptedMessage)
	p.mix.AddMessage(*encryptedMessage)
}

func (p *proxy) handleRequestWithDrop(conn net.Conn) {
	fmt.Println("Req no: " + strconv.Itoa(p.reqCount))
	p.reqCount++
	var dropProbability float32
	if p.reqCount > 100 {
		dropProbability = float32(float32(p.reqCount%100) / 100)
	} else {
		dropProbability = float32(float32(p.reqCount) / 100)
	}
	dropMessage := !mixes.PickTrueWithProbability(dropProbability)
	if !dropMessage {
		p.handleRequest(conn)
	} else {
		fmt.Println("Dropped request")
	}
}

func (p *proxy) forwardMessage(encryptedMessage mixes.EncryptedMessage) {
	decryptedMsg := mixes.DecryptWithPrivateKey(&encryptedMessage, p.privKey)
	recipientAddr := decryptedMsg.Addr
	recipientMessage := decryptedMsg.Unwrap()
	if recipientAddr != "" {
		mixes.SendMessage(&recipientMessage, recipientAddr)
	}
}

func (p *proxy) handleReqsReadyToForward(readyToForwardChannel chan mixes.MessageBatch) {
	for msgBatch := range readyToForwardChannel {
		for _, msg := range msgBatch.Messages {
			p.forwardMessage(msg)
		}
	}
}

func main() {
	mix, err, dropReqs := parseArguments(os.Args[1:])
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := ":8000"
	fmt.Printf("Starting proxy using private key: %s at %s\n", privateKeyPath, addr)
	privKey := mixes.ReadPrivateKey(privateKeyPath)
	p := proxy{mix, privKey, addr, dropReqs, 0}
	p.run()
}
