package simple

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vanillaiice/gocryptobot/gobinance"
)

type State struct {
	Symbol         string
	SymbolStepSize string
	LastBuyPrice   float64
	LastSellPrice  float64
	TradingBalance float64
	LastTx         Tx
}

func trade(botConfig *BotConfig, keys map[string]string, pricePtr *float64, received chan int) {
	symbol := strings.ToLower(botConfig.Base) + strings.ToLower(botConfig.Quote)
	stateFilePath := fmt.Sprintf("state/%s.json", symbol)

	var state State

	_, err := createDirIfNotExist("state")
	if err != nil {
		log.Fatalln(err)
	}

	exists, err := createFileIfNotExist(stateFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	if !exists {
		tradingBalanceFloat, err := strconv.ParseFloat(botConfig.TradingBalance, 64)
		if err != nil {
			log.Fatalln(err)
		}
		initialState := State{
			Symbol:         symbol,
			SymbolStepSize: "0",
			LastBuyPrice:   0,
			LastSellPrice:  0,
			TradingBalance: tradingBalanceFloat,
			LastTx:         TxFirst,
		}
		initialStateByte, err := json.Marshal(initialState)
		if err != nil {
			log.Fatalln(err)
		}
		err = os.WriteFile(stateFilePath, initialStateByte, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
		state = initialState
	} else {
		stateByte, err := os.ReadFile(stateFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("read state file %s\n", stateFilePath)
		err = json.Unmarshal(stateByte, &state)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var db *sql.DB

	if botConfig.LogDB {
		db, err = DBOpen(state.Symbol)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()
	}

	client := gobinance.NewClient(botConfig.ServerBaseEndpoint, keys["SECRET_KEY"], keys["API_KEY"])

	_ = <-received

	log.Printf("BOT: trading %.5f %s/%s, BUY margin @%s%%, SELL margin @%s%%, current price %.5f %s/%s\n", state.TradingBalance, botConfig.Base, botConfig.Quote, botConfig.PercentChangeBuy, botConfig.PercentChangeSell, *pricePtr, botConfig.Base, botConfig.Quote)

	for {
		switch state.LastTx {
		case TxFirst:
			switch botConfig.FirstTx {
			case TxBuy:
				err := orderMarket(botConfig, &state, client, gobinance.SideTypeBuy, db)
				if err != nil {
					log.Printf("[ERROR]: %s\n", err)
				}
			case TxSell:
				err := orderMarket(botConfig, &state, client, gobinance.SideTypeSell, db)
				if err != nil {
					log.Printf("[ERROR]: %s\n", err)
				}
			default:
				log.Fatalf("[ERROR]: unknown value %q for LastTx\n", state.LastTx)
			}
		case TxBuy:
			err := trySell(botConfig, &state, client, *pricePtr, db)
			if err != nil {
				log.Printf("[ERROR]: %s\n", err)
			}
		case TxSell:
			err := tryBuy(botConfig, &state, client, *pricePtr, db)
			if err != nil {
				log.Printf("[ERROR]: %s\n", err)
			}
		default:
			log.Fatalf("[ERROR]: unknown value %q for LastTx\n", state.LastTx)
		}

		_ = <-received
	}
}
