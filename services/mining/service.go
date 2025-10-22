package mining

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// New creates a new mining service
func New() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}
