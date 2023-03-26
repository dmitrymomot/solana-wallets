package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestService_Example(t *testing.T) {
	var (
		uid = uuid.New()
	)
	type args struct {
		ctx context.Context
		uid uuid.UUID
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "example",
			s:    NewService(),
			args: args{
				ctx: context.Background(),
				uid: uid,
			},
			want:    fmt.Sprintf("example: user id: %s", uid.String()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			got, err := s.Example(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Example() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.Example() = %v, want %v", got, tt.want)
			}
		})
	}
}
