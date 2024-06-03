package binance

import "testing"

func TestPlace(t *testing.T) {
	options := map[string]string{"side": "SELL", "symbol": "ETHUSDT", "quantity": "0.01000", "type": "MARKET"}

	order, err := place(c, options)
	if err != nil {
		t.Error(err)
	}

	if order.Status != "FILLED" {
		t.Errorf("want %q, got %q", "FILLED", order.Status)
	}

	options = map[string]string{"side": "BUY", "symbol": "ETHUSDT", "quantity": "0.01000", "type": "MARKET"}

	order, err = place(c, options)
	if err != nil {
		t.Error(err)
	}

	if order.Status != "FILLED" {
		t.Errorf("want %q, got %q", "FILLED", order.Status)
	}
}
