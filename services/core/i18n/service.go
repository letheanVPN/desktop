package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

// GetAvailableLanguages scans the embedded locales directory and returns a slice of supported language.Tag.
func GetAvailableLanguages() ([]language.Tag, error) {
	files, err := localeFS.ReadDir("locales")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded locales directory: %w", err)
	}

	var availableLangs []language.Tag
	for _, file := range files {
		lang := strings.TrimSuffix(file.Name(), ".json")
		tag := language.Make(lang)
		availableLangs = append(availableLangs, tag)
	}
	return availableLangs, nil
}

// DetectLanguage attempts to auto-detect the user's preferred language.
func DetectLanguage() (string, error) {
	langEnv := os.Getenv("LANG")
	if langEnv == "" {
		return "", nil // No LANG env var, not an error, just no detection
	}

	baseLang := strings.Split(langEnv, ".")[0]
	parsedLang, err := language.Parse(baseLang)
	if err != nil {
		return "", fmt.Errorf("failed to parse language tag '%s': %w", baseLang, err)
	}

	supported, err := GetAvailableLanguages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not get available languages for auto-detection: %v\n", err)
		return "", nil
	}
	if len(supported) == 0 {
		return "", nil
	}

	matcher := language.NewMatcher(supported)
	_, index, confidence := matcher.Match(parsedLang)

	if confidence >= language.Low {
		return supported[index].String(), nil
	}
	return "", nil
}

// NewService creates and initializes a new i18n service.
func NewService(defaultLang string) (*Service, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	availableLangs, err := GetAvailableLanguages()
	if err != nil {
		return nil, err
	}

	for _, lang := range availableLangs {
		filePath := fmt.Sprintf("locales/%s.json", lang.String())
		_, err := bundle.LoadMessageFileFS(localeFS, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load message file %s: %w", filePath, err)
		}
	}

	s := &Service{
		bundle:         bundle,
		availableLangs: availableLangs,
	}
	s.SetLanguage(defaultLang)

	return s, nil
}

// SetLanguage changes the active language for translations.
func (s *Service) SetLanguage(lang string) {
	s.localizer = i18n.NewLocalizer(s.bundle, lang)
}

// Translate a message by its ID using the configured language.
func (s *Service) Translate(messageID string) string {
	translation, err := s.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "i18n: translation for key \"%s\" not found\n", messageID)
		return messageID
	}
	return translation
}

// TranslateWithConfig translates a message using a LocalizeConfig.
func (s *Service) TranslateWithConfig(lc *i18n.LocalizeConfig) string {
	translation, err := s.localizer.Localize(lc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "i18n: translation for key \"%s\" not found\n", lc.MessageID)
		return lc.MessageID
	}
	return translation
}

// Bundle returns the underlying i18n.Bundle.
func (s *Service) Bundle() *i18n.Bundle {
	return s.bundle
}

// GetAllMessages returns all messages for a given language as a map.
// This is intended for frontend pre-loading to enable synchronous lookups.
func (s *Service) GetAllMessages(lang string) (map[string]string, error) {
	filePath := fmt.Sprintf("locales/%s.json", lang)
	data, err := localeFS.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locale file %s: %w", filePath, err)
	}

	messages := make(map[string]string)
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal locale file %s: %w", filePath, err)
	}

	return messages, nil
}

// AvailableLanguages returns the list of available languages.
func (s *Service) AvailableLanguages() []language.Tag {
	return s.availableLangs
}
