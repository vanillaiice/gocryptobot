/*
Copyright © 2024 vanillaiice <vanillaiice1@proton.me>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanillaiice/gocryptobot/strategy/simple"
)

// simpleCmd represents the simple command
var simpleCmd = &cobra.Command{
	Use:   "simple",
	Short: "Use the simple trading strategy",
	Long:  `The simple trading strategy uses the Binance API and Websocket to execute trades. The bot will buy and sell an asset based on the percentage change between the current price and the last price, specified by the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		var keys = []string{"SECRET_KEY", "API_KEY"}

		env, err := cmd.Flags().GetString("env")
		if err != nil {
			log.Fatal(err)
		}

		keysMap, err := godotenv.Read(env)
		if err != nil {
			log.Fatal(err)
		}

		for _, k := range keys {
			if keysMap[k] == "" {
				log.Fatalf("%s not found in .env file", k)
			}
		}

		requiredFlags := []string{
			"base",
			"quote",
			"percent-change-buy",
			"percent-change-sell",
			"trading-balance",
		}

		for _, flag := range requiredFlags {
			if !viper.IsSet(flag) {
				log.Fatalf("required flag %s not set", flag)
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		if viper.IsSet("save-receipt") {
			for _, flag := range []string{"db-backend", "db-url"} {
				if !viper.IsSet(flag) {
					log.Fatalf("required flag %s not set", flag)
				}
			}
		}

		botCfg := simple.BotCfg{}

		botCfg.Base = viper.GetString("base")
		botCfg.Quote = viper.GetString("quote")
		botCfg.PercentChangeBuy = float32(viper.GetFloat64("percent-change-buy"))
		botCfg.PercentChangeSell = float32(viper.GetFloat64("percent-change-sell"))
		botCfg.TrailingStopLossMargin = float32(viper.GetFloat64("trailing-stop-loss-margin"))
		botCfg.StopEntryPrice = float32(viper.GetFloat64("stop-entry-price"))
		botCfg.StopEntryPriceMargin = float32(viper.GetFloat64("stop-entry-price-margin"))
		botCfg.AdjustTradingBalanceProfit = viper.GetBool("adjust-trading-balance-profit")
		botCfg.AdjustTradingBalanceLoss = viper.GetBool("adjust-trading-balance-loss")
		botCfg.StopAfterTx = viper.GetInt("stop-after-tx")
		botCfg.MaxPriceTimestampMs = viper.GetInt("max-price-timestamp-ms")
		botCfg.TradingBalance = float32(viper.GetFloat64("trading-balance"))
		botCfg.ServerBaseEndpoint = viper.GetString("server-base-endpoint")
		botCfg.WsServerBaseEndpoint = viper.GetString("ws-server-base-endpoint")
		botCfg.StateFile = viper.GetString("state-file")
		botCfg.LogDb = viper.GetBool("save-receipt")
		botCfg.DbBackend = simple.DbBackend(viper.GetString("db-backend"))
		botCfg.DbUrl = viper.GetString("db-url")
		botCfg.FirstTx = simple.Tx(viper.GetString("first-tx"))
		botCfg.ApiKey = keysMap["API_KEY"]
		botCfg.SecretKey = keysMap["SECRET_KEY"]
		botCfg.LogLevel = viper.GetString("log-level")

		log.Fatal(simple.Run(&botCfg))
	},
}

func init() {
	rootCmd.AddCommand(simpleCmd)

	simpleCmd.Flags().StringP("base", "b", "", "base currency")
	simpleCmd.Flags().StringP("quote", "q", "", "quote currency")
	simpleCmd.Flags().Float32P("percent-change-buy", "B", 0, "percent change at which the bot should buy")
	simpleCmd.Flags().Float32P("percent-change-sell", "S", 0, "percent change at which the bot should sell")
	simpleCmd.Flags().Float32P("trailing-stop-loss-margin", "r", 0, "percent change between the current price and last buy price, at which the bot will sell to limit losses")
	simpleCmd.Flags().Float32P("stop-entry-price", "p", 0, "price at which the bot should buy")
	simpleCmd.Flags().Float32P("stop-entry-price-margin", "P", 0, "minimum percent change between the current price and the stop entry price, at which the bot will activate stop entry")
	simpleCmd.Flags().BoolP("adjust-trading-balance-profit", "a", false, "add profits to the trading balance")
	simpleCmd.Flags().BoolP("adjust-trading-balance-loss", "A", false, "subtract losses from the trading balance")
	simpleCmd.Flags().IntP("stop-after-tx", "x", -1, "stop the bot after specified number of transaction")
	simpleCmd.Flags().IntP("max-price-timestamp-ms", "m", 500, "maximum valid price timestamp in milliseconds")
	simpleCmd.Flags().Float32P("trading-balance", "t", 0, "trading balance of the bot")
	simpleCmd.Flags().StringP("server-base-endpoint", "E", "api.binance.com", "binance API server base endpoint")
	simpleCmd.Flags().StringP("ws-server-base-endpoint", "w", "api.binance.com", "binance websocket server base endpoint")
	simpleCmd.Flags().StringP("state-file", "f", "state.json", "json file storing the current state of the bot")
	simpleCmd.Flags().BoolP("save-receipt", "s", false, "save transaction receipts to a database")
	simpleCmd.Flags().StringP("db-backend", "d", "sqlite", "database backend (sqlite, postgres, bbolt, redis)")
	simpleCmd.Flags().StringP("db-url", "u", "state.db", "database URL")
	simpleCmd.Flags().StringP("log-level", "l", "info", "log level (disabled, debug, info, warn, error, fatal)")
	simpleCmd.Flags().StringP("first-tx", "F", "buy", "type of first transaction (buy or sell)")

	viper.BindPFlags(simpleCmd.Flags())
}
