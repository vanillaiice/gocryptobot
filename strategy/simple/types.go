package simple

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vanillaiice/gocryptobot/pkg/binance"
	"github.com/vanillaiice/gocryptobot/pkg/db"
)

// Tx is the type of transaction.
type Tx string

// Enum for transaction types.
const (
	txfirst Tx = "first"
	txBuy   Tx = "buy"
	txSell  Tx = "sell"
)

// OrderSide is the side of order.
type OrderSide string

// Enum for order sides.
const (
	orderSideBuy  OrderSide = "BUY"
	orderSideSell OrderSide = "SELL"
)

// DbBackend is the type of database backend.
type DbBackend string

// Enum for supported database backends.
const (
	dbBackendSqlite   DbBackend = "sqlite"
	dbBackendPostgres DbBackend = "postgres"
	dbBackendBbolt    DbBackend = "bbolt"
	dbBackendRedis    DbBackend = "redis"
)

// logLevelMap is the map of log levels.
var logLevelMap = map[string]zerolog.Level{
	"disabled": zerolog.Disabled,
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"fatal":    zerolog.FatalLevel,
}

// Logger is the logger to use.
var logger = log.Logger

// State holds the state of the bot.
type State struct {
	fileName string     // file name of the state file.
	Data     *StateData // state data.
}

// StateData holds the state data.
type StateData struct {
	Symbol         string  `json:"symbol"`         // Symbol is the trading pair symbol.
	SymbolStepSize float32 `json:"stepSize"`       // SymbolStepSize is the step size of the trading pair.
	LastBuyPrice   float32 `json:"lastBuyPrice"`   // LastBuyPrice is the last buy price.
	LastSellPrice  float32 `json:"lastSellPrice"`  // LastSellPrice is the last sell price.
	TradingBalance float32 `json:"tradingBalance"` // TradingBalance is the bot's trading balance.
	StopAfterTx    int     `json:"stopAfterTx"`    // StopAfterTx stops the bot after the specified number of transaction.
	LastTx         Tx      `json:"lastTx"`         // LastTx is the last transaction type.
}

// NewState creates a new state.
func NewState(botCfg *BotCfg, client *binance.Client) (*State, error) {
	var s State

	s.fileName = botCfg.StateFile

	if data, err := os.ReadFile(s.fileName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			stepSize, err := client.StepSize(symbolUpper)
			if err != nil {
				return nil, err
			}
			stepSizeValue, ok := stepSize[symbolUpper]
			if !ok {
				return nil, fmt.Errorf("step size not found for %s", symbolUpper)
			}
			stepSizeValueFloat, err := strconv.ParseFloat(stepSizeValue, 32)
			if err != nil {
				return nil, err
			}

			s.Data = &StateData{
				Symbol:         symbolUpper,
				SymbolStepSize: float32(stepSizeValueFloat),
				LastBuyPrice:   0,
				LastSellPrice:  0,
				TradingBalance: botCfg.TradingBalance,
				StopAfterTx:    botCfg.StopAfterTx,
				LastTx:         txfirst,
			}

			if err := s.Save(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if err = json.Unmarshal(data, &s.Data); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// Save saves the state to file.
func (s *State) Save() (err error) {
	data, err := json.MarshalIndent(s.Data, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileName, data, 0644)
}

// BotCfg is the configuration for the bot.
type BotCfg struct {
	Base                       string    `json:"base"`                       // Base is the base currency.
	Quote                      string    `json:"quote"`                      // Quote is the quote currency.
	PercentChangeBuy           float32   `json:"percentChangeBuy"`           // PercentChangeBuy is the percent change at which the bot should buy.
	PercentChangeSell          float32   `json:"percentChangeSell"`          // PercentChangeSell is the percent change at which the bot should sell.
	TrailingStopLossMargin     float32   `json:"trailingStopLossMargin"`     // TrailingStopLossMargin is the percent change between the current price and last buy price, at which the bot will sell to limit losses.
	StopEntryPrice             float32   `json:"stopEntryPrice"`             // StopEntryPrice is the entry price at which the bot should buy.
	StopEntryPriceMargin       float32   `json:"stopEntryPriceMargin"`       // StopEntryPriceMargin is the minimum percent change between the current price and the stop entry price, at which the bot will buy.
	AdjustTradingBalanceProfit bool      `json:"adjustTradingBalanceProfit"` // AdjustTradingBalanceProfit adds profits to the trading balance.
	AdjustTradingBalanceLoss   bool      `json:"adjustTradingBalanceLoss"`   // AdjustTradingBalanceLoss subtracts losses from the trading balance.
	StopAfterTx                int       `json:"stopAfterTx"`                // StopAfterTx stops the bot after the specified number of transaction.
	MaxPriceTimestampMs        int       `json:"maxPriceTimestampMs"`        // MaxPriceTimestampMs is the maximum valid price timestamp in milliseconds.
	TradingBalance             float32   `json:"tradingBalance"`             // TradingBalance is the bot's trading balance.
	ServerBaseEndpoint         string    `json:"serverBaseEndpoint"`         // ServerBaseEndpoint is the Binance API base endpoint.
	WsServerBaseEndpoint       string    `json:"wsServerBaseEndpoint"`       // WsServerBaseEndpoint is the Binance Websocket base endpoint.
	StateFile                  string    `json:"stateFile"`                  // StateFile stores the state of the bot.
	LogDb                      bool      `json:"logDb"`                      // LogDb makes the bot log transaction receipts to a database.
	DbBackend                  DbBackend `json:"dbBackend"`                  // DbBackend is the database backend.
	DbUrl                      string    `json:"dbUrl"`                      // DbUrl is the database URL.
	FirstTx                    Tx        `json:"firstTx"`                    // FirstTx is the type of the first transaction.
	SkipFirstTx                bool      `json:"skipFirstTx"`                // SkipFirstTx skips the first transaction.
	ApiKey                     string    `json:"apiKey"`                     // ApiKey is the Binance API key.
	SecretKey                  string    `json:"secretKey"`                  // SecretKey is the Binance secret key.
	LogLevel                   string    `json:"logLevel"`                   // LogLevel is the log level.
}

// DB is the database to use for storing receipts.
var receiptsDb db.DB

// symbolUpper is the symbol in uppercase.
var symbolUpper string

// priceMultiplier is a multiplier that accounts for the 0.1% trading fee.
const priceMultiplier float32 = 1.001
