package solana_wallet

import "crypto/sha256"

// hash hashes a byte slice with SHA-256 and returns the hash as a byte slice
func hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}
