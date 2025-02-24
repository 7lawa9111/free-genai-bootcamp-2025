-- Create words table
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    japanese TEXT NOT NULL,
    romaji TEXT NOT NULL,
    english TEXT NOT NULL,
    parts TEXT
);

-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- Create words_groups table
CREATE TABLE IF NOT EXISTS words_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER,
    group_id INTEGER,
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Create study_sessions table
CREATE TABLE IF NOT EXISTS study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    study_activity_id INTEGER,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Create study_activities table
CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER,
    group_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Create word_review_items table
CREATE TABLE IF NOT EXISTS word_review_items (
    word_id INTEGER,
    study_session_id INTEGER,
    correct BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
); 