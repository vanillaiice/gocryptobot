package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseTime is the response from the time endpoint.
type ResponseTime struct {
	ServerTime int64 `json:"serverTime"`
}

// time returns the server time.
func time(baseURL string) (time int64, err error) {
	var response ResponseTime

	resp, err := http.Get(fmt.Sprintf("https://%s/api/v3/time", baseURL))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}

	return response.ServerTime, nil
}
