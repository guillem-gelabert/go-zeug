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
