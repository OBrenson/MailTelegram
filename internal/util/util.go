package util

import (
	"YadnexTelegram/internal/configs"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

//Load configs from yaml file
func GetPostConfig(postConfPath string) configs.PostConfig {
	c := parseRawData(postConfPath)["config"]
	return configs.PostConfig{
		Addr:  c["addr"],
		Login: c["login"],
		Pass:  c["pass"],
	}
}

func GetTelConfig(path string) configs.TelegramConfig {
	c := parseRawData(path)["config"]
	return configs.TelegramConfig{
		BotName: c["token"],
	}
}

func parseRawData(path string) map[string]map[string]string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]map[string]string)

	if err = yaml.Unmarshal(data, &m); err != nil {
		log.Fatal(err)
	}
	return m
}
