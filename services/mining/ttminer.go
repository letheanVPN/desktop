package mining

import (
	"errors"
)

// TTMiner represents a TT-Miner
type TTMiner struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
	Path    string `json:"path"`
	Running bool   `json:"running"`
	Pid     int    `json:"pid"`
}

// NewTTMiner creates a new TT-Miner
func NewTTMiner() *TTMiner {
	return &TTMiner{
		Name:    "TT-Miner",
		Version: "latest",
		URL:     "https://github.com/TrailingStop/TT-Miner-release",
	}
}

// GetName returns the name of the miner
func (m *TTMiner) GetName() string {
	return m.Name
}

// Install the miner
func (m *TTMiner) Install() error {
	return errors.New("not implemented")
}

// Start the miner
func (m *TTMiner) Start(config *Config) error {
	return errors.New("not implemented")
}

// Stop the miner
func (m *TTMiner) Stop() error {
	return errors.New("not implemented")
}

// GetStats returns the stats for the miner
func (m *TTMiner) GetStats() (*PerformanceMetrics, error) {
	return nil, errors.New("not implemented")
}
