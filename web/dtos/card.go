package dto

// CardDTO represents a Word as returned to the client
type CardDTO struct {
	ID          int    `json:"id"`
	Article     string `json:"article"`
	Substantive string `json:"substantive"`
	WordID      int    `json:"word_id"`
	UserID      int    `json:"user_id"`
	Stage       string `json:"stage"`
}

// AnswerDTO represents an Answer to a Card as received from the client
type AnswerDTO struct {
	ID      int  `json:"id"`
	Correct bool `json:"correct"`
}
