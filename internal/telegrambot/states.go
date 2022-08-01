package telegrambot

import (
	"YadnexTelegram/internal/configs"
	"YadnexTelegram/internal/services"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type State string

const (
	Connect State = State(InitCommand)
	Filters State = "filters"
	Listen  State = "listen"
)

type UserState struct {
	State   State
	mail    services.ConsumerService
	filters []Filter
}

func NewUserState() *UserState {
	return &UserState{
		State: Connect,
	}
}

func (u *UserState) Connect(config configs.PostConfig) {
	c := services.NewPostService(config)
	u.mail = c
	u.State = Filters
}

func (u *UserState) Filter(addrs []string) {
	filters := make([]Filter, 0)
	for _, addr := range addrs {
		filters = append(filters, Filter{Addr: addr})
	}
	u.filters = filters
	u.State = Listen
}

func (u *UserState) Listen(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	h := &telegramHandler{Bot: bot, Update: update, filters: u.filters}
	go u.mail.Listen(h)
}

type telegramHandler struct {
	Bot     *tgbotapi.BotAPI
	Update  tgbotapi.Update
	filters []Filter
}

func (t *telegramHandler) Handle(message services.PostMessage) {
	if t.containMails(message.MailAddr) {
		msg := tgbotapi.NewMessage(t.Update.Message.Chat.ID, fmt.Sprintf("Mail: \n %s \n %s", message.MailAddr,
			message.Subject))
		if _, err := t.Bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func (t *telegramHandler) containMails(mails string) bool {
	for _, f := range t.filters {
		if strings.Contains(mails, f.Addr) {
			return true
		}
	}
	return false
}
