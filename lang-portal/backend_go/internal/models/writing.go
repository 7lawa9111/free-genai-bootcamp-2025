package models

type WritingPractice struct {
    ID        int64           `json:"id"`
    GroupID   int64          `json:"group_id"`
    Exercises []WriteExercise `json:"exercises"`
}

type WriteExercise struct {
    ID       int64    `json:"id"`
    Japanese string   `json:"japanese"`
    English  string   `json:"english"`
    Romaji   string   `json:"romaji"`
    Hints    []string `json:"hints,omitempty"`
}

type WritingResult struct {
    ActivityID      int64   `json:"activity_id"`
    WordID          int64   `json:"word_id"`
    WrittenText     string  `json:"written_text"`
    Accuracy        float64 `json:"accuracy"`
    ExercisesDone   int     `json:"exercises_done"`
    TotalExercises  int     `json:"total_exercises"`
    TimeTaken       float64 `json:"time_taken"`
} 