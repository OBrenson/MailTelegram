package main

import (
	"YadnexTelegram/internal/telegrambot"
	"YadnexTelegram/internal/util"
)

func main() {
	telegrambot.ManageTelegramBot(util.GetTelConfig("resources/configs.yaml"))
}
