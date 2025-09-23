package entity

type Node struct {
	CurrentPort              string            `json:"current_ip"`
	MemQueue                 []CertificateData `json:"mem_queue"`
	VerifyMemQueue           []Block           `json:"verify_mem_queue"`
	BlockChain               []Block           `json:"block_chain"`
	TempBlockForVerification Block             `json:"temp_queue"`
}
