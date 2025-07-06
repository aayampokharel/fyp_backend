package block

import (
	"github.com/aayampokharel/fyp/common"
	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils"
	"github.com/aayampokharel/fyp/utils/global"
)

func InsertBlock(certificateData models.CertificateData) error {
	// global array check garyo
	// if array is empty insert directly there , recalcualte the headers thing .
	// else call evaluete_new_header, and do stuffs . and insert the data there .
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

	merkelRootString, err := utils.CalculateMerkelRoot(lastBlock.CertificateData[:])
	if err != nil {
		utils.LogErrorWithContext("Merkel Root Calculation", err)
		return err
	}
	err = common.ProofOfWork(&lastBlock.Header)
	if err != nil {
		utils.LogErrorWithContext("Proof of Work", err)
		return err
	}
	lastBlock.Header.MerkleRoot = merkelRootString

	return nil
}
