package entity

type BroadcastStatusResponse struct {
	Status       int    `json:"status"`
	CurrentHash  string `json:"current_hash"`
	PreviousHash string `json:"previous_hash"`
}
