package models

type VocabularyQuiz struct {
    ID        int64  `json:"id"`
    GroupID   int64  `json:"group_id"`
    Questions []Word `json:"questions"`
}

type QuizQuestion struct {
    ID      int64    `json:"id"`
    Word    string   `json:"word"`
    Correct string   `json:"correct"`
    Options []string `json:"options"`
}

type QuizResult struct {
    ActivityID   int64   `json:"activity_id"`
    WordID      int64   `json:"word_id"`
    Answer      string  `json:"answer"`
    Correct     bool    `json:"correct"`
    Score       float64 `json:"score"`
    CorrectCount int    `json:"correct_count"`
    TotalCount   int    `json:"total_count"`
    TimeTaken    float64 `json:"time_taken"`
} 