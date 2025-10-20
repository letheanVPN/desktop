package blockchain

import (
	"fmt"

	"github.com/letheanVPN/desktop/services/config"
)

// Service manages different blockchain network implementations.
type Service struct {
	config        *config.Config
	networks      map[string]Network
	activeNetwork Network
}

// NewService creates a new, uninitialized blockchain service.
func NewService(cfg *config.Config) *Service {
	s := &Service{
		config:   cfg,
		networks: make(map[string]Network),
	}

	// Register all available network implementations here.
	s.networks["lthn"] = newLthnNetwork(cfg)

	return s
}

// ServiceStartup Startup is called by Wails. We can use it to connect to a default network.
func (s *Service) ServiceStartup() error {
	return s.Start("lthn") // Start the default Lethean network.
}

// Start activates a specific blockchain network.
func (s *Service) Start(networkID string) error {
	if s.activeNetwork != nil {
		s.activeNetwork.Disconnect()
	}

	network, exists := s.networks[networkID]
	if !exists {
		return fmt.Errorf("network '%s' is not supported", networkID)
	}

	if err := network.Connect(); err != nil {
		return fmt.Errorf("failed to connect to network '%s': %w", networkID, err)
	}

	s.activeNetwork = network
	fmt.Printf("Blockchain service started for network: %s\n", networkID)
	return nil
}

// FetchBlockData routes the call to the active network.
func (s *Service) FetchBlockData(identifier string) (string, error) {
	if s.activeNetwork == nil {
		return "", fmt.Errorf("no active blockchain network")
	}
	return s.activeNetwork.FetchBlockData(identifier)
}

// Install can be used to download and set up blockchain binaries.
func (s *Service) Install(networkID string) error {
	// TODO: Implement logic to download and verify blockchain binaries
	// for the given network, using the paths from s.config.
	fmt.Printf("Install requested for network: %s\n", networkID)
	return nil
}
