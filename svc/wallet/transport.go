package wallet

import (
	"context"
	"encoding/json"
	"fmt"
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

	r.Get("/generate", httptransport.NewServer(
		e.GenerateWallet,
		decodeEmptyRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/store", httptransport.NewServer(
		e.StoreWallet,
		decodeStoreWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/public/{id}", httptransport.NewServer(
		e.GetWallet,
		decodeGetWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/delete", httptransport.NewServer(
		e.DeleteWallet,
		decodeDeleteWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Patch("/update/name", httptransport.NewServer(
		e.UpdateWalletName,
		decodeUpdateWalletNameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Patch("/update/pin", httptransport.NewServer(
		e.ChangeWalletPin,
		decodeChangeWalletPinRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/export", httptransport.NewServer(
		e.ExportWallet,
		decodeExportWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/transaction/sign", httptransport.NewServer(
		e.SignTransaction,
		decodeSignTransactionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/message/sign", httptransport.NewServer(
		e.SignMessage,
		decodeSignMessageRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/transaction/sign/send", httptransport.NewServer(
		e.SignAndSendTransaction,
		decodeSignAndSendTransactionRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) || errors.Is(err, ErrInvalidPIN) {
		return http.StatusBadRequest, err.Error()
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden, err.Error()
	}
	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeEmptyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeStoreWalletRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StoreWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeGetWalletRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrInvalidParameter
	}

	return id, nil
}

func decodeDeleteWalletRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req DeleteWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeUpdateWalletNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateWalletNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeChangeWalletPinRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ChangeWalletPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeExportWalletRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ExportWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeSignTransactionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeSignMessageRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SignMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}

func decodeSignAndSendTransactionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SignAndSendTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	return req, nil
}
