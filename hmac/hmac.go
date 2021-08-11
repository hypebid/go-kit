package hype_hmac

import (
	"crypto/hmac"
	"crypto/sha256"
)

func SignHmacMessage(msg string, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	macBytes := mac.Sum(nil)

	return string(macBytes), nil
}
