package wallet

import (
	"context"

	"github.com/dmitrymomot/solana-wallets/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GenerateWallet         endpoint.Endpoint
		StoreWallet            endpoint.Endpoint
		GetWallet              endpoint.Endpoint
		DeleteWallet           endpoint.Endpoint
		UpdateWalletName       endpoint.Endpoint
		ChangeWalletPin        endpoint.Endpoint
		ExportWallet           endpoint.Endpoint
		SignTransaction        endpoint.Endpoint
		SignMessage            endpoint.Endpoint
		SignAndSendTransaction endpoint.Endpoint
	}
)

// Init endpoints
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GenerateWallet:         MakeGenerateWalletEndpoint(s),
		StoreWallet:            MakeStoreWalletEndpoint(s),
		GetWallet:              MakeGetWalletEndpoint(s),
		DeleteWallet:           MakeDeleteWalletEndpoint(s),
		UpdateWalletName:       MakeUpdateWalletNameEndpoint(s),
		ChangeWalletPin:        MakeChangeWalletPinEndpoint(s),
		ExportWallet:           MakeExportWalletEndpoint(s),
		SignTransaction:        MakeSignTransactionEndpoint(s),
		SignMessage:            MakeSignMessageEndpoint(s),
		SignAndSendTransaction: MakeSignAndSendTransactionEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GenerateWallet = mdw(e.GenerateWallet)
			e.StoreWallet = mdw(e.StoreWallet)
			e.GetWallet = mdw(e.GetWallet)
			e.DeleteWallet = mdw(e.DeleteWallet)
			e.UpdateWalletName = mdw(e.UpdateWalletName)
			e.ChangeWalletPin = mdw(e.ChangeWalletPin)
			e.ExportWallet = mdw(e.ExportWallet)
			e.SignTransaction = mdw(e.SignTransaction)
			e.SignMessage = mdw(e.SignMessage)
			e.SignAndSendTransaction = mdw(e.SignAndSendTransaction)
		}
	}

	return e
}

// MakeGenerateWalletEndpoint returns an endpoint function for the GenerateWallet method.
func MakeGenerateWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GenerateWallet(ctx)
	}
}

// StoreWalletRequest is a request for StoreWallet method
type StoreWalletRequest struct {
	UserID   string `json:"-" validate:"required" label:"User ID"`
	Name     string `json:"name" validate:"required|minLen:3|maxLen:50" label:"Name"`
	Pin      string `json:"pin" validate:"required|minLen:4|maxLen:50" label:"PIN Code"`
	Mnemonic string `json:"mnemonic" validate:"required" label:"Mnemonic"`
}

// MakeStoreWalletEndpoint returns an endpoint function for the StoreWallet method.
func MakeStoreWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(StoreWalletRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		if err := s.StoreWallet(ctx, req.UserID, req.Pin, req.Mnemonic, req.Name); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetWalletEndpoint returns an endpoint function for the GetWallet method.
func MakeGetWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userID, ok := request.(string)
		if !ok {
			return nil, ErrInvalidParameter
		}

		return s.GetWallet(ctx, userID)
	}
}

// DeleteWalletRequest is a request for DeleteWallet method
type DeleteWalletRequest struct {
	UserID string `json:"-" validate:"required" label:"User ID"`
	Pin    string `json:"pin" validate:"required" label:"PIN Code"`
}

// MakeDeleteWalletEndpoint returns an endpoint function for the DeleteWallet method.
func MakeDeleteWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteWalletRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		if err := s.DeleteWallet(ctx, req.UserID, req.Pin); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// UpdateWalletNameRequest is a request for UpdateWalletName method
type UpdateWalletNameRequest struct {
	UserID string `json:"-" validate:"required" label:"User ID"`
	Pin    string `json:"pin" validate:"required" label:"PIN Code"`
	Name   string `json:"name" validate:"required|minLen:3|maxLen:50" label:"Name"`
}

