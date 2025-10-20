package blockchain

// Network defines the standard interface for a blockchain network implementation.
// This allows different chains (e.g., Lethean, an alias sidechain) to be used interchangeably.
type Network interface {
	// Connect establishes a connection to the network's node or endpoint.
	Connect() error

	// Disconnect tears down the connection.
	Disconnect() error

	// FetchBlockData retrieves data for a specific block, given a hash or height.
	FetchBlockData(identifier string) (string, error)
}
