package solanawallet

import (
	"fmt"

	"github.com/portto/solana-go-sdk/pkg/hdwallet"
	"github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

// Predefined mnemonic lengths
const (
	MnemonicLength12 MnemonicLength = 128 // 128 bits of entropy, 12 words
	MnemonicLength24 MnemonicLength = 256 // 256 bits of entropy, 24 words
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
