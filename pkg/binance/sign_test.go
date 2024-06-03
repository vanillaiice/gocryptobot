package binance

import "testing"

func TestSign(t *testing.T) {
	sig := sign([]byte("foobarbaz69420"), []byte("fizzbuzz"))

	if sig != "a914ce04c4ec17de64dc0c6132dffaa38724ead867efb52e36d416b6be10d177" {
		t.Error("Unexpected error")
	}
}