// MakeUpdateWalletNameEndpoint returns an endpoint function for the UpdateWalletName method.
func MakeUpdateWalletNameEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateWalletNameRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		if err := s.UpdateWalletName(ctx, req.UserID, req.Pin, req.Name); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// ChangeWalletPinRequest is a request for ChangeWalletPin method
type ChangeWalletPinRequest struct {
	UserID string `json:"-" validate:"required" label:"User ID"`
	Pin    string `json:"pin" validate:"required" label:"PIN Code"`
	NewPin string `json:"new_pin" validate:"required|minLen:4|maxLen:50" label:"New PIN Code"`
}

// MakeChangeWalletPinEndpoint returns an endpoint function for the ChangeWalletPin method.
func MakeChangeWalletPinEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ChangeWalletPinRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		if err := s.ChangeWalletPin(ctx, req.UserID, req.Pin, req.NewPin); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// ExportWalletRequest is a request for ExportWallet method
type ExportWalletRequest struct {
	UserID string `json:"-" validate:"required" label:"User ID"`
	Pin    string `json:"pin" validate:"required" label:"PIN Code"`
}

// MakeExportWalletEndpoint returns an endpoint function for the ExportWallet method.
func MakeExportWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ExportWalletRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		return s.ExportWallet(ctx, req.UserID, req.Pin)
	}
}

// SignTransactionRequest is a request for SignTransaction method
type SignTransactionRequest struct {
	UserID string `json:"-" validate:"required" label:"User ID"`
	Pin    string `json:"pin" validate:"required" label:"PIN Code"`
	Tx     string `json:"tx" validate:"required" label:"Base64 encoded transaction"`
}

// MakeSignTransactionEndpoint returns an endpoint function for the SignTransaction method.
func MakeSignTransactionEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SignTransactionRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		return s.SignTransaction(ctx, req.UserID, req.Pin, req.Tx)
	}
}

type (
	// SignMessageRequest is a request for SignMessage method
	SignMessageRequest struct {
		UserID string `json:"-" validate:"required" label:"User ID"`
		Pin    string `json:"pin" validate:"required" label:"PIN Code"`
		Msg    string `json:"msg" validate:"required" label:"Message"`
	}

	// SignMessageResponse is a response for SignMessage method
	SignMessageResponse struct {
		Signature string `json:"signature" label:"Signature"`
		Msg       string `json:"msg" label:"Message"`
	}
)

// MakeSignMessageEndpoint returns an endpoint function for the SignMessage method.
func MakeSignMessageEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SignMessageRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		msg, sig, err := s.SignMessage(ctx, req.UserID, req.Pin, req.Msg)
		if err != nil {
			return nil, err
		}

		return SignMessageResponse{
			Signature: sig,
			Msg:       msg,
		}, nil
	}
}

type (
	// SignAndSendTransactionRequest is a request for SignAndSendTransaction method
	SignAndSendTransactionRequest struct {
		UserID string `json:"-" validate:"required" label:"User ID"`
		Pin    string `json:"pin" validate:"required" label:"PIN Code"`
		Tx     string `json:"tx" validate:"required" label:"Base64 encoded transaction"`
	}

	// SignAndSendTransactionResponse is a response for SignAndSendTransaction method
	SignAndSendTransactionResponse struct {
		TxSignature string `json:"tx_signature" label:"Transaction signature"`
	}
)

// MakeSignAndSendTransactionEndpoint returns an endpoint function for the SignAndSendTransaction method.
func MakeSignAndSendTransactionEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SignAndSendTransactionRequest)
		if !ok {
			return nil, ErrInvalidParameter
		}
		if v := validator.ValidateStruct(req); v != nil {
			return nil, validator.NewValidationError(v)
		}

		sig, err := s.SignAndSendTransaction(ctx, req.UserID, req.Pin, req.Tx)
		if err != nil {
			return nil, err
		}

		return SignAndSendTransactionResponse{
			TxSignature: sig,
		}, nil
	}
}
