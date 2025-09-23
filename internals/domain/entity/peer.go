package entity

type PeerResponse struct {
	Message string `json:"message"`
	Port    int    `json:"port"`
	Status  int    `json:"status"`
	//Signature string `json:"signature"` if required later on .
}
