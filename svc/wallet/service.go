package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmitrymomot/solana-wallets/internal/solanawallet"
	"github.com/dmitrymomot/solana-wallets/internal/utils"
	wallet_repository "github.com/dmitrymomot/solana-wallets/svc/wallet/repository"
	"github.com/portto/solana-go-sdk/types"
)

type (
	// Service interface
	Service interface {
		// Generate new wallet
		GenerateWallet(ctx context.Context) (Wallet, error)
		// Store wallet
		StoreWallet(ctx context.Context, uid, pin, mnemonic, name string) error
		// Get wallet by user id
		GetWallet(ctx context.Context, uid string) (Wallet, error)
		// Delete wallet by user id
		DeleteWallet(ctx context.Context, uid string, pin string) error
		// Update wallet name
		UpdateWalletName(ctx context.Context, uid string, pin string, name string) error
		// Change wallet pin
		ChangeWalletPin(ctx context.Context, uid string, pin string, newPin string) error
		// Export wallet
		ExportWallet(ctx context.Context, uid string, pin string) (Wallet, error)
		// Sign transaction and return signed transaction as base64 string
		SignTransaction(ctx context.Context, uid string, pin string, base64Tx string) (string, error)
		// Sign message and return signed message as base64 string
		SignMessage(ctx context.Context, uid string, pin string, base64Msg string) (msg, signature string, err error)
		// Sign and send transaction, return transaction signature
		SignAndSendTransaction(ctx context.Context, uid string, pin string, base64Tx string) (string, error)
	}

	// service struct
	service struct {
		repo   walletRepository
		wallet solanaWallet
		solana solanaClient
	}

	walletRepository interface {
		CreateWallet(ctx context.Context, arg wallet_repository.CreateWalletParams) (wallet_repository.Wallet, error)
		DeleteWallet(ctx context.Context, userID string) error
		GetWallet(ctx context.Context, userID string) (wallet_repository.Wallet, error)
		GetWalletByPublicKey(ctx context.Context, publicKey string) (wallet_repository.Wallet, error)
		UpdateWallet(ctx context.Context, arg wallet_repository.UpdateWalletParams) (wallet_repository.Wallet, error)
	}

	solanaWallet interface {
		EnctyptMnemonic(mnemonic, pin string) (string, error)
		DecryptMnemonic(encrypted, pin string) (string, error)
	}

	solanaClient interface {
		SignTransaction(ctx context.Context, wallet types.Account, txSource string) (string, error)
		SendTransaction(ctx context.Context, txSource string, i ...uint8) (string, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo walletRepository, wallet solanaWallet, solana solanaClient) Service {
	return &service{
		repo:   repo,
		wallet: wallet,
		solana: solana,
	}
}

// Generate new wallet
func (s *service) GenerateWallet(ctx context.Context) (Wallet, error) {
	mnemonic, err := solanawallet.NewMnemonic(solanawallet.MnemonicLength12)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	acc, err := solanawallet.DeriveAccountFromMnemonicBip44(mnemonic)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to derive account from mnemonic: %w", err)
	}

	return Wallet{
		Name:       "Wallet",
		Mnemonic:   mnemonic,
		PublicKey:  acc.PublicKey.ToBase58(),
		PrivateKey: utils.BytesToBase58(acc.PrivateKey),
	}, nil
}

// Store wallet
func (s *service) StoreWallet(ctx context.Context, uid, pin, mnemonic, name string) error {
	acc, err := solanawallet.DeriveAccountFromMnemonicBip44(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to derive account from mnemonic: %w", err)
	}

	if name == "" {
		name = "Wallet"
	}

	encrypted, err := s.wallet.EnctyptMnemonic(mnemonic, pin)
	if err != nil {
		return fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	if _, err := s.repo.CreateWallet(ctx, wallet_repository.CreateWalletParams{
		UserID:    uid,
		Name:      name,
		PublicKey: acc.PublicKey.ToBase58(),
		Mnemonic:  encrypted,
	}); err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

// Get wallet by user id
func (s *service) GetWallet(ctx context.Context, uid string) (Wallet, error) {
	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Wallet{}, ErrNotFound
		}
		return Wallet{}, fmt.Errorf("failed to get wallet: %w", err)
	}

	return Wallet{Name: w.Name, PublicKey: w.PublicKey}, nil
}

// Delete wallet by user id
func (s *service) DeleteWallet(ctx context.Context, uid string, pin string) error {
	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	mnemonic, err := s.wallet.DecryptMnemonic(w.Mnemonic, pin)
	if err != nil || mnemonic == "" {
		return ErrInvalidPIN
	}

	if err := s.repo.DeleteWallet(ctx, uid); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to delete wallet: %w", err)
		}
	}

	return nil
}

