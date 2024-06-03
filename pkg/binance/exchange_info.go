package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ResponseExchangeInfo is the response from the exchange info endpoint.
type ResponseExchangeInfo struct {
	Timezone        string            `json:"timezone"`
	ServerTime      int64             `json:"serverTime"`
	RateLimits      []RateLimits      `json:"rateLimits"`
	ExchangeFilters []ExchangeFilters `json:"exchangeFilters"`
	Symbols         []Symbols         `json:"symbols"`
}

// RateLimits is the rate limits for the exchange info endpoint.
type RateLimits struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// ExchangeFilters is the exchange filters for the exchange info endpoint.
type ExchangeFilters struct{}

// Symbols is the symbols for the exchange info endpoint.
type Symbols struct {
	Symbol                          string    `json:"symbol"`
	Status                          string    `json:"status"`
	BaseAsset                       string    `json:"baseAsset"`
	BaseAssetPrecision              int       `json:"baseAssetPrecision"`
	QuoteAsset                      string    `json:"quoteAsset"`
	QuotePrecision                  int       `json:"quotePrecision"`
	QuoteAssetPrecision             int       `json:"quoteAssetPrecision"`
	OrderTypes                      []string  `json:"orderTypes"`
	IcebergAllowed                  bool      `json:"icebergAllowed"`
	OcoAllowed                      bool      `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed      bool      `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool      `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool      `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed            bool      `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool      `json:"isMarginTradingAllowed"`
	Filters                         []Filters `json:"filters"`
	Permissions                     []string  `json:"permissions"`
	DefaultSelfTradePreventionMode  string    `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string  `json:"allowedSelfTradePreventionModes"`
}

// Filters is the filters for the exchange info endpoint.
type Filters struct {
	FilterType            string `json:"filterType"`
	PriceFilter           string `json:"priceFilter"`
	MinPrice              string `json:"minPrice"`
	MaxPrice              string `json:"maxPrice"`
	TickSize              string `json:"tickSize"`
	MinQty                string `json:"minQty"`
	MaxQty                string `json:"maxQty"`
	StepSize              string `json:"stepSize"`
	Limit                 int    `json:"limit"`
	MinTrailingAboveDelta int    `json:"minTrailingAboveDelta"`
	MaxTrailingAboveDelta int    `json:"maxTrailingAboveDelta"`
	MinTrailingBelowDelta int    `json:"minTrailingBelowDelta"`
	MaxTrailingBelowDelta int    `json:"maxTrailingBelowDelta"`
	BidMultiplierUp       string `json:"bidMultiplierUp"`
	BidMultiplierDown     string `json:"bidMultiplierDown"`
	AskMultiplierUp       string `json:"askMultiplierUp"`
	AskMultiplierDown     string `json:"askMultiplierDown"`
	AvgPriceMins          int    `json:"avgPriceMins"`
	MinNotional           string `json:"minNotional"`
	MaxNotional           string `json:"maxNotional"`
	ApplyMinToMarket      bool   `json:"applyMinToMarket"`
	ApplyMaxToMarket      bool   `json:"applyMaxToMarket"`
	MaxNumOrders          int    `json:"maxNumOrders"`
	MaxNumAlgoOrders      int    `json:"maxNumAlgoOrders"`
}

// exchangeInfo returns the exchange info.
func exchangeInfo(baseURL string, symbols ...string) (*ResponseExchangeInfo, error) {
	var response ResponseExchangeInfo
	var s string

	if len(symbols) > 0 {
		s = "?symbols=" + formSymbolArray(symbols...)
	}

	resp, err := http.Get(fmt.Sprintf("https://%s/api/v3/exchangeInfo%s", baseURL, s))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// stepSize returns the step size for the symbol(s).
func stepSize(baseURL string, symbols ...string) (map[string]string, error) {
	exchangeInfo, err := exchangeInfo(baseURL, symbols...)
	if err != nil {
		return nil, err
	}

	response := map[string]string{}

	for _, val := range exchangeInfo.Symbols {
		response[val.Symbol] = val.Filters[1].StepSize
	}

	return response, nil
}

// formSymbolArray forms an array of symbols from a string.
func formSymbolArray(symbols ...string) (symbs string) {
	symbs = "["

	for _, s := range symbols {
		symbs += fmt.Sprintf("%q,", s)
	}

	symbs = strings.TrimRight(symbs, ",")

	symbs += "]"

	return
}
