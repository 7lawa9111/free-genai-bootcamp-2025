package models

type Word struct {
    ID       int    `json:"id"`
    Japanese string `json:"japanese"`
    Romaji   string `json:"romaji"`
    English  string `json:"english"`
}

type WordWithStats struct {
    Word
    CorrectCount int `json:"correct_count"`
    WrongCount   int `json:"wrong_count"`
}

type WordResponse struct {
    Japanese string `json:"japanese"`
    Romaji   string `json:"romaji"`
    English  string `json:"english"`
    Stats    struct {
        CorrectCount int `json:"correct_count"`
        WrongCount   int `json:"wrong_count"`
    } `json:"stats"`
    Groups []Group `json:"groups"`
} 