// Update wallet name
func (s *service) UpdateWalletName(ctx context.Context, uid, pin, name string) error {
	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	mnemonic, err := s.wallet.DecryptMnemonic(w.Mnemonic, pin)
	if err != nil || mnemonic == "" {
		return ErrInvalidPIN
	}

	if name == "" {
		return fmt.Errorf("name is required")
	}

	if _, err := s.repo.UpdateWallet(ctx, wallet_repository.UpdateWalletParams{
		UserID:   uid,
		Name:     name,
		Mnemonic: w.Mnemonic,
	}); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// Change wallet pin
func (s *service) ChangeWalletPin(ctx context.Context, uid string, pin string, newPin string) error {
	if pin == newPin {
		return fmt.Errorf("new pin must be different from the old one")
	}
	if newPin == "" || len(newPin) < 4 {
		return fmt.Errorf("new pin is required and must be at least 4 characters long")
	}

	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	mnemonic, err := s.wallet.DecryptMnemonic(w.Mnemonic, pin)
	if err != nil || mnemonic == "" {
		return ErrInvalidPIN
	}

	encrypted, err := s.wallet.EnctyptMnemonic(mnemonic, newPin)
	if err != nil {
		return fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	if _, err := s.repo.UpdateWallet(ctx, wallet_repository.UpdateWalletParams{
		UserID:   uid,
		Name:     w.Name,
		Mnemonic: encrypted,
	}); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// Export wallet
func (s *service) ExportWallet(ctx context.Context, uid string, pin string) (Wallet, error) {
	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Wallet{}, ErrNotFound
		}
		return Wallet{}, fmt.Errorf("failed to get wallet: %w", err)
	}

	mnemonic, err := s.wallet.DecryptMnemonic(w.Mnemonic, pin)
	if err != nil || mnemonic == "" {
		return Wallet{}, ErrInvalidPIN
	}

	acc, err := solanawallet.DeriveAccountFromMnemonicBip44(mnemonic)
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to derive account from mnemonic: %w", err)
	}

	return Wallet{
		Name:       w.Name,
		PublicKey:  acc.PublicKey.ToBase58(),
		PrivateKey: utils.BytesToBase58(acc.PrivateKey),
		Mnemonic:   mnemonic,
	}, nil
}

// Sign message and return signed message as base64 string
func (s *service) SignMessage(ctx context.Context, uid string, pin string, base64Msg string) (msg, signature string, err error) {
	acc, err := s.getAccount(ctx, uid, pin)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	decodedMsg, err := utils.Base64ToBytes(base64Msg)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	return base64Msg, utils.BytesToBase64(acc.Sign(decodedMsg)), nil
}

// Sign transaction and return signed transaction as base64 string
func (s *service) SignTransaction(ctx context.Context, uid string, pin string, base64Tx string) (string, error) {
	acc, err := s.getAccount(ctx, uid, pin)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	signedTx, err := s.solana.SignTransaction(ctx, acc, base64Tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// Sign and send transaction, return transaction signature
func (s *service) SignAndSendTransaction(ctx context.Context, uid string, pin string, base64Tx string) (string, error) {
	signedTx, err := s.SignTransaction(ctx, uid, pin, base64Tx)
	if err != nil {
		return "", err
	}

	txSignature, err := s.solana.SendTransaction(ctx, signedTx, 2)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return txSignature, nil
}

// get decoded wallet account
func (s *service) getAccount(ctx context.Context, uid string, pin string) (types.Account, error) {
	w, err := s.repo.GetWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Account{}, ErrNotFound
		}
		return types.Account{}, fmt.Errorf("failed to get wallet: %w", err)
	}

	mnemonic, err := s.wallet.DecryptMnemonic(w.Mnemonic, pin)
	if err != nil || mnemonic == "" {
		return types.Account{}, ErrInvalidPIN
	}

	acc, err := solanawallet.DeriveAccountFromMnemonicBip44(mnemonic)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to derive account from mnemonic: %w", err)
	}

	return acc, nil
}
