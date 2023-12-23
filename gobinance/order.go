package gobinance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

func place(baseURL, secretKey, apiKey string, options map[string]string) (ResponseOrder, error) {
	var response ResponseOrder
	client := &http.Client{}
	timestamp, err := time(baseURL)
	if err != nil {
		return response, err
	}
	options["timestamp"] = fmt.Sprintf("%d", timestamp)
	req := formOrderRequest(options)
	sig := sign([]byte(secretKey), []byte(req))
	data := []byte(fmt.Sprintf("%s&signature=%s", req, sig))
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/api/v3/order", baseURL), bytes.NewBuffer(data))
	if err != nil {
		return response, err
	}
	request.Header.Add("X-MBX-APIKEY", apiKey)
	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		return response, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	return response, nil
}

func formOrderRequest(m map[string]string) string {
	var request string
	for k, v := range m {
		request += fmt.Sprintf("%s=%s&", k, v)
	}
	return strings.TrimRight(request, "&")
}
