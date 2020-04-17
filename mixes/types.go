package mixes

type Mix interface {
	Forward()
	AddMessage(EncryptedMessage)
	CleanUp()
	GetMessages() []EncryptedMessage
}

type Message struct {
	Content string `json:"content"`
	Addr    string `json:"address"`
}

type EncryptedMessage struct {
	Data     []byte `json:"data"`
	Nonce    []byte `json:"nonce"`
	Password []byte `json:"password"`
}

type MessageBatch struct {
	Messages []EncryptedMessage
}

type MixNew interface {
	AddMessage(EncryptedMessage)
	ReadyToForwardChannel() chan MessageBatch
}

type ReqSender interface {
	AddRequest(EncryptedMessage)
}

type ReqReciever interface {
	ProcessRequest(EncryptedMessage)
}
