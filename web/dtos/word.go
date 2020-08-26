package dto

// WordDTO represents a Word as returned to the client
type WordDTO struct {
	ID          int    `json:"id"`
	Article     string `json:"article"`
	Substantive string `json:"substantive"`
	// Priority represents the frequency of a word, low priority value means high frequency
	Priority int `json:"priority"`
}
