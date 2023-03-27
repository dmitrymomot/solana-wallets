package wallet

import (
	"context"
	"net/http"

	"github.com/dmitrymomot/solana-wallets/internal/httpencoder"
	"github.com/go-chi/chi/v5"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/{wallet}/tokens", httptransport.NewServer(
		e.GetBalance,
		decodeGetBalanceRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet}/tokens/{mint}", httptransport.NewServer(
		e.GetTokenBalance,
		decodeGetTokenBalanceRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet}/assets", httptransport.NewServer(
		e.GetAssets,
		decodeGetAssetsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet}/nfts", httptransport.NewServer(
		e.GetNFTs,
		decodeGetNFTsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

// decodeGetBalanceRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetBalanceRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	wallet := chi.URLParam(req, "wallet")
	if wallet == "" {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid wallet")
	}
	return wallet, nil
}

// decodeGetTokenBalanceRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetTokenBalanceRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	wallet := chi.URLParam(req, "wallet")
	if wallet == "" {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid wallet")
	}

	mint := chi.URLParam(req, "mint")
	if mint == "" {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid mint")
	}

	return GetTokenBalanceRequest{
		WalletAddr: wallet,
		TokenMint:  mint,
	}, nil
}

// decodeGetAssetsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetAssetsRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	wallet := chi.URLParam(req, "wallet")
	if wallet == "" {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid wallet")
	}

	return wallet, nil
}

// decodeGetNFTsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetNFTsRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	wallet := chi.URLParam(req, "wallet")
	if wallet == "" {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid wallet")
	}

	return wallet, nil
}
