package utils

import "github.com/mr-tron/base58"

// Base58ToBytes converts base58 string to bytes.
func Base58ToBytes(s string) ([]byte, error) {
	return base58.Decode(s)
}

// BytesToBase58 converts bytes to base58 string.
func BytesToBase58(b []byte) string {
	return base58.Encode(b)
}
