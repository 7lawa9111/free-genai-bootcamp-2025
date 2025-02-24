DROP TABLE IF EXISTS vocabulary_quiz_answers;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activities;

CREATE TABLE study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    activity_type TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    score REAL,
    settings TEXT,
    confidence_score FLOAT,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    study_activity_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
);

CREATE TABLE vocabulary_quiz_answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    answer TEXT NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);

CREATE INDEX idx_study_activities_type ON study_activities(activity_type);
CREATE INDEX idx_study_activities_completed ON study_activities(completed_at);
CREATE INDEX idx_study_activities_group_id ON study_activities(group_id);
CREATE INDEX idx_study_sessions_activity_id ON study_sessions(study_activity_id);
CREATE INDEX idx_quiz_answers_session_id ON vocabulary_quiz_answers(study_session_id); 