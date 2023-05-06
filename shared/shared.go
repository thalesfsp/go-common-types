package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateHash returns a sha256 hash of the value.
func GenerateHash[T any](value T) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", value)))

	return hex.EncodeToString(hash[:])
}
