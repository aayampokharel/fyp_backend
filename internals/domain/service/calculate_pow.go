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
	hashedPowParams, err := common.HashData(powParams)
	if err != nil {
		return -1, err
	}

	nonce := 0
	for {
		hashedVal, err := common.HashData(hashedPowParams + strconv.Itoa(nonce))
		if err != nil {
			return -1, err
		}
		if hashedVal[:4] == powRuleString {
			return nonce, nil
		}
		nonce++

	}
}
