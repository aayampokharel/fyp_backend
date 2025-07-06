package block

import "errors"

var (
	ErrInvalidCertificate = errors.New("invalid certificate data")
	ErrMerkleRootFailed   = errors.New("failed to compute Merkle root")
	ErrPoWFailed          = errors.New("proof of work failed")
)
