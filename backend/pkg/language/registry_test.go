package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegistry(t *testing.T) {
	registry := GetRegistry()
	assert.NotNil(t, registry)

	// Should return the same instance
	registry2 := GetRegistry()
	assert.Equal(t, registry, registry2)
}

func TestRegistry_Get_VerifiedLanguage(t *testing.T) {
	registry := GetRegistry()

	// Japanese is a verified language
	ja := registry.Get("ja")
	assert.NotNil(t, ja)
	assert.Equal(t, "ja", ja.Code)
	assert.Equal(t, "Japanese", ja.Name)
	assert.Equal(t, "日本語", ja.NativeName)
	assert.Equal(t, TierVerified, ja.SupportTier)
	assert.True(t, ja.SupportsPronunciation)
}

func TestRegistry_Get_SupportedLanguage(t *testing.T) {
	registry := GetRegistry()

	// Korean is a supported language
	ko := registry.Get("ko")
	assert.NotNil(t, ko)
	assert.Equal(t, "ko", ko.Code)
	assert.Equal(t, "Korean", ko.Name)
	assert.Equal(t, TierSupported, ko.SupportTier)
}

func TestRegistry_Get_UnknownLanguage(t *testing.T) {
	registry := GetRegistry()

	// Kurdish (ku) is not in the registry but should work as experimental
	ku := registry.Get("ku")
	assert.NotNil(t, ku)
	assert.Equal(t, "ku", ku.Code)
	assert.Equal(t, TierExperimental, ku.SupportTier)
	assert.Contains(t, ku.Notes, "not pre-verified")

	// Esperanto (eo) - also experimental
	eo := registry.Get("eo")
	assert.NotNil(t, eo)
	assert.Equal(t, "eo", eo.Code)
	assert.Equal(t, TierExperimental, eo.SupportTier)

	// Any valid language code should return experimental
	xyz := registry.Get("xy")
	assert.NotNil(t, xyz)
	assert.Equal(t, TierExperimental, xyz.SupportTier)
}

func TestRegistry_Get_CaseInsensitive(t *testing.T) {
	registry := GetRegistry()

	// Should work regardless of case
	ja1 := registry.Get("JA")
	ja2 := registry.Get("ja")
	ja3 := registry.Get("Ja")

	assert.Equal(t, ja1.Code, ja2.Code)
	assert.Equal(t, ja2.Code, ja3.Code)
}

func TestRegistry_Get_EmptyCode(t *testing.T) {
	registry := GetRegistry()

	empty := registry.Get("")
	assert.Nil(t, empty)

	spaces := registry.Get("   ")
	assert.Nil(t, spaces)
}

func TestRegistry_GetAll(t *testing.T) {
	registry := GetRegistry()

	all := registry.GetAll()
	assert.NotEmpty(t, all)

	// Should have at least the verified languages
	assert.GreaterOrEqual(t, len(all), 9) // 9 verified languages

	// Check that Japanese is in the list
	found := false
	for _, lang := range all {
		if lang.Code == "ja" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestRegistry_GetByTier(t *testing.T) {
	registry := GetRegistry()

	verified := registry.GetByTier(TierVerified)
	assert.NotEmpty(t, verified)

	// All returned should be verified
	for _, lang := range verified {
		assert.Equal(t, TierVerified, lang.SupportTier)
	}

	supported := registry.GetByTier(TierSupported)
	assert.NotEmpty(t, supported)

	// All returned should be supported
	for _, lang := range supported {
		assert.Equal(t, TierSupported, lang.SupportTier)
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry() // Use fresh registry for this test

	// Register a new language
	registry.Register(&Info{
		Code:                  "test",
		Name:                  "Test Language",
		NativeName:            "Test",
		SupportTier:           TierExperimental,
		SupportsPronunciation: true,
	})

	lang := registry.Get("test")
	assert.NotNil(t, lang)
	assert.Equal(t, "Test Language", lang.Name)
}

func TestIsValidCode(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{"ja", true},
		{"en", true},
		{"JA", true},  // Case insensitive check happens in Get()
		{"eng", true}, // ISO 639-2 (3 letters)
		{"jpn", true},
		{"a", false},     // Too short
		{"abcd", false},  // Too long
		{"12", false},    // Numbers
		{"j1", false},    // Mixed
		{"", false},      // Empty
		{"ja-JP", false}, // Locale format (not pure language code)
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := IsValidCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDynamicLanguageSupport(t *testing.T) {
	registry := GetRegistry()

	// The key feature: ANY language code should return something usable
	// This allows users to try any language their LLM supports

	minorLanguages := []string{
		"ku",  // Kurdish
		"eo",  // Esperanto
		"cy",  // Welsh
		"ga",  // Irish
		"mt",  // Maltese
		"is",  // Icelandic
		"sq",  // Albanian
		"mk",  // Macedonian
		"sr",  // Serbian
		"hr",  // Croatian
		"sl",  // Slovenian
		"sk",  // Slovak
		"bg",  // Bulgarian
		"lt",  // Lithuanian
		"lv",  // Latvian
		"et",  // Estonian
		"ka",  // Georgian
		"hy",  // Armenian
		"az",  // Azerbaijani
		"uz",  // Uzbek
		"kk",  // Kazakh
		"tg",  // Tajik
		"mn",  // Mongolian
		"my",  // Burmese
		"km",  // Khmer
		"lo",  // Lao
		"si",  // Sinhala
		"ne",  // Nepali
		"bn",  // Bengali
		"ta",  // Tamil
		"te",  // Telugu
		"ml",  // Malayalam
		"kn",  // Kannada
		"gu",  // Gujarati
		"mr",  // Marathi
		"pa",  // Punjabi
		"ur",  // Urdu
		"sw",  // Swahili
		"am",  // Amharic
		"yo",  // Yoruba
		"ig",  // Igbo
		"zu",  // Zulu
		"af",  // Afrikaans
	}

	for _, code := range minorLanguages {
		t.Run(code, func(t *testing.T) {
			lang := registry.Get(code)
			assert.NotNil(t, lang, "Language %s should return non-nil", code)
			assert.Equal(t, code, lang.Code)
			// Either it's in our known list or it's experimental
			assert.Contains(t, []SupportTier{TierVerified, TierSupported, TierExperimental}, lang.SupportTier)
		})
	}
}
