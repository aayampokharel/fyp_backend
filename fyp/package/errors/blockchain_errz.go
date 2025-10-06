package err

import (
	"fmt"
	errorz "project/package/error_blueprint"
	"strings"
)

func ErrWithMoreInfo(err error, erStr ...string) error {
	mergedErStr := strings.Join(erStr, ", ")
	if err == nil {
		return fmt.Errorf("error: %v", mergedErStr)
	}
	return fmt.Errorf("%w: %v", err, mergedErStr)
}

var (
	ErrEmptyBlockChain      = errorz.Status400BadRequest.Wrap("Empty Blockchain")
	ErrGenesisBlockUpdate   = errorz.Status400BadRequest.Wrap("Cannot update Genesis Block")
	ErrBlockNumberMismatch  = errorz.Status400BadRequest.Wrap("Block number mismatch")
	ErrGenesisBlockInsert   = errorz.Status400BadRequest.Wrap("Cannot insert genesis block at index other than 0")
	ErrNotEnoughBlocks      = errorz.Status400BadRequest.Wrap("Not enough blocks")
	ErrArrayOutOfBound      = errorz.Status406NotAcceptable.Wrap("index not within the range of 0-3")
	ErrInvalidBlockNumber   = errorz.Status400BadRequest.Wrap("Invalid Block Number")
	ErrInvalidHash          = errorz.Status400BadRequest.Wrap("Invalid Hash")
	ErrGenesisBlockMismatch = errorz.Status400BadRequest.Wrap("Genesis Block Mismatch")
	ErrMarshaling           = errorz.Status400BadRequest.Wrap("Error marshaling")
	ErrEmptyFields          = errorz.Status400BadRequest.Wrap("Empty Fields")
	ErrEmptyPOWRules        = errorz.Status400BadRequest.Wrap("Empty POW Rule String in .env file")
	ErrEnvParsing           = errorz.Status400BadRequest.Wrap("Error parsing .env file")
	ErrIntParse             = errorz.Status400BadRequest.Wrap("Error parsing int from srring")
	ErrTcpListen            = errorz.Status400BadRequest.Wrap("Error in TCP Listen")
)
