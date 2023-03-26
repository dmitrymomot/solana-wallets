package solana_client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

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

// hash hashes a byte slice with SHA-256 and returns the hash as a byte slice
func hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
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

// encrypt string to base64 crypto using AES-256
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt base64 crypto to decrypted string
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
