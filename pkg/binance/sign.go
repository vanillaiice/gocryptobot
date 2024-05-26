package gobinance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func sign(secretKey, message []byte) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}
