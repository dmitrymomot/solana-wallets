package solana_wallet

import (
	"fmt"

	"github.com/mr-tron/base58"
)

type (
	// Client is the main struct for the Solana client
	Client struct {
		salt string
	}
)

// NewClient creates a new Solana client
func NewClient(salt string) *Client {
	if l := len(salt); l < 16 {
		panic(fmt.Sprintf("invalid salt length: %d, must be more than 16", l))
	}
	return &Client{salt: salt}
}

// GenerateMnemonic generates a new mnemonic phrase, encode it to base58 and hash with secret key.
func (c *Client) EnctyptMnemonic(mnemonic, pin string) (string, error) {
	return c.encryptMnemonic(mnemonic, pin)
}

// DecryptMnemonic decrypts a base58 encoded string with AES-256-GCM and returns the decrypted mnemonic phrase
func (c *Client) DecryptMnemonic(encrypted, pin string) (string, error) {
	return c.decryptMnemonic(encrypted, pin)
}

// signingKey returns a signing key for the given pin.
// The resulting key is used to encrypt and decrypt mnemonic phrases.
func (c *Client) signingKey(pin string) []byte {
	return hash([]byte(pin + c.salt))
}

// encryptMnemonic encrypts a mnemonic phrase with AES-256-GCM and returns the encrypted data as a base58 encoded string
func (c *Client) encryptMnemonic(mnemonic, pin string) (string, error) {
	encrypted, err := encrypt([]byte(mnemonic), c.signingKey(pin))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	return base58.Encode(encrypted), nil
}

// decryptMnemonic decrypts a base58 encoded string with AES-256-GCM and returns the decrypted mnemonic phrase
func (c *Client) decryptMnemonic(encrypted, pin string) (string, error) {
	decoded, err := base58.Decode(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base58 string: %w", err)
	}

	decrypted, err := decrypt(decoded, c.signingKey(pin))
	if err != nil {
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	return string(decrypted), nil
}
