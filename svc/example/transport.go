package example

import (
	"context"
	"net/http"

	"github.com/dmitrymomot/go-api-server/internal/httpencoder"
	"github.com/go-chi/chi/v5"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
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

	r.Get("/{uid}", httptransport.NewServer(
		e.Example,
		decodeExampleRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeExampleRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	uid, err := uuid.Parse(chi.URLParam(req, "uid"))
	if err != nil {
		return nil, errors.Wrap(ErrInvalidParameter, "invalid user id")
	}
	return uid, nil
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
