// Package language provides dynamic language support for HaiLanGo.
// Instead of hardcoding supported languages, this package allows ANY language
// that the underlying LLM/TTS/STT APIs support.
package language

import (
	"strings"
	"sync"
)

// SupportTier indicates the level of support for a language
type SupportTier string

const (
	// TierVerified - Language has been tested and verified for quality
	TierVerified SupportTier = "verified"
	// TierSupported - Language is known to be supported by APIs but not fully tested
	TierSupported SupportTier = "supported"
	// TierExperimental - Language may work but quality is not guaranteed
	TierExperimental SupportTier = "experimental"
)

// Info contains metadata about a language
type Info struct {
	Code                  string      `json:"code"`                    // ISO 639-1 code (e.g., "ja", "en", "ku")
	Name                  string      `json:"name"`                    // English name
	NativeName            string      `json:"native_name"`             // Name in the language itself
	SupportTier           SupportTier `json:"support_tier"`            // Level of support
	TTSVoices             []string    `json:"tts_voices,omitempty"`    // Known TTS voices
	SupportsPronunciation bool        `json:"supports_pronunciation"`  // STT pronunciation evaluation support
	Notes                 string      `json:"notes,omitempty"`         // Additional notes
}

// Registry manages language information dynamically
type Registry struct {
	mu        sync.RWMutex
	languages map[string]*Info // code -> Info
}

// Global registry instance
var (
	globalRegistry *Registry
	once           sync.Once
)

// GetRegistry returns the global language registry
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = NewRegistry()
		globalRegistry.initWellKnownLanguages()
	})
	return globalRegistry
}

// NewRegistry creates a new language registry
func NewRegistry() *Registry {
	return &Registry{
		languages: make(map[string]*Info),
	}
}

// initWellKnownLanguages initializes languages that have been verified
func (r *Registry) initWellKnownLanguages() {
	// Verified languages - tested and confirmed to work well
	verified := []*Info{
		{Code: "ja", Name: "Japanese", NativeName: "日本語", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"ja-JP-Neural2-B", "ja-JP-Neural2-C"}},
		{Code: "en", Name: "English", NativeName: "English", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"en-US-Neural2-A", "en-US-Neural2-C"}},
		{Code: "zh", Name: "Chinese", NativeName: "中文", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"zh-CN-Neural2-A", "zh-CN-Neural2-B"}},
		{Code: "ru", Name: "Russian", NativeName: "Русский", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"ru-RU-Wavenet-A", "ru-RU-Wavenet-B"}},
		{Code: "es", Name: "Spanish", NativeName: "Español", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"es-ES-Neural2-A", "es-ES-Neural2-B"}},
		{Code: "fr", Name: "French", NativeName: "Français", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"fr-FR-Neural2-A", "fr-FR-Neural2-B"}},
		{Code: "de", Name: "German", NativeName: "Deutsch", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"de-DE-Neural2-A", "de-DE-Neural2-B"}},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"pt-BR-Neural2-A", "pt-BR-Neural2-B"}},
		{Code: "it", Name: "Italian", NativeName: "Italiano", SupportTier: TierVerified, SupportsPronunciation: true, TTSVoices: []string{"it-IT-Neural2-A", "it-IT-Neural2-B"}},
	}

	// Supported languages - known to work but less tested
	supported := []*Info{
		{Code: "fa", Name: "Persian", NativeName: "فارسی", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"fa-IR-Wavenet-A"}},
		{Code: "he", Name: "Hebrew", NativeName: "עברית", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"he-IL-Wavenet-A"}},
		{Code: "tr", Name: "Turkish", NativeName: "Türkçe", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"tr-TR-Wavenet-A"}},
		{Code: "ko", Name: "Korean", NativeName: "한국어", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"ko-KR-Neural2-A"}},
		{Code: "ar", Name: "Arabic", NativeName: "العربية", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"ar-XA-Wavenet-A"}},
		{Code: "hi", Name: "Hindi", NativeName: "हिन्दी", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"hi-IN-Neural2-A"}},
		{Code: "th", Name: "Thai", NativeName: "ไทย", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"th-TH-Neural2-C"}},
		{Code: "vi", Name: "Vietnamese", NativeName: "Tiếng Việt", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"vi-VN-Neural2-A"}},
		{Code: "nl", Name: "Dutch", NativeName: "Nederlands", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"nl-NL-Neural2-A"}},
		{Code: "pl", Name: "Polish", NativeName: "Polski", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"pl-PL-Neural2-A"}},
		{Code: "uk", Name: "Ukrainian", NativeName: "Українська", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"uk-UA-Wavenet-A"}},
		{Code: "cs", Name: "Czech", NativeName: "Čeština", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"cs-CZ-Wavenet-A"}},
		{Code: "sv", Name: "Swedish", NativeName: "Svenska", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"sv-SE-Neural2-A"}},
		{Code: "da", Name: "Danish", NativeName: "Dansk", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"da-DK-Neural2-D"}},
		{Code: "fi", Name: "Finnish", NativeName: "Suomi", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"fi-FI-Neural2-A"}},
		{Code: "no", Name: "Norwegian", NativeName: "Norsk", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"nb-NO-Neural2-A"}},
		{Code: "el", Name: "Greek", NativeName: "Ελληνικά", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"el-GR-Neural2-A"}},
		{Code: "hu", Name: "Hungarian", NativeName: "Magyar", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"hu-HU-Wavenet-A"}},
		{Code: "ro", Name: "Romanian", NativeName: "Română", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"ro-RO-Wavenet-A"}},
		{Code: "id", Name: "Indonesian", NativeName: "Bahasa Indonesia", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"id-ID-Neural2-A"}},
		{Code: "ms", Name: "Malay", NativeName: "Bahasa Melayu", SupportTier: TierSupported, SupportsPronunciation: true, TTSVoices: []string{"ms-MY-Neural2-A"}},
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, lang := range verified {
		r.languages[lang.Code] = lang
	}
	for _, lang := range supported {
		r.languages[lang.Code] = lang
	}
}

