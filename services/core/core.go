package core

import (
	"embed"
	"fmt"
	"sync"

	"github.com/letheanVPN/desktop/services/core/config"
	"github.com/letheanVPN/desktop/services/core/display"
	"github.com/letheanVPN/desktop/services/core/i18n"
	"github.com/letheanVPN/desktop/services/docs"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Service provides access to all core application services.
type Service struct {
	app            *application.App
	configService  *config.Service
	i18nService    *i18n.Service
	displayService *display.Service
	docsService    *docs.Service
}

var (
	instance *Service
	once     sync.Once
	initErr  error
)

// New performs Phase 1 of initialization: Instantiation.
// It creates the raw service objects without wiring them together.
func New(assets embed.FS) *Service {
	once.Do(func() {
		// Instantiate services in the correct order of dependency.
		configService, err := config.NewService()
		if err != nil {
			initErr = fmt.Errorf("failed to initialize config service: %w", err)
			return
		}

		i18nService, err := i18n.NewService(configService.Get().Language)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize i18n service: %w", err)
			return
		}

		displayService := display.NewService(display.ClientHub, assets)
		docsService := docs.NewService(assets)

		instance = &Service{
			configService:  configService,
			i18nService:    i18nService,
			displayService: displayService,
			docsService:    docsService,
		}
	})

	if initErr != nil {
		panic(initErr) // A failure in a core service is fatal.
	}

	return instance
}

// Setup performs Phase 2 of initialization: Wiring.
// It injects the required dependencies into each service.
func Setup(app *application.App) {
	if instance == nil {
		panic("core.Setup() called before core.New() was successfully initialized")
	}
	instance.app = app

	// Wire the services with their dependencies.
	instance.displayService.Setup(app, instance.configService, instance.i18nService)
	instance.docsService.Setup(app, instance.displayService)
}

// App returns the global application instance.
func App() *application.App {
	if instance == nil || instance.app == nil {
		panic("core.App() called before core.Setup() was successfully initialized")
	}
	return instance.app
}

// Config returns the singleton instance of the ConfigService.
func Config() *config.Service {
	if instance == nil {
		panic("core.Config() called before core.New() was successfully initialized")
	}
	return instance.configService
}

// I18n returns the singleton instance of the i18n.Service.
func I18n() *i18n.Service {
	if instance == nil {
		panic("core.I18n() called before core.New() was successfully initialized")
	}
	return instance.i18nService
}

// Display returns the singleton instance of the display.Service.
func Display() *display.Service {
	if instance == nil {
		panic("core.Display() called before core.New() was successfully initialized")
	}
	return instance.displayService
}

// Docs returns the singleton instance of the DocsService.
func Docs() *docs.Service {
	if instance == nil {
		panic("core.Docs() called before core.New() was successfully initialized")
	}
	return instance.docsService
}
