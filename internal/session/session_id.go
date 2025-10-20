package session

import (
	"crypto/rand"
	"encoding/hex"
)

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
