package gobinance

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestPlace(t *testing.T) {
	var keysMap map[string]string
	keys := []string{"SECRET_KEY", "API_KEY"}
	keysMap, err := godotenv.Read()
	if err != nil {
		t.Error(err)
	}
	for _, k := range keys {
		if keysMap[k] == "" {
			t.Errorf("%s not found in .env file\n", k)
		}
	}

	// test sell
	options := map[string]string{"side": "SELL", "symbol": "ETHUSDT", "quantity": "0.01000", "type": "MARKET"}
	order, err := place("testnet.binance.vision", keysMap["SECRET_KEY"], keysMap["API_KEY"], options)
	if err != nil {
		t.Error(err)
	}
	if order.Status != "FILLED" {
		t.Errorf("want %q, got %q", "FILLED", order.Status)
	}

	// test buy
	options = map[string]string{"side": "BUY", "symbol": "ETHUSDT", "quantity": "0.01000", "type": "MARKET"}
	order, err = place("testnet.binance.vision", keysMap["SECRET_KEY"], keysMap["API_KEY"], options)
	if err != nil {
		t.Error(err)
	}
	if order.Status != "FILLED" {
		t.Errorf("want %q, got %q", "FILLED", order.Status)
	}
}

/*
func TestFormOrderRequest(t *testing.T) {
	m := map[string]string{"foo": "bar", "baz": "buzz", "fizz": "fuzz"}
	mm := formOrderRequest(m)
	mmm := "foo=bar&baz=buzz&fizz=fuzz"
	if mm != mmm {
		t.Errorf("want %q, got %q", mmm, mm)
	}
}
*/