// Get returns language info for a given code.
// If the language is not in the registry, it returns an experimental entry.
// This allows ANY language code to be used - the LLM will handle it if supported.
func (r *Registry) Get(code string) *Info {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" {
		return nil
	}

	r.mu.RLock()
	if info, ok := r.languages[code]; ok {
		r.mu.RUnlock()
		return info
	}
	r.mu.RUnlock()

	// Return an experimental entry for unknown languages
	// This allows the system to try ANY language that LLM APIs might support
	return &Info{
		Code:                  code,
		Name:                  "Unknown (" + code + ")",
		NativeName:            code,
		SupportTier:           TierExperimental,
		SupportsPronunciation: true, // Let the API decide
		Notes:                 "This language is not pre-verified. Quality may vary depending on LLM/TTS/STT API support.",
	}
}

// Register adds or updates a language in the registry
func (r *Registry) Register(info *Info) {
	if info == nil || info.Code == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.languages[info.Code] = info
}

// GetAll returns all registered languages (not including dynamic experimental ones)
func (r *Registry) GetAll() []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Info, 0, len(r.languages))
	for _, info := range r.languages {
		result = append(result, info)
	}
	return result
}

// GetByTier returns languages filtered by support tier
func (r *Registry) GetByTier(tier SupportTier) []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Info, 0)
	for _, info := range r.languages {
		if info.SupportTier == tier {
			result = append(result, info)
		}
	}
	return result
}

// IsValidCode checks if a language code appears to be valid (ISO 639-1 format)
// Note: This does NOT reject the language - it only validates the format
func IsValidCode(code string) bool {
	code = strings.ToLower(strings.TrimSpace(code))
	// ISO 639-1 codes are 2 letters, ISO 639-2 are 3 letters
	if len(code) < 2 || len(code) > 3 {
		return false
	}
	for _, c := range code {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}
