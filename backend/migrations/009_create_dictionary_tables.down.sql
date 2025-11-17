-- Drop Dictionary tables
DROP INDEX IF EXISTS idx_dictionary_cache_last_accessed;
DROP INDEX IF EXISTS idx_dictionary_cache_word_lang;

DROP TABLE IF EXISTS dictionary_supported_languages;
DROP TABLE IF EXISTS dictionary_cache;
