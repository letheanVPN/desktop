package display

import "github.com/wailsapp/wails/v3/pkg/application"

// OpenWindow creates and shows a new webview window.
// This function is callable from the frontend.
func (s *Service) OpenWindow(name string, options application.WebviewWindowOptions) {
	// Check if a window with that name already exists
	if _, exists := s.windowHandles[name]; exists {
		// You might want to focus the existing window instead of creating a new one
		s.windowHandles[name].Focus()
		return
	}

	window := s.app.Window.NewWithOptions(options)
	s.windowHandles[name] = window
	window.Show()
}

// SelectDirectory opens a directory selection dialog and returns the selected path.
func (s *Service) SelectDirectory() (string, error) {
	dialog := application.OpenFileDialog()
	dialog.SetTitle("Select Project Directory")
	if path, err := dialog.PromptForSingleSelection(); err == nil {
		// Use selected directory path
		return path, nil
	}
	return "", nil
}
