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
	ErrEmptyBlockChain              = errorz.Status400BadRequest.Wrap("Empty Blockchain")
	ErrEmptyInstitutionInfo         = errorz.Status400BadRequest.Wrap("Institution info cannot be empty")
	ErrEmptyUserEmail               = errorz.Status400BadRequest.Wrap("User email cannot be empty")
	ErrGenesisBlockUpdate           = errorz.Status400BadRequest.Wrap("Cannot update Genesis Block")
	ErrBlockNumberMismatch          = errorz.Status400BadRequest.Wrap("Block number mismatch")
	ErrNumberMismatch               = errorz.Status400BadRequest.Wrap("Number mismatch")
	ErrInvalidType                  = errorz.Status400BadRequest.Wrap("Invalid type")
	ErrGenesisBlockInsert           = errorz.Status400BadRequest.Wrap("Cannot insert genesis block at index other than 0")
	ErrNotEnoughBlocks              = errorz.Status400BadRequest.Wrap("Not enough blocks")
	ErrArrayOutOfBound              = errorz.Status406NotAcceptable.Wrap("index not within the range of 0-3")
	ErrLeaderPort                   = errorz.Status406NotAcceptable.Wrap("leader port wasnot initialized")
	ErrUserDoesnotExist             = errorz.Status400BadRequest.Wrap("User does not exist")
	ErrPythonScriptReturnedEmpty    = errorz.Status400BadRequest.Wrap("python script returned empty result")
	ErrInstitutionAlreadyVerified   = errorz.Status406NotAcceptable.Wrap("institution is already verified.")
	ErrInvalidBlockNumber           = errorz.Status400BadRequest.Wrap("Invalid Block Number")
	ErrInvalidHash                  = errorz.Status400BadRequest.Wrap("Invalid Hash")
	ErrInvalidBase64                = errorz.Status400BadRequest.Wrap("Invalid base64 string")
	ErrGenesisBlockMismatch         = errorz.Status400BadRequest.Wrap("Genesis Block Mismatch")
	ErrMarshaling                   = errorz.Status400BadRequest.Wrap("Error marshaling")
	ErrEmptyFields                  = errorz.Status400BadRequest.Wrap("Empty Fields")
	ErrEmptyPOWRules                = errorz.Status400BadRequest.Wrap("Empty POW Rule String in .env file")
	ErrEnvParsing                   = errorz.Status400BadRequest.Wrap("Error parsing .env file")
	ErrFileParsing                  = errorz.Status400BadRequest.Wrap("Error parsing file")
	ErrFileExecuting                = errorz.Status400BadRequest.Wrap("Error executing file")
	ErrIntParse                     = errorz.Status400BadRequest.Wrap("Error parsing int from srring")
	ErrTcpListen                    = errorz.Status400BadRequest.Wrap("Error in TCP Listen")
	ErrInstitutionAlreadyRegistered = errorz.Status400BadRequest.Wrap("Institution already registered")
	ErrCannotConvertToBool          = errorz.Status400BadRequest.Wrap("cannot convert to bool type")
	ErrCannotConvertToInt           = errorz.Status400BadRequest.Wrap("cannot convert to int type")
	ErrCannotConvertToFloat         = errorz.Status400BadRequest.Wrap("cannot convert to floattype")
	ErrEmptyString                  = errorz.Status400BadRequest.Wrap("empty string provided")
	ErrWritingZip                   = errorz.Status400BadRequest.Wrap("Error while writing in zip")
	ErrPDFDataIsNil                 = errorz.Status400BadRequest.Wrap("PDF Data is empty")
	ErrClosingZipWriter             = errorz.Status400BadRequest.Wrap("Error closing zip writer")
)
