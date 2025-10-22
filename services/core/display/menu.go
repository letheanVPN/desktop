package display

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// buildMenu creates and sets the main application menu.
func (s *Service) buildMenu() {
	appMenu := s.app.Menu.New()
	if runtime.GOOS == "darwin" {
		appMenu.AddRole(application.AppMenu)
	}
	appMenu.AddRole(application.FileMenu)
	appMenu.AddRole(application.ViewMenu)
	appMenu.AddRole(application.EditMenu)

	workspace := appMenu.AddSubmenu("Workspace")
	workspace.Add("New").OnClick(func(ctx *application.Context) { /* TODO */ })
	workspace.Add("List").OnClick(func(ctx *application.Context) { /* TODO */ })

	// Add brand-specific menu items
	if s.brand == DeveloperHub {
		appMenu.AddSubmenu("Developer")
	}

	appMenu.AddRole(application.WindowMenu)
	appMenu.AddRole(application.HelpMenu)

	s.app.Menu.Set(appMenu)
}
