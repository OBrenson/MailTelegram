package telegrambot

import (
	"YadnexTelegram/internal/configs"
	"YadnexTelegram/internal/services"
	"fmt"
	"strings"
)

type State string

//States enum
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

//Prepare UserState for connecting with mail
func (u *UserState) Connect(config configs.PostConfig) {
	c := services.NewPostService(config)
	u.mail = c
	u.State = Filters
}

//Adding mail filters
func (u *UserState) Filter(addrs []string) {
	filters := make([]Filter, 0)
	for _, addr := range addrs {
		filters = append(filters, Filter{Addr: addr})
	}
	u.filters = filters
	u.State = Listen
}

//Start listening
func (u *UserState) Listen(chatId int64, msgQue *chan botMsg) {
	h := &telegramHandler{chatId: chatId, msgQue: msgQue, filters: u.filters}
	go u.mail.Listen(h)
}

type telegramHandler struct {
	chatId  int64
	msgQue  *chan botMsg
	filters []Filter
}

func (t *telegramHandler) Handle(message services.PostMessage) {
	if t.containMails(message.MailAddr) {
		*t.msgQue <- botMsg{chatId: t.chatId,
			msg: fmt.Sprintf("Mail: \n %s \n %s", message.MailAddr, message.Subject)}
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
