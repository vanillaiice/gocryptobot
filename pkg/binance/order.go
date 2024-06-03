package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseOrder is the response from the order endpoint.
type ResponseOrder struct {
	Code               int    `json:"code"`
	Msg                string `json:"msg"`
	Symbol             string `json:"symbol"`
	OrderId            int    `json:"orderId"`
	OrderListId        int    `json:"orderListId"`
	ClientOrderId      string `json:"clientOrderId"`
	TransactTime       int64  `json:"transactTime"`
	Price              string `json:"price"`
	OrigQty            string `json:"origQty"`
	ExecutedQty        string `json:"executedQty"`
	CumulativeQuoteQty string `json:"cumulativeQuoteQty"`
	Status             string `json:"status"`
	TimeInForce        string `json:"timeInForce"`
	Type               string `json:"type"`
	Side               string `json:"side"`
	WorkingTime        int64  `json:"workingTime"`
	Fills              []struct {
		Price                   string `json:"price"`
		Qty                     string `json:"qty"`
		Commission              string `json:"commission"`
		CommissionAsset         string `json:"commissionAsset"`
		TradeId                 int    `json:"tradeId"`
		SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	}
}

// place places an order.
func place(c *Client, options map[string]string) (*ResponseOrder, error) {
	var response ResponseOrder
	client := &http.Client{}

	timestamp, err := time(c.BaseURL)
	if err != nil {
		return nil, err
	}

	options["timestamp"] = fmt.Sprintf("%d", timestamp)

	req := formOrderRequest(options)

	sig := sign([]byte(c.SecretKey), []byte(req))

	data := []byte(fmt.Sprintf("%s&signature=%s", req, sig))

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/api/v3/order", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-MBX-APIKEY", c.ApiKey)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response, json.Unmarshal(body, &response)
}

// formOrderRequest forms the order request from its parameters.
func formOrderRequest(opts map[string]string) (s string) {
	lenOpts := len(opts)
	var i int

	for k, v := range opts {
		if i == lenOpts-1 {
			s += fmt.Sprintf("%s=%s", k, v)
		} else {
			s += fmt.Sprintf("%s=%s&", k, v)
		}

		i++
	}

	return s
}
