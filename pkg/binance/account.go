package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseAccount is the response from the account endpoint.
type ResponseAccount struct {
	Balances []struct {
		Asset string `json:"asset"`
		Free  string `json:"free"`
	} `json:"balances"`
}

// info returns the account information.
func info(c *Client) (*ResponseAccount, error) {
	var response ResponseAccount
	client := &http.Client{}

	timestamp, err := time(c.BaseURL)
	if err != nil {
		return nil, err
	}

	req := fmt.Sprintf("timestamp=%d", timestamp)
	sig := sign([]byte(c.SecretKey), []byte(req))

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/api/v3/account?%s&signature=%s", c.BaseURL, req, sig), nil)
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

	err = json.Unmarshal(body, &response)

	return &response, err
}

// infoPretty returns the account information in a pretty format.
func infoPretty(c *Client) (s string, err error) {
	info, err := info(c)
	if err != nil {
		return
	}

	for _, b := range info.Balances {
		s += fmt.Sprintf("%s: %s\n", b.Asset, b.Free)
	}

	return
}
