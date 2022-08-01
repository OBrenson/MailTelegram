package a_test

import (
	configs "YadnexTelegram/internal/configs"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testing"
)

func TestConfParsing(t *testing.T) {
	p := GetPostConfig("./resources/configs.yaml")
	fmt.Println(p)

	tel := GetTelConfig("./resources/configs.yaml")
	fmt.Println(tel)

}
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
