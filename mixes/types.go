package mixes

type Mix interface {
	Forward()
	AddMessage(Message)
	CleanUp()
	GetMessages() []Message
}

type EncryptedMessage struct {
	Content  []byte `json:"content"`
	Nonce    []byte `json:"nonce"`
	Password []byte `json:"password"`
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
