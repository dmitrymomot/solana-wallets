package wallet

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/solana/metadata"
	"github.com/dmitrymomot/solana/token_metadata"
	"github.com/dmitrymomot/solana/types"
)

type (
	// Service interface
	Service interface {
		// GetSOLBalance returns the SOL balance of a wallet
		GetSOLBalance(ctx context.Context, walletAddr string) (Balance, error)
		// Get balance for specified token
		GetTokenBalance(ctx context.Context, walletAddr, tokenMint string) (types.TokenAmount, error)
		// Get balance for all fungible tokens
		GetFungibleTokens(ctx context.Context, walletAddr string) ([]Balance, error)
		// Get balance for all fungible assets
		GetFungibleAssets(ctx context.Context, walletAddr string) ([]Balance, error)
		// Get all non-fungible tokens
		GetNonFungibleTokens(ctx context.Context, walletAddr string) ([]token_metadata.Metadata, error)
	}

	// service struct
	service struct {
		solana solanaClient
	}

	// solana rpc client interface
	solanaClient interface {
		GetSOLBalance(ctx context.Context, base58Addr string) (uint64, error)
		GetTokenBalance(ctx context.Context, base58Addr, base58MintAddr string) (types.TokenAmount, error)
		GetFungibleTokensList(ctx context.Context, walletAddr string) ([]types.TokenAccount, error)
		GetFungibleTokenMetadata(ctx context.Context, base58MintAddr string) (result *metadata.Metadata, err error)
		GetNonFungibleTokensList(ctx context.Context, walletAddr string) ([]types.TokenAccount, error)
		GetTokenMetadata(ctx context.Context, base58MintAddr string) (*token_metadata.Metadata, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(solana solanaClient) Service {
	return &service{solana: solana}
}

// GetSOLBalance returns the SOL balance of a wallet
func (s *service) GetSOLBalance(ctx context.Context, walletAddr string) (Balance, error) {
	balance, err := s.solana.GetSOLBalance(ctx, walletAddr)
	if err != nil {
		return Balance{}, fmt.Errorf("failed to get SOL balance: %w", err)
	}

	metadata, err := s.solana.GetFungibleTokenMetadata(ctx, types.WrappedSOLMint)
	if err != nil {
		return Balance{}, fmt.Errorf("failed to get SOL metadata: %w", err)
	}

	return Balance{
		Pubkey:   walletAddr,
		Mint:     "SOL",
		IsNative: true,
		Balance:  types.NewDefaultTokenAmount(balance),
		Metadata: metadata,
	}, nil
}

// Get balance for specified token
func (s *service) GetTokenBalance(ctx context.Context, walletAddr, tokenMint string) (types.TokenAmount, error) {
	balance, err := s.solana.GetTokenBalance(ctx, walletAddr, tokenMint)
	if err != nil {
		return types.TokenAmount{}, fmt.Errorf("failed to get token balance: %w", err)
	}

	return balance, nil
}

// Get balance for all fungible tokens
func (s *service) GetFungibleTokens(ctx context.Context, walletAddr string) ([]Balance, error) {
	accounts, err := s.solana.GetFungibleTokensList(ctx, walletAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get fungible tokens list: %w", err)
	}

	result := make([]Balance, 0, len(accounts))
	for _, account := range accounts {
		metadata, _ := s.solana.GetFungibleTokenMetadata(ctx, account.Mint.ToBase58()) // ignore error, because it's not critical

		result = append(result, Balance{
			Pubkey:   account.Pubkey.ToBase58(),
			Mint:     account.Mint.ToBase58(),
			IsNative: account.IsNative,
			Balance:  account.Balance,
			Metadata: metadata,
		})
	}

	return result, nil
}

// Get balance for all fungible assets
func (s *service) GetFungibleAssets(ctx context.Context, walletAddr string) ([]Balance, error) {
	accounts, err := s.solana.GetNonFungibleTokensList(ctx, walletAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get fungible assets list: %w", err)
	}

	result := make([]Balance, 0, len(accounts))
	for _, account := range accounts {
		if account.IsNative || account.Mint.ToBase58() == types.WrappedSOLMint || account.Balance.Decimals > 0 {
			continue // skip SOL, wrapped SOL and fungible tokens (decimals > 0)
		}

		metadata, err := s.solana.GetTokenMetadata(ctx, account.Mint.ToBase58())
		if err != nil {
			continue // skip if metadata is not found
		}
		if metadata.Data == nil {
			continue // skip if metadata is not found
		}
		if metadata.TokenStandard != token_metadata.TokenStandardFungibleAsset.String() {
			continue // skip if token standard is not FungibleAsset
		}

		result = append(result, Balance{
			Pubkey:   account.Pubkey.ToBase58(),
			Mint:     account.Mint.ToBase58(),
			IsNative: account.IsNative,
			Balance:  account.Balance,
			Metadata: metadata.Data,
		})
	}

	return result, nil
}

// Get all non-fungible tokens
func (s *service) GetNonFungibleTokens(ctx context.Context, walletAddr string) ([]token_metadata.Metadata, error) {
	accounts, err := s.solana.GetNonFungibleTokensList(ctx, walletAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get non-fungible tokens list: %w", err)
	}

	result := make([]token_metadata.Metadata, 0, len(accounts))
	for _, account := range accounts {
		if account.IsNative || account.Mint.ToBase58() == types.WrappedSOLMint || account.Balance.Decimals > 0 {
			continue // skip SOL, wrapped SOL and fungible tokens (decimals > 0)
		}

		metadata, err := s.solana.GetTokenMetadata(ctx, account.Mint.ToBase58())
		if err != nil {
			continue // skip if metadata is not found
		}
		if metadata.Data == nil {
			continue // skip if metadata is not found
		}
		if metadata.TokenStandard != token_metadata.TokenStandardNonFungible.String() && metadata.TokenStandard != token_metadata.TokenStandardNonFungibleEdition.String() && metadata.TokenStandard != token_metadata.TokenStandardProgrammableNonFungible.String() {
			continue // skip if token standard is not FungibleAsset
		}

		result = append(result, *metadata)
	}

	return result, nil
}
