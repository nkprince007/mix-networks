package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/nkprince007/mix-networks/mixes"
)

var privateKeyPath string

func usage() {
	programName := os.Args[0]
	fmt.Printf("Usage: %s <key> <port>\n", programName)
	fmt.Println("key:\tPath to RSA private key to be used for this node")
	fmt.Println("port:\tThe port number to start TCP listener on")
	fmt.Println()
}

func parseArguments(args []string) (port int, privKey string, err error) {
	if len(args) != 2 {
		usage()
		err = errors.New("Invalid number of arguments")
		return
	}

	port, err = strconv.Atoi(args[0])
	if err != nil {
		return
	}

	privKey = args[1]
	fileInfo, err := os.Stat(privKey)
	if err == nil && fileInfo.IsDir() {
		err = fmt.Errorf("Expected file, got directory instead at %s", privKey)
	}
	return
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var msg mixes.EncryptedMessage
	err := decoder.Decode(&msg)
	if err != nil {
		fmt.Println("Invalid message: ", err.Error())
		conn.Write([]byte(err.Error()))
		return
	}

	// TODO: Decrypt message, add it to mix and forward when mix conditions are met
}

func main() {
	port, privateKeyPath, err := parseArguments(os.Args[1:])
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := "127.0.0.1:" + string(port)
	fmt.Printf("Starting proxy using private key: %s at %s\n", privateKeyPath, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleRequest(conn)
	}
}
