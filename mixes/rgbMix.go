package mixes

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type RgbMix struct {
	mu                            sync.Mutex
	inputMsgs                     []EncryptedMessage
	readyToForwardChannel         chan MessageBatch
	PeriodMillis                  time.Duration
	periodTicker                  *time.Ticker
	c                             ReqCounter
	ProxyPublicKey                *rsa.PublicKey
	ProxyPrivateKey               *rsa.PrivateKey
	GreenMessageRecieverPublicKey *rsa.PublicKey
}

const (
	seedRedMessageCount      = int(20)
	proxyAddr                = ":8000"
	proxySecretRedMessage    = "fed"
	proxySecretGreenMessage  = "green"
	greenMessageRecieverAddr = ":8001"
)

type ReqCounter struct {
	NoOfRedMessagesSentInLastSession int
	NoOfMessageSentInLastSession     int
	totalNoOfMessagesRecieved        int
	redMessagesRecieved              int
	blackMessagesRecieved            int
}

func (m *RgbMix) Init() {
	if m.inputMsgs == nil {
		m.inputMsgs = []EncryptedMessage{}
	}

	if m.readyToForwardChannel == nil {
		m.readyToForwardChannel = make(chan MessageBatch, 200) //TODO: Arbitrary size, should we make it configurable?
	}
	m.ProxyPublicKey = ReadPublicKey("../.keys/proxy/public.pem")
	m.ProxyPrivateKey = ReadPrivateKey("../.keys/proxy/private.pem")
	m.GreenMessageRecieverPublicKey = ReadPublicKey("../.keys/recipient/public.pem")
	m.initPeriodTicker()
}

func (m *RgbMix) initPeriodTicker() {
	m.periodTicker = time.NewTicker(m.PeriodMillis)
	go func() {
		for {
			select {
			case <-m.periodTicker.C:
				noOfRedMessagesSent := m.injectRedMessages()
				var noOfGreenMessagesSent int
				if m.c.NoOfMessageSentInLastSession != 0 {
					noOfGreenMessagesSent = m.injectGreenMessages()
				}
				totalNoOfMessagesSent := len(m.inputMsgs)
				fmt.Println("---------Sending----------")
				fmt.Println("Green " + strconv.Itoa(noOfGreenMessagesSent))
				fmt.Println("Red " + strconv.Itoa(noOfRedMessagesSent))
				fmt.Println("Black " + strconv.Itoa(totalNoOfMessagesSent-noOfRedMessagesSent-noOfGreenMessagesSent))
				fmt.Println("-------------------------------")
				m.Forward()
				m.c = ReqCounter{
					NoOfRedMessagesSentInLastSession: noOfRedMessagesSent,
					NoOfMessageSentInLastSession:     totalNoOfMessagesSent,
				}
				fmt.Println("Reset counter")
			}
		}
	}()
}

func (m *RgbMix) ReadyToForwardChannel() chan MessageBatch {
	return m.readyToForwardChannel
}

func (m *RgbMix) Forward() {
	m.readyToForwardChannel <- MessageBatch{Messages: shuffle(m.inputMsgs)}
	m.inputMsgs = []EncryptedMessage{}
}

func (m *RgbMix) injectRedMessages() int {
	var noOfRedMessagesToSend int
	if m.c.NoOfMessageSentInLastSession == 0 {
		noOfRedMessagesToSend = seedRedMessageCount
	} else {
		R, B, r := m.calculateVariables()
		noOfRedMessagesToSend = int((R + B) * r)
	}
	for i := 0; i < noOfRedMessagesToSend; i++ {
		m.inputMsgs = append(m.inputMsgs, m.ComposeRedMessage())
	}
	return noOfRedMessagesToSend
}

func (m *RgbMix) injectGreenMessages() int {
	R, B, r := m.calculateVariables()
	fmt.Println(fmt.Sprintf("R:- %f | B :- %f | r :- %f", R, B, r))
	G := int((((R + B) * r) * (1 - r)) / r)
	for i := 0; i < G; i++ {
		m.inputMsgs = append(m.inputMsgs, m.ComposeGreenMessage())
	}
	return G
}

func (m *RgbMix) calculateVariables() (float32, float32, float32) {
	return float32(m.c.redMessagesRecieved), //R
		float32(m.c.totalNoOfMessagesRecieved - m.c.redMessagesRecieved), //B
		float32(float32(m.c.NoOfRedMessagesSentInLastSession) / float32(m.c.NoOfMessageSentInLastSession)) //r
}

func (m *RgbMix) ComposeGreenMessage() EncryptedMessage {
	encryptedMsg1 := EncryptWithPublicKey(&Message{
		Content: proxySecretGreenMessage,
		Addr:    "",
	}, m.GreenMessageRecieverPublicKey)
	wrappedEncryptedMsg1 := encryptedMsg1.Wrap(greenMessageRecieverAddr)
	return EncryptWithPublicKey(&wrappedEncryptedMsg1, m.ProxyPublicKey)
}

func (m *RgbMix) ComposeRedMessage() EncryptedMessage {
	encryptedMsg1 := EncryptWithPublicKey(&Message{
		Content: proxySecretRedMessage,
		Addr:    "",
	}, m.ProxyPublicKey)
	wrappedEncryptedMsg1 := encryptedMsg1.Wrap(proxyAddr)
	return EncryptWithPublicKey(&wrappedEncryptedMsg1, m.ProxyPublicKey)
}

func (m *RgbMix) ProcessRequest(msg EncryptedMessage) bool {
	decryptedMsg := DecryptWithPrivateKey(&msg, m.ProxyPrivateKey)
	isRedMessage := false
	if decryptedMsg.Content == proxySecretRedMessage {
		m.c.redMessagesRecieved++
		isRedMessage = true
	} else {
		m.c.blackMessagesRecieved++
	}
	m.c.totalNoOfMessagesRecieved++
	fmt.Println("Black :- " + strconv.Itoa(m.c.blackMessagesRecieved) +
		" Red:- " + strconv.Itoa(m.c.redMessagesRecieved))
	return isRedMessage
}

func (m *RgbMix) AddMessage(msg EncryptedMessage) {
	m.mu.Lock()
	isRedMessage := m.ProcessRequest(msg)
	if !isRedMessage {
		m.inputMsgs = append(m.inputMsgs, msg)
	}
	m.mu.Unlock()
}
