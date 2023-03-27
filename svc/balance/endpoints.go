package wallet

import (
	"context"

	"github.com/dmitrymomot/solana/token_metadata"
	"github.com/dmitrymomot/solana/types"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetBalance      endpoint.Endpoint
		GetAssets       endpoint.Endpoint
		GetNFTs         endpoint.Endpoint
		GetTokenBalance endpoint.Endpoint
	}

	BalanceResponse struct {
		Balances []Balance `json:"balances"`
	}
)

// Init endpoints
func MakeEndpoints(s Service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetBalance:      MakeGetBalanceEndpoint(s),
		GetAssets:       MakeGetAssetsEndpoint(s),
		GetNFTs:         MakeGetNFTsEndpoint(s),
		GetTokenBalance: MakeGetTokenBalanceEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetBalance = mdw(e.GetBalance)
			e.GetAssets = mdw(e.GetAssets)
			e.GetNFTs = mdw(e.GetNFTs)
			e.GetTokenBalance = mdw(e.GetTokenBalance)
		}
	}

	return e
}

// MakeGetBalanceEndpoint returns an endpoint function for the GetBalance method.
func MakeGetBalanceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		walletAddr, ok := req.(string)
		if !ok {
			return nil, ErrInvalidParameter
		}

		result := []Balance{}

		if solBalance, err := s.GetSOLBalance(ctx, walletAddr); err == nil {
			result = append(result, solBalance)
		}

		if splBalace, err := s.GetFungibleTokens(ctx, walletAddr); err == nil {
			result = append(result, splBalace...)
		}

		return BalanceResponse{Balances: result}, nil
	}
}

// MakeGetAssetsEndpoint returns an endpoint function for the GetAssets method.
func MakeGetAssetsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		walletAddr, ok := req.(string)
		if !ok {
			return nil, ErrInvalidParameter
		}

		result, err := s.GetFungibleAssets(ctx, walletAddr)
		if err != nil {
			return nil, err
		}

		return BalanceResponse{Balances: result}, nil
	}
}

// NFTsResponse is a response for the GetNFTs method.
type NFTsResponse struct {
	NFTs []token_metadata.Metadata `json:"nfts"`
}

// MakeGetNFTsEndpoint returns an endpoint function for the GetNFTs method.
func MakeGetNFTsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		walletAddr, ok := req.(string)
		if !ok {
			return nil, ErrInvalidParameter
		}

		result, err := s.GetNonFungibleTokens(ctx, walletAddr)
		if err != nil {
			return nil, err
		}

		return NFTsResponse{NFTs: result}, nil
	}
}

type (
	// GetTokenBalanceRequest is a request payload for the GetTokenBalance method.
	GetTokenBalanceRequest struct {
		WalletAddr string `json:"wallet_addr"`
		TokenMint  string `json:"mint"`
	}

	TokenBalanceResponse struct {
		Balance types.TokenAmount `json:"balance"`
	}
)

// MakeGetTokenBalanceEndpoint returns an endpoint function for the GetTokenBalance method.
func MakeGetTokenBalanceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetTokenBalanceRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if req.WalletAddr == "" || req.TokenMint == "" {
			return nil, ErrInvalidParameter
		}

		result, err := s.GetTokenBalance(ctx, req.WalletAddr, req.TokenMint)
		if err != nil {
			return nil, err
		}

		return TokenBalanceResponse{Balance: result}, nil
	}
}
