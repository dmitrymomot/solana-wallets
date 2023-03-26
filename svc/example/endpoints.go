package example

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		Example endpoint.Endpoint
	}

	service interface {
		Example(ctx context.Context, uid uuid.UUID) (string, error)
	}
)

// Init endpoints
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		Example: MakeExampleEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.Example = mdw(e.Example)
		}
	}

	return e
}

// MakeExampleEndpoint ...
func MakeExampleEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		uid, ok := req.(uuid.UUID)
		if !ok {
			return nil, ErrInvalidParameter
		}

		resp, err := s.Example(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
