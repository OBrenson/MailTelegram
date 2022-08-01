package telegrambot

import (
	"YadnexTelegram/internal/configs"
	"YadnexTelegram/internal/services"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type State string

const (
	Connect State = State(InitCommand)
	Filters State = "filters"
	Listen  State = "listen"
)

type UserState struct {
	State State
	mail  services.ConsumerService
}

func NewUserState() *UserState {
	return &UserState{
		State: Connect,
	}
}

func (u *UserState) Connect(config configs.PostConfig, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	c := services.NewPostService(config)
	u.mail = c
	h := &telegramHandler{Bot: bot, Update: update}
	go c.Listen(h)
	u.State = Listen
}

type telegramHandler struct {
	Bot    *tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (t *telegramHandler) Handle(message services.PostMessage) {
	msg := tgbotapi.NewMessage(t.Update.Message.Chat.ID, fmt.Sprintf("Mail: \n %s \n %s", message.MailAddr,
		message.Subject))
	msg.ReplyToMessageID = t.Update.Message.MessageID
	if _, err := t.Bot.Send(msg); err != nil {
		log.Println(err)
	}
}
