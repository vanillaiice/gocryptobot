package gobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ResponseAccount struct {
	Balances []struct {
		Asset string `json:"asset"`
		Free  string `json:"free"`
	} `json:"balances"`
}

func info(baseURL, secretKey, apiKey string) (ResponseAccount, error) {
	var response ResponseAccount
	client := &http.Client{}
	timestamp, err := time(baseURL)
	if err != nil {
		return response, err
	}
	req := fmt.Sprintf("timestamp=%d", timestamp)
	sig := sign([]byte(secretKey), []byte(req))
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/api/v3/account?%s&signature=%s", baseURL, req, sig), nil)
	if err != nil {
		return response, err
	}
	request.Header.Add("X-MBX-APIKEY", apiKey)
	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	return response, nil
}

func infoPretty(baseURL, secretKey, apiKey string) (string, error) {
	info, err := info(baseURL, secretKey, apiKey)
	if err != nil {
		return "", err
	}

	var infoPretty []string
	for _, b := range info.Balances {
		infoPretty = append(infoPretty, fmt.Sprintf("%s: %s", b.Asset, b.Free))
	}
	return strings.Join(infoPretty, "\n"), nil
}
