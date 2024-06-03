package binance

import (
	"fmt"
	"testing"
)

func TestInfo(t *testing.T) {
	info, err := info(c)
	if err != nil {
		t.Error(err)
	}

	if fmt.Sprintf("%v", info) == "" {
		t.Error("Unexpected error")
	}
}

func TestInfoPretty(t *testing.T) {
	info, err := infoPretty(c)
	if err != nil {
		t.Error(err)
	}

	if info == "" {
		t.Error("Unexpected error")
	}
}
