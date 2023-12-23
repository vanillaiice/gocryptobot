package simple

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/vanillaiice/gocryptobot/gobinance"
)

func tryBuy(botConfig *BotConfig, state *State, client *gobinance.Client, price float64, db *sql.DB) error {
	percentChangeBuy, err := strconv.ParseFloat(botConfig.PercentChangeBuy, 32)
	if err != nil {
		return err
	}

	if percentChange(state.LastSellPrice, price) >= percentChangeBuy {
		err := orderMarket(botConfig, state, client, gobinance.SideTypeBuy, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func trySell(botConfig *BotConfig, state *State, client *gobinance.Client, price float64, db *sql.DB) error {
	percentChangeSell, err := strconv.ParseFloat(botConfig.PercentChangeSell, 32)
	if err != nil {
		return err
	}
	if percentChange(state.LastBuyPrice, price) >= percentChangeSell {
		err := orderMarket(botConfig, state, client, gobinance.SideTypeSell, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func orderMarket(botConfig *BotConfig, state *State, client *gobinance.Client, sideType gobinance.Side, db *sql.DB) error {
	orderResponse, err := client.Order(sideType, gobinance.OrderTypeMarket, strings.ToUpper(state.Symbol), fmt.Sprintf("%.5f", state.TradingBalance))

	if err != nil {
		return err
	}

	if orderResponse.Status != "FILLED" {
		return errors.New(fmt.Sprintf("code: %d, msg: %s", orderResponse.Code, orderResponse.Msg))
	}

	txPrice, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 32)
	if err != nil {
		return err
	}

	qty, err := strconv.ParseFloat(orderResponse.ExecutedQty, 32)
	if err != nil {
		return err
	}

	var profit float64

	if sideType == gobinance.SideTypeSell {
		if state.LastBuyPrice != 0 {
			profit = (txPrice - state.LastBuyPrice) * qty
		} else {
			profit = 0
		}
	} else {
		profit = 0
	}

	if sideType == gobinance.SideTypeBuy {
		state.LastTx = TxBuy
		state.LastBuyPrice = txPrice
		log.Printf("BOUGHT %.5f %s @%.5f %s\n", state.TradingBalance, strings.ToUpper(botConfig.Base), txPrice, strings.ToUpper(state.Symbol))
	} else if sideType == gobinance.SideTypeSell {
		state.LastTx = TxSell
		state.LastSellPrice = txPrice
		log.Printf("SOLD %.5f %s @%.5f %s, PROFIT %.5f\n", state.TradingBalance, strings.ToUpper(botConfig.Base), txPrice, strings.ToUpper(state.Symbol), profit)
	} else {
		return errors.New(fmt.Sprintf("[ERROR]: unknown side type %s", sideType))
	}

	err = dumpState(state)
	if err != nil {
		return err
	}

	receipt := &Receipt{Tx(sideType), float32(qty), float32(txPrice), float32(profit), time.Now().Unix()}

	err = DBInsertReceipt(receipt, db)
	if err != nil {
		return err
	}

	return nil
}
