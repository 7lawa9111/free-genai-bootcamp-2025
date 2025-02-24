CREATE INDEX idx_study_activities_type ON study_activities(type);
CREATE INDEX idx_study_activities_completed ON study_activities(completed_at);
CREATE INDEX idx_study_activities_group_id ON study_activities(group_id);
CREATE INDEX idx_study_sessions_activity_id ON study_sessions(study_activity_id);
CREATE INDEX idx_quiz_answers_session_id ON vocabulary_quiz_answers(study_session_id); 