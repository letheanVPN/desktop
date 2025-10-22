package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Service provides internationalization and localization.
// It is the concrete implementation for the i18n service.
type Service struct {
	bundle         *i18n.Bundle
	localizer      *i18n.Localizer
	availableLangs []language.Tag
}
