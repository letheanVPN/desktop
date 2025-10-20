package local

// Medium implements the filesystem.Medium interface for the local disk.
type Medium struct {
	root string
}
