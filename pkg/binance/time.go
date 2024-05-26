package gobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseTime struct {
	ServerTime int64 `json:"serverTime"`
}

func time(baseURL string) (int64, error) {
	var response ResponseTime
	resp, err := http.Get(fmt.Sprintf("https://%s/api/v3/time", baseURL))
	if err != nil {
		return response.ServerTime, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.ServerTime, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response.ServerTime, err
	}
	return response.ServerTime, nil
}
