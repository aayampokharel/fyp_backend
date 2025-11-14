package service

import (
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
	err "project/package/errors"
	"strconv"
)

func (s *Service) VerifyBlock(currentBlock, previousBlock entity.Block) error {
	// 	// ⏺️ Previous Block Hash - Chain continuity

	calculatedMerkleRoot, er := s.CalculateMerkleRoot(currentBlock.CertificateData)
	if er != nil {
		return er
	}
	if calculatedMerkleRoot != currentBlock.Header.MerkleRoot {
		s.Logger.Errorln("[verify_block_chain] Error: calculateMerkleRoot Error::", er)
		return err.ErrVerifyingBlock
	}

	powStructureParam, er := entity.NewPowStructure(currentBlock.Header.BlockNumber, previousBlock.Header.PreviousHash, calculatedMerkleRoot, currentBlock.Header.TimeStamp)
	if er != nil {
		s.Logger.Errorln("[verify_block_chain] Error: NewPowStructure Error::", er)
		return err.ErrVerifyingBlock
	}

	env, er := config.NewEnv()
	if er != nil {
		s.Logger.Errorln(er)
		return er
	}
	powRuleString := env.GetValueForKey(constants.PowNumberRuleKey)
	nonce, powPuzzleHash, er := s.CalculatePOW(powStructureParam, powRuleString)
	if er != nil {
		s.Logger.Errorln("[verify_block_chain] Error: CalculatePOW error::", er)
		return err.ErrVerifyingBlock
	}

	if strconv.Itoa(nonce) != currentBlock.Header.Nonce || powPuzzleHash != currentBlock.Header.CurrentHash {
		s.Logger.Errorln("[verify_block_chain] Error: pow_verification error::", er)
		return err.ErrVerifyingBlock
	}

	if currentBlock.Header.PreviousHash != previousBlock.Header.CurrentHash {

		return err.ErrVerifyingBlock
	}
	return nil

}
