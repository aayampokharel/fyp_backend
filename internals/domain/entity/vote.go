package entity

type VoteBody struct {
	CurrentHash string `json:"hash"`
	BoolValue   bool   `json:"bool_value"`
	PortInt     string `json:"port"`
}
