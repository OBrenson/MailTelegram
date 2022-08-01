package a_test

import (
	"YadnexTelegram/internal/telegrambot"
	"YadnexTelegram/internal/util"
	"testing"
)

func TestTelegramInit(t *testing.T) {
	telegrambot.ManageTelegramBot(util.GetTelConfig("./../resources/configs.yaml"))
}
