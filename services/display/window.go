package display

import "github.com/wailsapp/wails/v3/pkg/application"

// OpenWindow creates and shows a new webview window.
// This function is callable from the frontend.
func (s *Service) OpenWindow(app *application.App, name string, options application.WebviewWindowOptions) {
	// Check if a window with that name already exists
	if _, exists := s.windowHandles[name]; exists {
		// You might want to focus the existing window instead of creating a new one
		s.windowHandles[name].Focus()
		return
	}

	window := app.Window.NewWithOptions(options)
	s.windowHandles[name] = window
	window.Show()
}
