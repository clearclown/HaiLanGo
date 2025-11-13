package models

import "time"

// WordDefinition represents a single definition of a word
type WordDefinition struct {
	Definition string   `json:"definition"`
	Examples   []string `json:"examples,omitempty"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

// WordMeaning represents a meaning of a word with a specific part of speech
type WordMeaning struct {
	PartOfSpeech string           `json:"partOfSpeech"`
	Definitions  []WordDefinition `json:"definitions"`
}

// WordPhonetic represents pronunciation information
type WordPhonetic struct {
	Text     string `json:"text"`               // Phonetic text (e.g., "/həˈloʊ/")
	AudioURL string `json:"audioUrl,omitempty"` // URL to audio pronunciation
}

// WordEntry represents a complete dictionary entry for a word
type WordEntry struct {
	Word      string         `json:"word"`
	Phonetics []WordPhonetic `json:"phonetics,omitempty"`
	Meanings  []WordMeaning  `json:"meanings"`
	Language  string         `json:"language,omitempty"`
	SourceAPI string         `json:"sourceApi,omitempty"` // Which API provided this data
	FetchedAt time.Time      `json:"fetchedAt"`
}
