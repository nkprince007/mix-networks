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
	seedRedMessageCount      = int(5)
	proxyAddr                = ":8000"
	proxySecretRedMessage    = "red"
	proxySecretGreenMessage  = "green"
	greenMessageRecieverAddr = "8001"
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
}

func initPeriodTicker(m *RgbMix) {
	m.periodTicker = time.NewTicker(m.PeriodMillis)
	go func() {
		for {
			select {
			case <-m.periodTicker.C:
				fmt.Println("Tick. The buffer size is " + strconv.Itoa(len(m.inputMsgs)))
				noOfRedMessagesSent := m.injectRedMessages()
				if &m.c != nil {
					m.injectGreenMessages()
				}
				totalNoOfMessagesSent := len(m.inputMsgs)
				m.Forward()
				m.c = ReqCounter{
					NoOfRedMessagesSentInLastSession: noOfRedMessagesSent,
					NoOfMessageSentInLastSession:     totalNoOfMessagesSent,
				}
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
	if &m.c == nil {
		noOfRedMessagesToSend = seedRedMessageCount
	} else {
		R := float32(m.c.redMessagesRecieved)
		B := float32(m.c.NoOfMessageSentInLastSession - m.c.NoOfRedMessagesSentInLastSession)
		r := float32(float32(m.c.NoOfRedMessagesSentInLastSession) / float32(m.c.NoOfMessageSentInLastSession))
		noOfRedMessagesToSend = int((R + B) * r)
	}
	for i := 0; i < noOfRedMessagesToSend; i++ {
		m.inputMsgs = append(m.inputMsgs, m.ComposeRedMessage())
	}
	return noOfRedMessagesToSend
}

func (m *RgbMix) injectGreenMessages() {
	R := float32(m.c.redMessagesRecieved)
	B := float32(m.c.NoOfMessageSentInLastSession - m.c.NoOfRedMessagesSentInLastSession)
	r := float32(float32(m.c.NoOfRedMessagesSentInLastSession) / float32(m.c.NoOfMessageSentInLastSession))
	G := int((((R + B) * r) * (1 - r)) / r)
	for i := 0; i < G; i++ {
		m.inputMsgs = append(m.inputMsgs, m.ComposeGreenMessage())
	}
}

func (m *RgbMix) ComposeGreenMessage() EncryptedMessage {
	return EncryptWithPublicKey(&Message{Content: proxySecretGreenMessage, Addr: greenMessageRecieverAddr}, m.GreenMessageRecieverPublicKey)
}

func (m *RgbMix) ComposeRedMessage() EncryptedMessage {
	return EncryptWithPublicKey(&Message{Content: proxySecretRedMessage, Addr: proxyAddr}, m.ProxyPublicKey)
}

func (m *RgbMix) ProcessRequest(msg EncryptedMessage) {
	decryptedMsg, err := DecryptWithPrivateKey(&msg, m.ProxyPrivateKey)
	if err != nil && decryptedMsg.Content == proxySecretRedMessage {
		m.c.redMessagesRecieved++
	} else {
		m.c.blackMessagesRecieved++
	}
	m.c.totalNoOfMessagesRecieved++
}

func (m *RgbMix) AddMessage(msg EncryptedMessage) {
	m.mu.Lock()
	m.ProcessRequest(msg)
	m.inputMsgs = append(m.inputMsgs, msg)
	m.mu.Unlock()
}
