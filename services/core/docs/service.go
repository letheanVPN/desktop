package docs

import (
	"net/http"
	"strings"

	"github.com/letheanVPN/desktop/services/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// NewService creates a new, un-wired documentation service.
func NewService() *Service {
	return &Service{
		config: &Config{},
	}
}

// Setup injects the required dependencies into the service.
func (s *Service) Setup(app *application.App, displayService core.DisplayService) {
	s.app = app
	s.displayService = displayService
}

// OpenDocsWindow opens a new window with the documentation.
func (s *Service) OpenDocsWindow(path ...string) {
	url := "/docs/"
	if len(path) > 0 {
		fullPath := path[0]
		if strings.Contains(fullPath, "#") {
			parts := strings.SplitN(fullPath, "#", 2)
			pagePath := parts[0]
			fragment := parts[1]
			url += pagePath + "/#" + fragment
		} else {
			url += fullPath
		}
	}

	// Use the injected displayService, not the global core.Display().
	s.displayService.OpenWindow("docs", application.WebviewWindowOptions{
		Title:       "Lethean Documentation",
		Height:      600,
		Width:       1000,
		URL:         url,
		AlwaysOnTop: true,
		Frameless:   false,
	})
}

// ServeHTTP serves the embedded documentation assets.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(docsStatic).ServeHTTP(w, r)
}
