package a_test

import (
	"YadnexTelegram/internal/telegrambot"
	"testing"
)

func TestTelegramInit(t *testing.T) {
	telegrambot.ManageTelegramBot(GetTelConfig("./resources/configs.yaml"))
}
