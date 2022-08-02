package telegrambot

import (
	"YadnexTelegram/internal/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"sync"
)

type TelegramInitService struct {
	Config configs.TelegramConfig
}

const (
	InitCommand   string = "/connectmail"
	ListenCommand        = "/listen"
)

type lockingPool struct {
	pool  map[int64]*UserState
	mutex sync.Mutex
}

type botMsg struct {
	chatId int64
	msg    string
}

func ManageTelegramBot(config configs.TelegramConfig) {
	bot, err := tgbotapi.NewBotAPI(config.BotName)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	lp := lockingPool{
		pool: map[int64]*UserState{},
	}

	msgQueue := make(chan botMsg)
	go sendMsg(bot, msgQueue)
	for update := range updates {
		if update.Message != nil {
			lp.mutex.Lock()
			if us, ok := lp.pool[update.Message.From.ID]; ok {
				switch us.State {
				case Connect:
					mailData := strings.Split(update.Message.Text, ", ")
					us.Connect(configs.PostConfig{
						mailData[0],
						mailData[1],
						mailData[2],
					})
					msgQueue <- botMsg{chatId: update.Message.Chat.ID,
						msg: "Expecting emails addresses in format: addr1@mail.com, addr2@mail.com, ..."}
				case Filters:
					addrs := strings.Split(update.Message.Text, ", ")
					if len(addrs) != 0 {
						us.Filter(addrs)
						msgQueue <- botMsg{chatId: update.Message.Chat.ID,
							msg: "Filters are created, expecting /listen command"}
					} else {
						msgQueue <- botMsg{chatId: update.Message.Chat.ID,
							msg: "Expecting emails addresses in format: addr1@mail.com, addr2@mail.com, ..."}
					}
				case Listen:
					if update.Message.Text == ListenCommand {
						us.Listen(update.Message.Chat.ID, &msgQueue)
						msgQueue <- botMsg{chatId: update.Message.Chat.ID,
							msg: "Now your mail is listening"}
					}
				}
			} else {
				if update.Message.Text == InitCommand {
					lp.pool[update.Message.From.ID] = NewUserState()
					msgQueue <- botMsg{chatId: update.Message.Chat.ID,
						msg: "Expecting mail address in format: address, login, password"}
				} else {
					msgQueue <- botMsg{chatId: update.Message.Chat.ID,
						msg: "In's need to be executed " + InitCommand}
				}
			}
			lp.mutex.Unlock()
		}
	}
}

func sendMsg(bot *tgbotapi.BotAPI, msgQue chan botMsg) {
	for bm := range msgQue {
		msg := tgbotapi.NewMessage(bm.chatId, bm.msg)
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
