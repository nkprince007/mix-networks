package mixes

type Message struct {
	Content string `json:"content"`
	Addr    string `json:"address"`
	Delay   uint64 `json:"delay"`
}

type EncryptedMessage struct {
	Data     []byte `json:"data"`
	Nonce    []byte `json:"nonce"`
	Password []byte `json:"password"`
}

type MessageBatch struct {
	Messages []EncryptedMessage
}

type Mix interface {
	AddMessage(EncryptedMessage)
	ReadyToForwardChannel() chan MessageBatch
	Init()
}
