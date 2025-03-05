CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY,
    group_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id)
); 