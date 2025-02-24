package models

type WordMatchingActivity struct {
	ID        int64    `json:"id"`
	GroupID   int64    `json:"group_id"`
	WordPairs []WordPair `json:"word_pairs"`
}

type WordPair struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	English  string `json:"english"`
	Matched  bool   `json:"matched"`
}

type WordMatchingResult struct {
	ActivityID   int64  `json:"activity_id"`
	Correct      bool   `json:"correct"`
	TimeTaken    int    `json:"time_taken_seconds"`
	MatchedPairs int    `json:"matched_pairs"`
	TotalPairs   int    `json:"total_pairs"`
} 