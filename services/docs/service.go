package docs

import (
	"net/http"
	"strings"

	"github.com/letheanVPN/desktop/services/display"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func (s *Service) Get() *Config {
	return s.config
}

// NewService creates and initializes a new documentation service.
func NewService(displaySvc *display.Service) (*Service, error) {
	return &Service{
		config: &Config{
			"",
		},
		display: displaySvc,
	}, nil
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
	s.display.OpenWindow("docs", application.WebviewWindowOptions{
		Title:       "Lethean Documentation",
		Height:      600,
		Width:       1000,
		URL:         url,
		AlwaysOnTop: true,
		Frameless:   false,
	})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(docsStatic).ServeHTTP(w, r)
	//application.AssetFileServerFS(docsStatic).ServeHTTP(w, r)
}
