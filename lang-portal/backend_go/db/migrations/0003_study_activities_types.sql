-- Add tables for sentence construction
CREATE TABLE IF NOT EXISTS sentences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    japanese TEXT NOT NULL,
    english TEXT NOT NULL,
    words TEXT NOT NULL,  -- JSON array of words for construction
    hints TEXT,          -- JSON array of hints
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sentences_groups (
    sentence_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (sentence_id, group_id),
    FOREIGN KEY (sentence_id) REFERENCES sentences(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS sentence_review_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    sentence_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (sentence_id) REFERENCES sentences(id)
);

-- Add indices for better performance
CREATE INDEX idx_sentences_groups_group ON sentences_groups(group_id);
CREATE INDEX idx_sentence_review_items_session ON sentence_review_items(study_session_id);
CREATE INDEX idx_sentence_review_correct ON sentence_review_items(correct);

-- Add columns for activity-specific scores
ALTER TABLE study_activities ADD COLUMN accuracy_score FLOAT;
ALTER TABLE study_activities ADD COLUMN confidence_score FLOAT;

-- Add indices for activity types and scores
CREATE INDEX idx_study_activities_accuracy ON study_activities(accuracy_score);
CREATE INDEX idx_study_activities_confidence ON study_activities(confidence_score); 