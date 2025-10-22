package core

import (
	"fmt"
	"sync"

	"github.com/letheanVPN/desktop/services/config"
	"github.com/letheanVPN/desktop/services/core/i18n"
)

// Service provides access to all core application services.
type Service struct {
	configService *config.Service
	i18nService   *i18n.Service // Changed to concrete type
}

var (
	instance *Service
	once     sync.Once
	initErr  error
)

// New initializes the core service singleton.
func New() *Service {
	once.Do(func() {
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

		instance = &Service{
			configService: configService,
			i18nService:   i18nService,
		}
	})

	return instance
}

// Config returns the singleton instance of the ConfigService.
func Config() *config.Service {
	if instance == nil {
		panic("core.Config() called before core.NewService() was successfully initialized")
	}
	return instance.configService
}

// I18n returns the singleton instance of the i18n.Service.
func I18n() *i18n.Service { // Changed to return concrete type
	if instance == nil {
		panic("core.I18n() called before core.NewService() was successfully initialized")
	}
	return instance.i18nService
}
