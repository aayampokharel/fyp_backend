package block

import (
	"github.com/aayampokharel/fyp/common"
	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils"
	"github.com/aayampokharel/fyp/utils/global"
)

func InsertBlock(certificateData models.CertificateData) error {
	//! use MUTEX for global slice . for accidental causes  !!! its much easier  .
	if len(global.BlockChain) == 0 {
		global.BlockChain = append(global.BlockChain, utils.CreateGenesisBlock())
	}

	lastBlock := &global.BlockChain[len(global.BlockChain)-1]
	lastBlockCertificatesLength := common.CalculateLengthOfCertificateArray(lastBlock.CertificateData)
	if lastBlockCertificatesLength == 4 {
		global.BlockChain = append(global.BlockChain, utils.InsertEmptyBlock())
		lastBlock = &global.BlockChain[len(global.BlockChain)-1]
		lastBlockCertificatesLength = 0
	}
	lastBlock.CertificateData[lastBlockCertificatesLength] = certificateData

	merkelRootString, err := common.CalculateMerkelRoot(lastBlock.CertificateData)
	if err != nil {
		utils.LogErrorWithContext("Merkel Root Calculation", err)
		return ErrMerkleRootFailed
	}
	err = common.ProofOfWork(&lastBlock.Header)
	if err != nil {
		utils.LogErrorWithContext("Proof of Work", err)
		return ErrPoWFailed
	}
	if len(global.BlockChain) > 1 {
		lastBlock.Header.PreviousHash = global.BlockChain[len(global.BlockChain)-2].Header.CurrentHash
	} else {
		lastBlock.Header.PreviousHash = lastBlock.Header.CurrentHash
	}
	lastBlock.Header.MerkleRoot = merkelRootString
	return nil
}
