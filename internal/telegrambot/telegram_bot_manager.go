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
					sendMsg("Expecting emails addresses ", bot, update)
				case Filters:

					us.Filter(strings.Split(update.Message.Text, ", "))
					sendMsg("Filters are created, expecting /listen command", bot, update)
				case Listen:
					if update.Message.Text == ListenCommand {
						us.Listen(bot, update)
						sendMsg("Now your mail is listening", bot, update)
					}
				}
			} else {
				if update.Message.Text == InitCommand {
					lp.pool[update.Message.From.ID] = NewUserState()
					sendMsg("Expecting mail address in format: address, login, password", bot, update)
				} else {
					sendMsg("In's need to be executed "+InitCommand, bot, update)
				}
			}
			lp.mutex.Unlock()
		}
	}
}

func sendMsg(mes string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, mes)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
