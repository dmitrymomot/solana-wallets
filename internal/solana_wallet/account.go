package solana_wallet

import (
	"fmt"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// NewAccount creates a new Solana account
func NewAccount() types.Account {
	return types.NewAccount()
}

// ToBase58 converts an Solana account to a base58 encoded string
func AccountToBase58(a types.Account) string {
	return base58.Encode(a.PrivateKey)
}

// FromBase58 creates an Solana account from a base58 encoded string
func AccountFromBase58(s string) (types.Account, error) {
	b, err := base58.Decode(s)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to decode base58 string to solana account: %w", err)
	}

	return types.AccountFromBytes(b)
}

// AccountFromString creates an Solana account from a base58 encoded string.
// Alias for AccountFromBase58
func AccountFromString(s string) (types.Account, error) {
	return AccountFromBase58(s)
}

// PublicKeyFromBase58 converts a base58 encoded public key to a PublicKey type
// Alias for PublicKeyFromString
func PublicKeyFromBase58(s string) common.PublicKey {
	return PublicKeyFromString(s)
}

// PublicKeyFromString converts a string to a PublicKey type
// Wrapper around the common.PublicKeyFromString function from the solana-go-sdk
func PublicKeyFromString(s string) common.PublicKey {
	return common.PublicKeyFromString(s)
}

// ValidateSolanaWalletAddr validates a Solana wallet address.
// Returns an error if the address is invalid, nil otherwise.
func ValidateSolanaWalletAddr(addr string) error {
	d, err := base58.Decode(addr)
	if err != nil {
		return fmt.Errorf("invalid wallet public key: %w", err)
	}

	if len(d) != common.PublicKeyLength {
		return fmt.Errorf("invalid wallet public key length: %w", err)
	}

	if _, err := new(edwards25519.Point).SetBytes(d); err != nil {
		return fmt.Errorf("invalid wallet public key: %w", err)
	}

	return nil
}

// DeriveTokenAccount derives an associated token account from a Solana account and a mint address.
// This is a wrapper around the FindAssociatedTokenAddress function from the solana-go-sdk.
// base58WalletAddr is the base58 encoded address of the Solana account.
// base58MintAddr is the base58 encoded address of the token mint.
// The function returns the base58 encoded address of the token account or an error.
func DeriveTokenAccount(base58WalletAddr, base58MintAddr string) (common.PublicKey, error) {
	return DeriveTokenAccountPubkey(
		PublicKeyFromBase58(base58WalletAddr),
		PublicKeyFromBase58(base58MintAddr),
	)
}

// DeriveTokenAccountPubkey derives an associated token account from a Solana account and a mint address.
// This is a wrapper around the FindAssociatedTokenAddress function from the solana-go-sdk.
func DeriveTokenAccountPubkey(wallet, mint common.PublicKey) (common.PublicKey, error) {
	ata, _, err := common.FindAssociatedTokenAddress(wallet, mint)
	if err != nil {
		return common.PublicKey{}, fmt.Errorf("failed to derive token account: %w", err)
	}

	return ata, nil
}
