package entity

type RetrieveBlockRequest struct {
	Id string `json:"id"`
}

type RetrievedBlockResponse struct {
	Status       int    `json:"status"`
	CurrentBlock Block  `json:"block"`
	Message      string `json:"message"`
}
