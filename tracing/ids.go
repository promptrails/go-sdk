package tracing

import (
	"crypto/rand"
	"encoding/hex"
)

// generateTraceID returns a 32-character hex trace ID (16 random bytes).
func generateTraceID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// generateSpanID returns a 16-character hex span ID (8 random bytes).
func generateSpanID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
