package services

import (
	"YadnexTelegram/internal/configs"
)

type PostMessage struct {
	MailAddr string
	Subject  string
}

type MessageHandler interface {
	Handle(message PostMessage)
}

type ConsumerService interface {
	Connect() error
	Listen(handler MessageHandler) error
}

func NewPostService(config configs.PostConfig) ConsumerService {
	y := &PostConsumer{
		config: config,
	}
	return y
}
