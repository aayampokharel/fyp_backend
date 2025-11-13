package service

import (
	"project/internals/domain/entity"
	err "project/package/errors"
	"project/package/utils/common"
	"strconv"
)

func (s *Service) CalculatePOW(powParams entity.PowStructure, powRuleString string) (nonce int, currentHash string, er error) {

	if powParams.BlockMerkleRoot == "" || powParams.PreviousHash == "" || powParams.BlockNumber == 0 {
		return -1, "", err.ErrEmptyFields

	}
	powRuleLengthString := len(powRuleString)
	if powRuleLengthString == 0 {
		return -1, "", err.ErrEmptyPOWRules
	}
	hashedPowParams, _, err := common.HashData(powParams)
	if err != nil {
		return -1, "", err
	}

	for nonce := 0; ; nonce++ {
		hashedVal, _, err := common.HashData(hashedPowParams + strconv.Itoa(nonce))
		if err != nil {
			return -1, "", err
		}
		if hashedVal[:powRuleLengthString] == powRuleString {
			return nonce, hashedVal, nil
		}

	}
}
