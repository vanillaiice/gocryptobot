package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// sign signs the message with the secret key usng HMAC-SHA256.
func sign(secretKey, message []byte) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}
