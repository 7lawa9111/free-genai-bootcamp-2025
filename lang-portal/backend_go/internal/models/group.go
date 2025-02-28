package models

type Group struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	WordCount int    `json:"word_count,omitempty"`
}

type GroupResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stats struct {
		TotalWordCount int `json:"total_word_count"`
	} `json:"stats"`
} 