package pricews

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
)

type PriceResponse struct {
	Id     int `json:"id"`
	Status int `json:"status"`
	Result struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
}

/*
type WriteMsg struct {
	Id     int    `json:"id"`
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Params struct {
	Symbol string `json:"symbol"`
}
*/

const (
	SlowDown = 429
	IpBanned = 418
	Ok       = 200
)

func listen(conn *ws.Conn, symbol string, pricePtr *float64, writeMsgIntervalMs int, received ...chan int) {
	go readMsg(conn, symbol, pricePtr, received...)
	ticker := time.NewTicker(time.Duration(writeMsgIntervalMs) * time.Millisecond)
	for {
		err := writeMsg(conn, symbol)
		if err != nil {
			log.Fatalln(err)
		}
		_ = <-ticker.C
	}
}

func writeMsg(conn *ws.Conn, symbol string) error {
	err := conn.WriteMessage(ws.TextMessage, []byte(fmt.Sprintf(`{"id": 0, "method": "ticker.price", "params": {"symbol": %q}}`, symbol)))
	if err != nil {
		return err
	}
	return nil
}

func readMsg(conn *ws.Conn, symbol string, pricePtr *float64, received ...chan int) {
	var priceResponse PriceResponse
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal(msg, &priceResponse)
		if err != nil {
			log.Fatalln(err)
		}

		switch priceResponse.Status {
		case Ok:
			*pricePtr, err = strconv.ParseFloat(priceResponse.Result.Price, 64)
			if err != nil {
				log.Fatalln(err)
			}
			if len(received) != 0 {
				received[0] <- 0
			}
		case SlowDown:
			log.Printf("received code %d, slow down api calls\n", SlowDown)
			os.Exit(0)
		case IpBanned:
			log.Printf("received code %d, IP address banned\n", IpBanned)
			os.Exit(0)
		default:
			log.Fatalf("received unhandled code %d\n", priceResponse.Status)
		}
	}
}

func Start(baseEndpoint, symbol string, pricePtr *float64, writeMsgIntervalMs int, received ...chan int) {
	conn, _, err := ws.DefaultDialer.Dial(fmt.Sprintf("wss://%s/ws-api/v3", baseEndpoint), nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	listen(conn, strings.ToUpper(symbol), pricePtr, writeMsgIntervalMs, received...)
}
