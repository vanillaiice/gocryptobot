package gobinance

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestInfo(t *testing.T) {
	var keysMap map[string]string
	keys := []string{"SECRET_KEY", "API_KEY"}
	keysMap, err := godotenv.Read()
	if err != nil {
		t.Error(err)
	}
	for _, k := range keys {
		if keysMap[k] == "" {
			t.Errorf("%s not found in .env file\n", k)
		}
	}

	info, err := info("testnet.binance.vision", keysMap["SECRET_KEY"], keysMap["API_KEY"])
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%v", info) == "" {
		t.Error("Unexpected error")
	}
}

func TestInfoPretty(t *testing.T) {
	var keysMap map[string]string
	keys := []string{"SECRET_KEY", "API_KEY"}
	keysMap, err := godotenv.Read()
	if err != nil {
		t.Error(err)
	}
	for _, k := range keys {
		if keysMap[k] == "" {
			t.Errorf("%s not found in .env file\n", k)
		}
	}

	info, err := infoPretty("testnet.binance.vision", keysMap["SECRET_KEY"], keysMap["API_KEY"])
	if err != nil {
		t.Error(err)
	}
	if info == "" {
		t.Error("Unexpected error")
	}
}
