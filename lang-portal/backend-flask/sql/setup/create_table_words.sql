CREATE TABLE IF NOT EXISTS words (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  kanji TEXT NOT NULL,
  romaji TEXT NOT NULL,  -- This is currently romaji, not kana
  english TEXT NOT NULL, -- This is currently english, not meaning
  parts TEXT NOT NULL   -- JSON string of word parts
);