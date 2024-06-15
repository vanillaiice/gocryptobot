package simple

import (
	"fmt"
	"math"
	"time"

	"github.com/vanillaiice/gocryptobot/pkg/binance"
	"github.com/vanillaiice/gocryptobot/pkg/pricews"
)

func trade(botCfg *BotCfg, state *State, client *binance.Client, c <-chan pricews.PriceData) (err error) {
	logger.Info().Msgf(
		"trading %.5f %s/%s, BUY margin @%.3f%%, SELL margin @%.3f%%",
		botCfg.TradingBalance,
		botCfg.Base,
		botCfg.Quote,
		botCfg.PercentChangeBuy,
		botCfg.PercentChangeSell,
	)

	var priceData pricews.PriceData

	for {
		priceData = <-c

		timestampDiff := math.Abs(float64(time.Now().UnixMilli() - priceData.Timestamp))
		if timestampDiff > float64(botCfg.MaxPriceTimestampMs) {
			logger.Debug().Msgf("price timestamp too old: %.3f (max: %d)", timestampDiff, botCfg.MaxPriceTimestampMs)
			continue
		}

		switch state.Data.LastTx {
		case txfirst:
			switch botCfg.FirstTx {
			case txBuy:
				if !botCfg.SkipFirstTx {
					if err = tryBuy(botCfg, state, client, -1); err != nil {
						return
					}
				} else {
					state.Data.LastTx = txBuy
					if err = state.Save(); err != nil {
						return
					}
				}
			case txSell:
				if !botCfg.SkipFirstTx {
					if err = orderMarket(botCfg, state, client, orderSideSell); err != nil {
						return
					}
				} else {
					state.Data.LastTx = txSell
					if err = state.Save(); err != nil {
						return
					}
				}
			default:
				return fmt.Errorf("invalid first tx: %s", botCfg.FirstTx)
			}
		case txBuy:
			if err = trySell(botCfg, state, client, priceData.Price*priceMultiplierSell); err != nil {
				return
			}
		case txSell:
			if err = tryBuy(botCfg, state, client, priceData.Price*priceMultiplierBuy); err != nil {
				return
			}
		default:
			return fmt.Errorf("invalid last tx: %s", state.Data.LastTx)
		}

		if botCfg.StopAfterTx >= 0 && state.Data.StopAfterTx == 0 {
			return fmt.Errorf("stopped bot after %d tx", botCfg.StopAfterTx)
		}
	}
}
