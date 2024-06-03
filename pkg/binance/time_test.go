package binance

import "testing"

func TestTime(t *testing.T) {
	time, err := time("testnet.binance.vision")
	if err != nil {
		t.Error(err)
	}

	if time == 0 {
		t.Error(err)
	}
}
