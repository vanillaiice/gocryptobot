package gobinance

import "testing"

func TestExchangeInfo(t *testing.T) {
	_, err := exchangeInfo("testnet.binance.vision", "BTCUSDT", "TRXUSDT")
	if err != nil {
		t.Error(err)
	}

	_, err = exchangeInfo("testnet.binance.vision", "BTCUSDT")
	if err != nil {
		t.Error(err)
	}

	_, err = exchangeInfo("testnet.binance.vision")
	if err != nil {
		t.Error(err)
	}
}

func TestFormSymbolArray(t *testing.T) {
	s := formSymbolArray([]string{"BTCUSDT", "TRXUSDT", "ETHUSDT"})
	if s != `["BTCUSDT","TRXUSDT","ETHUSDT"]` {
		t.Errorf("want %q, got %q", `["BTCUSDT","TRXUSDT","ETHUSDT"]`, s)
	}
}

// func testStepSize(t *testing.T) {}
