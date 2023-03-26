package solana_client

import (
	"fmt"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/hdwallet"
	"github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

// Predefined mnemonic lengths
const (
	MnemonicLength12 MnemonicLength = 128 // 128 bits of entropy
	MnemonicLength24 MnemonicLength = 256 // 256 bits of entropy
)

// Mnemonic length type
type MnemonicLength int

// NewMnemonic generates a new mnemonic phrase
func NewMnemonic(len MnemonicLength) (string, error) {
	entropy, err := bip39.NewEntropy(int(len))
	if err != nil {
		return "", fmt.Errorf("failed to create new entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to create new mnemonic: %w", err)
	}

	return mnemonic, nil
}

// DeriveAccountFromMnemonicBip44 derives an Solana account from a mnemonic phrase
// Compatible with BIP44 (phantom wallet, etc.)
func DeriveAccountFromMnemonicBip44(mnemonic string) (types.Account, error) {
	acc, err := deriveFromMnemonicBip44(mnemonic, 0)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive account from mnemonic: %w", err)
	}

	return acc, nil
}

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

// deriveFromMnemonicBip44 derives an Solana account from a mnemonic phrase
// Compatible with BIP44 (phantom wallet)
func deriveFromMnemonicBip44(mnemonic string, path int) (types.Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create seed from mnemonic: %w", err)
	}

	derivedKey, err := hdwallet.Derived(fmt.Sprintf("m/44'/501'/%d'/0'", path), seed)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive key from seed: %w", err)
	}

	account, err := types.AccountFromSeed(derivedKey.PrivateKey)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create account from seed: %w", err)
	}

	return account, nil
}