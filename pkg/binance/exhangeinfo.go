package gobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ResponseExchangeInfo struct {
	Timezone        string            `json:"timezone"`
	ServerTime      int64             `json:"serverTime"`
	RateLimits      []RateLimits      `json:"rateLimits"`
	ExchangeFilters []ExchangeFilters `json:"exchangeFilters"`
	Symbols         []Symbols         `json:"symbols"`
}

type RateLimits struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type ExchangeFilters struct{}

type Symbols struct {
	Symbol                          string    `json:"symbol"`
	Status                          string    `json:"status"`
	BaseAsset                       string    `json:"baseAsset"`
	BaseAssetPrecision              int       `json:"baseAssetPrecision"`
	QuoteAsset                      string    `json:"quoteAsset"`
	QuotePrecision                  int       `json:"quoteAssetPrecision"`
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

func exchangeInfo(baseURL string, symbols ...string) (ResponseExchangeInfo, error) {
	var response ResponseExchangeInfo
	req := ""
	if len(symbols) >= 1 {
		req = "?" + formSymbolArray(symbols)
	}
	resp, err := http.Get(fmt.Sprintf("https://%s/api/v3/exchangeInfo%s", baseURL, req))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func formSymbolArray(symbols []string) string {
	symbs := "["
	for _, s := range symbols {
		symbs += fmt.Sprintf("%q,", s)
	}
	symbs = strings.TrimRight(symbs, ",")
	symbs += "]"
	return symbs
}

func stepSize(baseURL string, symbols ...string) (map[string]string, error) {
	var response map[string]string
	exchangeInfo, err := exchangeInfo(baseURL, symbols...)
	if err != nil {
		return response, err
	}

	for _, i := range exchangeInfo.Symbols {
		response[i.Symbol] = i.Filters[1].StepSize
	}
	return response, nil
}
