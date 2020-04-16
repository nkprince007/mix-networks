package mixes

type Mix interface {
	Forward()
	AddMessage(Message)
	CleanUp()
	GetMessages() []Message
}

type EncryptedMessage struct {
	Content  string `json:"content"`
	Password string `json:"password"`
}
