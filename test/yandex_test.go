package a_test

import (
	"YadnexTelegram/internal/services"
	"YadnexTelegram/internal/util"
	"testing"
)

var test *testing.T

func TestPost(t *testing.T) {
	test = t
	p := util.GetPostConfig("./resources/configs.yaml")
	y := services.NewPostService(p).(services.ConsumerService)
	err := y.Listen(&HandlerMock{})
	if err != nil {
		t.Fatal(err)
	}
}

type HandlerMock struct {
}

func (h HandlerMock) Handle(message services.PostMessage) {
	test.Log(message.MailAddr)
	test.Log(message.Subject)
}

func (h HandlerMock) Produce(message services.ProducerMessage) error {
	return nil
}

func (h HandlerMock) Connect() error {
	return nil
}
