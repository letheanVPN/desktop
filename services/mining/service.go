package mining

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// New creates a new mining service
func New() *Service {
	return &Service{
		//config:        &config.Config{},
		//networks:      make(map[string]Network),
		//activeNetwork: nil,
	}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.XMRig = NewXMRigMiner()
	return nil
}

// InstallXMRig installs the latest version of XMRig
func (s *Service) InstallXMRig() error {
	return s.XMRig.Install()
}

// StartXMRig starts the XMRig miner
func (s *Service) StartXMRig(config *Config) error {
	return s.XMRig.Start(config)
}

// StopXMRig stops the XMRig miner
func (s *Service) StopXMRig() error {
	return s.XMRig.Stop()
}

// GetXMRigStats returns the stats for the XMRig miner
func (s *Service) GetXMRigStats() (*Stats, error) {
	return s.XMRig.GetStats()
}
