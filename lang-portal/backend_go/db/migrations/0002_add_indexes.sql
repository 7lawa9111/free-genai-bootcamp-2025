-- Add indexes for frequently queried columns and foreign keys

-- Words table indexes
CREATE INDEX IF NOT EXISTS idx_words_japanese ON words(japanese);
CREATE INDEX IF NOT EXISTS idx_words_romaji ON words(romaji);

-- Groups table indexes
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);

-- Words_groups table indexes
CREATE INDEX IF NOT EXISTS idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX IF NOT EXISTS idx_words_groups_group_id ON words_groups(group_id);

-- Study sessions table indexes
CREATE INDEX IF NOT EXISTS idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_created_at ON study_sessions(created_at);
CREATE INDEX IF NOT EXISTS idx_study_sessions_activity_id ON study_sessions(study_activity_id);

-- Study activities table indexes
CREATE INDEX IF NOT EXISTS idx_study_activities_group_id ON study_activities(group_id);
CREATE INDEX IF NOT EXISTS idx_study_activities_created_at ON study_activities(created_at);

-- Word review items table indexes
CREATE INDEX IF NOT EXISTS idx_word_review_items_word_id ON word_review_items(word_id);
CREATE INDEX IF NOT EXISTS idx_word_review_items_session_id ON word_review_items(study_session_id);
CREATE INDEX IF NOT EXISTS idx_word_review_items_created_at ON word_review_items(created_at); 