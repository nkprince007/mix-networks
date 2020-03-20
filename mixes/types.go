package mixes

type Mix interface {
	Forward()
	AddMessage(Message)
	CleanUp()
	GetMessages() []Message
}
