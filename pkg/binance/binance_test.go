package binance

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const ApiUrl = "testnet.binance.vision"

var keysMap = map[string]string{}

var c *Client

func TestMain(m *testing.M) {
	var err error

	keys := []string{"SECRET_KEY", "API_KEY"}

	keysMap, err = godotenv.Read("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range keys {
		if keysMap[k] == "" {
			log.Fatalf("%s not found in .env file\n", k)
		}
	}

	c = &Client{
		BaseURL:   ApiUrl,
		SecretKey: keysMap[keys[0]],
		ApiKey:    keysMap[keys[1]],
	}

	os.Exit(m.Run())
}
