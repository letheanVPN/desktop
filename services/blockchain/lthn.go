package blockchain

import (
	"context"
	"encoding/json"
	"fmt"

	lthn "github.com/letheanVPN/blockchain/utils/sdk/client/go"
	"github.com/letheanVPN/desktop/services/core/config"
)

// lthnNetwork implements the Network interface for the Lethean blockchain.
type lthnNetwork struct {
	config    *config.Config
	apiClient *lthn.APIClient
}

// newLthnNetwork creates a new instance of the Lethean network handler.
func newLthnNetwork(cfg *config.Config) *lthnNetwork {
	return &lthnNetwork{
		config: cfg,
	}
}

// Connect sets up the connection to the Lethean node.
func (n *lthnNetwork) Connect() error {
	// TODO: The endpoint should be loaded from a workspace or app config.
	configuration := lthn.NewConfiguration()
	// For example: configuration.Host = n.config.GetLetheanNodeHost()
	n.apiClient = lthn.NewAPIClient(configuration)
	fmt.Println("Connected to Lethean network implementation.")
	return nil
}

// Disconnect tears down the connection.
func (n *lthnNetwork) Disconnect() error {
	n.apiClient = nil
	fmt.Println("Disconnected from Lethean network implementation.")
	return nil
}

// FetchBlockData retrieves block data from the Lethean network.
func (n *lthnNetwork) FetchBlockData(identifier string) (string, error) {
	if n.apiClient == nil {
		return "", fmt.Errorf("not connected to Lethean network")
	}

	resp, _, err := n.apiClient.BlockUtilsSdkClientGo.GetBlock(context.Background(), identifier).Execute()
	if err != nil {
		return "", fmt.Errorf("error when calling `GetBlock`: %w", err)
	}

	// Marshal the response to a JSON string to return to the frontend.
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal block details to JSON: %w", err)
	}

	return string(jsonResponse), nil
}
