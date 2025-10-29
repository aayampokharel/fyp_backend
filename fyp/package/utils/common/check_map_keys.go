package common

import (
	err "project/package/errors"
	"strconv"
)

func CheckMapKeysReturnValues(providedMap map[string]string, keys ...string) map[string]string {
	for _, key := range keys {
		_, exists := providedMap[key]
		if !exists {
			return nil
		}
	}
	return providedMap
}

func CheckBoolFromString(boolStr string) (bool, error) {
	if boolStr == "" {
		return false, err.ErrEmptyString
	}
	val, er := strconv.ParseBool(boolStr)
	if er != nil {
		return false, err.ErrCannotConvertToBool
	}
	return val, nil
}

func CheckIntFromString(intStr string) (int, error) {
	if intStr == "" {
		return 0, err.ErrEmptyString
	}
	val, er := strconv.Atoi(intStr)
	if er != nil {
		return 0, err.ErrCannotConvertToInt
	}
	return int(val), nil
}
