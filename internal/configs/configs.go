package configs

//Config for mail.
type PostConfig struct {
	Addr  string
	Login string
	Pass  string
}

//Config for telegram bot.
type TelegramConfig struct {
	BotName string
}
