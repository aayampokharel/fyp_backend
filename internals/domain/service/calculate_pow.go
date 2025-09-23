package service

import (
	"os"
	"project/internals/domain/entity"
	err "project/package/errors"
	"project/package/utils/common"
	"strconv"
)

func CalculatePOW(powParams entity.PowStructure) (int, error) {

	if powParams.BlockMerkleRoot == "" || powParams.PreviousHash == "" || powParams.BlockNumber == 0 {
		return -1, err.ErrEmptyFields

	}
	powRuleString := os.Getenv("POW_NUMBER_RULE")
	powRuleLengthString := len(powRuleString)
	nonce := 0
	if powRuleLengthString == 0 {
		return -1, err.ErrEmptyPOWRules
	}
	hashedPowParams, err := common.HashData(powParams)
	if err != nil {
		return -1, err
	}

	for {
		hashedVal, err := common.HashData(hashedPowParams + strconv.Itoa(nonce))
		if err != nil {
			return -1, err
		}
		if hashedVal[:powRuleLengthString] == powRuleString {
			return nonce, nil
		}
		nonce++

	}
}
