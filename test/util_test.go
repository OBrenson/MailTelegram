package a_test

import (
	"YadnexTelegram/internal/util"
	"fmt"
	"testing"
)

func TestConfParsing(t *testing.T) {
	p := util.GetPostConfig("./resources/configs.yaml")
	fmt.Println(p)

	tel := util.GetTelConfig("./resources/configs.yaml")
	fmt.Println(tel)

}
