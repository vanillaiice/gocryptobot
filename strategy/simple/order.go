package simple

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vanillaiice/gocryptobot/pkg/binance"
	"github.com/vanillaiice/gocryptobot/pkg/db"
)

// tryBuy tries to place a buy order.
// It checks if the percentage change between the last sell price
// and the current price is greater than the specified percentage change.
func tryBuy(botCfg *BotCfg, state *State, client *binance.Client, price float32) (err error) {
	pChange := percentChange(state.Data.LastSellPrice, price)

	if botCfg.StopEntryPrice != 0 && !shouldBuyStopEntry(botCfg, pChange) {
		log.Debug().Msgf(
			"NOT buying, percent change between price and stop entry price @%.5f",
			pChange,
		)
		return
	} else if !shouldBuy(botCfg, pChange) {
		log.Debug().Msgf(
			"NOT buying price @%.5f %s/%s, percent change @%.5f",
			price,
			botCfg.Base,
			botCfg.Quote,
			pChange,
		)
		return
	}

	return orderMarket(botCfg, state, client, orderSideBuy)
}

// trySell tries to place a sell order.
// It checks if the percentage change between the last buy price
// and the current price is greater than the specified percentage change.
func trySell(botCfg *BotCfg, state *State, client *binance.Client, price float32) (err error) {
	pChange := percentChange(state.Data.LastBuyPrice, price)

	if botCfg.TrailingStopLossMargin != 0 && shouldSellStopLoss(botCfg, pChange) {
		logger.Warn().Msgf(
			"Triggering trailing stop loss order @%.5f %s/%s, percent change @%.5f, stop loss @%.5f",
			price,
			botCfg.Base,
			botCfg.Quote,
			pChange,
			botCfg.TrailingStopLossMargin,
		)
	} else if !shouldSell(botCfg, pChange) {
		log.Debug().Msgf(
			"NOT selling price @%.5f %s/%s, percent change @%.5f",
			price,
			botCfg.Base,
			botCfg.Quote,
			pChange,
		)
		return
	}

	return orderMarket(botCfg, state, client, orderSideSell)
}

// shouldBuy checks if the bot should buy. It compares the percentage change
// between the last sell price and the current price.
func shouldBuy(botCfg *BotCfg, pChange float32) bool {
	return pChange >= botCfg.PercentChangeBuy
}

// shouldBuyStopEntry checks if the bot should buy if a stop entry price is set.
// It compares the percentage change between the stop entry price and the current price.
func shouldBuyStopEntry(botCfg *BotCfg, price float32) bool {
	pChange := percentChange(botCfg.StopEntryPrice, price)
	if pChange < 0 {
		pChange = -pChange
	}

	return pChange <= botCfg.StopEntryPriceMargin
}

// shouldSell checks if the bot should sell. It compares the percentage change
// between the last buy price and the current price.
func shouldSell(botCfg *BotCfg, pChange float32) bool {
	return pChange <= -botCfg.PercentChangeSell
}

// shouldSellStopLoss checks if the bot should sell if a trailing stop loss is set.
// It compares the percentage change between the last buy price and the current price.
func shouldSellStopLoss(botCfg *BotCfg, pChange float32) bool {
	return pChange >= botCfg.TrailingStopLossMargin
}

// orderMarket places a market order.
func orderMarket(botCfg *BotCfg, state *State, client *binance.Client, side OrderSide) (err error) {
	tempStateData := *state.Data

	orderResponse, err := client.Order(map[string]string{
		"side":     string(side),
		"symbol":   symbolUpper,
		"quantity": fmt.Sprintf("%f", roundStepSize(tempStateData.TradingBalance, tempStateData.SymbolStepSize)),
		"type":     "MARKET",
	})

	if err != nil {
		return
	}

	if orderResponse.Status != "FILLED" {
		msg := fmt.Sprintf("code: %d, msg: %s", orderResponse.Code, orderResponse.Msg)
		return errors.New(msg)
	}

	txPrice, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 32)
	if err != nil {
		return
	}

	qty, err := strconv.ParseFloat(orderResponse.ExecutedQty, 32)
	if err != nil {
		return
	}

	var profit float32

	switch side {
	case orderSideBuy:
		profit = 0

		tempStateData.LastTx = txBuy
		tempStateData.LastBuyPrice = float32(txPrice)

		log.Info().Msgf(
			"BOUGHT %.5f %s @%.5f %s",
			tempStateData.TradingBalance,
			botCfg.Base,
			txPrice,
			tempStateData.Symbol,
		)
	case orderSideSell:
		if tempStateData.LastBuyPrice != 0 {
			profit = (float32(txPrice) - tempStateData.LastBuyPrice) * float32(qty)
		} else {
			profit = 0
		}

		tempStateData.LastTx = txSell
		tempStateData.LastSellPrice = float32(txPrice)

		if botCfg.AdjustTradingBalanceProfit && profit > 0 {
			tempStateData.TradingBalance += profit
		} else if botCfg.AdjustTradingBalanceLoss && profit < 0 {
			tempStateData.TradingBalance += profit
		}

		log.Info().Msgf(
			"SOLD %.5f %s @%.5f %s, PROFIT %.5f",
			tempStateData.TradingBalance,
			botCfg.Base,
			txPrice,
			tempStateData.Symbol,
			profit,
		)
	default:
		return fmt.Errorf("invalid side: %s", side)
	}

	receipt := &db.Receipt{
		Symbol:    symbolUpper,
		TxType:    string(side),
		Quantity:  float32(qty),
		Price:     float32(txPrice),
		Profit:    float32(profit),
		Timestamp: time.Now().Unix(),
	}

	if botCfg.StopAfterTx >= 0 && tempStateData.StopAfterTx > 0 {
		tempStateData.StopAfterTx--
	}

	state.Data = &tempStateData
	if err = state.Save(); err != nil {
		return
	}

	if botCfg.LogDb {
		err = receiptsDb.InsertReceipt(receipt)
	}

	return
}

// percentChange returns the absolute value of the percent change between a and b.
func percentChange(a, b float32) float32 {
	return ((a - b) / a) * 100
}

// roundStepSize rounds a given number to a specific step size.
// Ported from python-binance by Sam McHardy, MIT License.
func roundStepSize(num float32, stepSize float32) float64 {
	return float64(num) - math.Mod(float64(num), float64(stepSize))
}
