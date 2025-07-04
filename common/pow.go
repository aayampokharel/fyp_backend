package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils"
)

func ProofOfWork(header *models.Header) error {
	nonce := 0
	tempHeader := struct {
		BlockNumber  uint64 `json:"block_number"`
		TimeStamp    string `json:"timestamp"`
		PreviousHash string `json:"previous_hash"`
		MerkleRoot   string `json:"merkle_root"`
		Nonce        string `json:"nonce"`
	}{
		header.BlockNumber,
		header.TimeStamp.Format(time.RFC3339),
		header.PreviousHash,
		header.MerkleRoot,
		"",
	}
	for {
		header.Nonce = strconv.Itoa(nonce)
		hash, err := utils.CalculateHashHex(tempHeader)
		if err != nil {
			return fmt.Errorf("error in calculating hash: %w", err)
		}
		if hash[:4] == "0000" {
			break
		}
		nonce++
	}

	return nil
}
