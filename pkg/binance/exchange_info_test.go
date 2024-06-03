package binance

import "testing"

func TestExchangeInfo(t *testing.T) {
	_, err := exchangeInfo(ApiUrl, "BTCUSDT", "TRXUSDT")
	if err != nil {
		t.Error(err)
	}

	_, err = exchangeInfo(ApiUrl, "BTCUSDT")
	if err != nil {
		t.Error(err)
	}
}

func TestStepSize(t *testing.T) {
	stepSize, err := stepSize(ApiUrl, "BTCUSDT", "TRXUSDT")
	if err != nil {
		t.Error(err)
	}

	want := "0.00001000"
	if stepSize["BTCUSDT"] != want {
		t.Errorf("want %q, got %q", want, stepSize["BTCUSDT"])
	}

	want = "0.10000000"
	if stepSize["TRXUSDT"] != want {
		t.Errorf("want %q, got %q", want, stepSize["TRXUSDT"])
	}
}

func TestFormSymbolArray(t *testing.T) {
	s := formSymbolArray("BTCUSDT", "TRXUSDT", "ETHUSDT")

	want := `["BTCUSDT","TRXUSDT","ETHUSDT"]`

	if s != want {
		t.Errorf("want %q, got %q", want, s)
	}
}
