package example_worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type (
	// Worker is a task handler for email delivery.
	Worker struct {
		mail mailer
	}

	mailer interface {
		SendExample(ctx context.Context, exampleID, email, role, merchantName string) error
	}
)

// NewWorker creates a new email task handler.
func NewWorker(mail mailer) *Worker {
	return &Worker{mail: mail}
}

// Schedule schedules tasks for the worker.
func (w *Worker) Schedule(s *asynq.Scheduler) {
	s.Register("@every 1h", asynq.NewTask(SendExampleTask, nil))
}

// Register registers task handlers for email delivery.
func (w *Worker) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(SendExampleTask, w.SendExampleEmail)
}

// SendExampleEmail sends an example email.
func (w *Worker) SendExampleEmail(ctx context.Context, t *asynq.Task) error {
	var p SendExamplePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := w.mail.SendExample(ctx, p.ExampleID, p.Email, p.Role, p.MerchantName); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
