package simple

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vanillaiice/gocryptobot/pricews"
	"github.com/vanillaiice/gocryptobot/strategy/simple/config"
)

type Tx string

const (
	TxFirst Tx = "first"
	TxSell  Tx = "sell"
	TxBuy   Tx = "buy"
)

type BotConfig struct {
	Base                 string `json:"base"`
	Quote                string `json:"quote"`
	PercentChangeBuy     string `json:"percentChangeBuy"`
	PercentChangeSell    string `json:"percentChangeSell"`
	TradingBalance       string `json:"tradingBalance"`
	ServerBaseEndpoint   string `json:"serverBaseEndpoint"`
	WsServerBaseEndpoint string `json:"wsServerBaseEndpoint"`
	Testnet              bool   `json:"testnet"`
	LogDB                bool   `json:"logDB"`
	DecisionIntervalMs   int    `json:"decisionintervalMs"`
	FirstTx              Tx     `json:"firstTx"`
}

func Start(configPath string, keys map[string]string) {
	var botConfig *BotConfig
	if configPath == "" {
		var createFileAns string
		fmt.Print("no config file provided, create one ? (y/N): ")
		fmt.Scanln(&createFileAns)
		if strings.ToLower(createFileAns) == "no" || strings.ToLower(createFileAns) == "n" {
			fmt.Println("not creating config file")
			os.Exit(0)
		}
		configStr := config.New()
		err := json.Unmarshal([]byte(configStr), &botConfig)
		if err != nil {
			log.Fatalln(err)
		}
		configFile, err := os.Create(fmt.Sprintf("%s.json", fmt.Sprintf("%s%s", strings.ToLower(botConfig.Base), strings.ToLower(botConfig.Quote))))
		if err != nil {
			log.Fatalln(err)
		}
		defer configFile.Close()
		_, err = configFile.WriteString(configStr)
		if err != nil {
			log.Fatalln(err)
		}
	} else if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s does not exist\n", configPath)
		os.Exit(1)
	} else {
		configFileByte, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal(configFileByte, &botConfig)
		if err != nil {
			log.Fatalln(err)
		}
	}

	received := make(chan int, 1)
	var price float64
	go pricews.Start(botConfig.WsServerBaseEndpoint, strings.ToLower(botConfig.Base)+strings.ToLower(botConfig.Quote), &price, botConfig.DecisionIntervalMs, received)
	trade(botConfig, keys, &price, received)
}
