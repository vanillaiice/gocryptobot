package pricews

import (
	"strconv"

	ws "github.com/gorilla/websocket"
)

// PriceResponse is the response from the price endpoint.
type PriceResponse struct {
	EventType           string `json:"e"` // Event type
	EventTime           int64  `json:"E"` // Event time
	Symbol              string `json:"s"` // Symbol
	PriceChange         string `json:"p"` // Price change
	PriceChangePercent  string `json:"P"` // Price change percent
	LastPrice           string `json:"c"` // Last price
	StatisticsCloseTime int    `json:"C"` // Statistics close time
}

// PriceData holds the price at a specified timestamp.
type PriceData struct {
	Price     float32 // price
	Timestamp int64   // timestamp
}

// Listen listens for price updates on the price endpoint.
func listen(conn *ws.Conn, symbol string, c chan<- PriceData) (err error) {
	errChan := make(chan error)

	go func() {
		errChan <- readConn(conn, symbol, c)
	}()

	return <-errChan
}

// readConn reads from the price endpoint.
func readConn(conn *ws.Conn, symbol string, c chan<- PriceData) (err error) {
	var priceResponse PriceResponse

	for {
		err = conn.ReadJSON(&priceResponse)
		if err != nil {
			return err
		}

		p, err := strconv.ParseFloat(priceResponse.LastPrice, 32)
		if err != nil {
			return err
		}

		c <- PriceData{
			Price:     float32(p),
			Timestamp: priceResponse.EventTime,
		}
	}
}

// Run starts the websocket connection to the price endpoint.
func Run(endpoint, symbol string, c chan<- PriceData) (err error) {
	conn, _, err := ws.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@ticker", nil)
	if err != nil {
		return
	}
	defer conn.Close()

	return listen(conn, symbol, c)
}
