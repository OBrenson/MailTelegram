package services

import (
	"YadnexTelegram/internal/configs"
)

type PostMessage struct {
	MailAddr string
	Subject  string
}

type ProducerMessage struct {
}

type MessageHandler interface {
	Handle(message PostMessage)
	ProducerService
}

type ResourceService interface {
	Connect() error
}

type ConsumerService interface {
	ResourceService
	Listen(handler MessageHandler) error
}

type ProducerService interface {
	ResourceService
	Produce(message ProducerMessage) error
}

func NewPostService(config configs.PostConfig) ConsumerService {
	y := &PostConsumer{
		config: config,
	}
	return y
}
