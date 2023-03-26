package example

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService() *Service {
	return &Service{}
}

// Example ...
func (s *Service) Example(ctx context.Context, uid uuid.UUID) (string, error) {
	return fmt.Sprintf("example: user id: %s", uid.String()), nil
}
