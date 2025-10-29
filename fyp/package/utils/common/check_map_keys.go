package common

import (
	err "project/package/errors"
	"strconv"
)

func CheckMapKeysReturnValues(providedMap map[string]string, keys []string) (map[string]string, error) {
	for _, key := range keys {
		providedMapValue, exists := providedMap[key]
		if !exists {
			return nil, err.ErrWithMoreInfo(nil, "key doesnot exist")
		}
		if providedMapValue == "" {
			return nil, err.ErrWithMoreInfo(nil, key+"is required")
		}

	}
	return providedMap, nil
}

func ConvertToBool(boolStr string) (bool, error) {
	if boolStr == "" {
		return false, err.ErrEmptyString
	}
	val, er := strconv.ParseBool(boolStr)
	if er != nil {
		return false, err.ErrCannotConvertToBool
	}
	return val, nil
}

func ConvertToInt(intStr string) (int, error) {
	if intStr == "" {
		return 0, err.ErrEmptyString
	}
	val, er := strconv.Atoi(intStr)
	if er != nil {
		return 0, err.ErrCannotConvertToInt
	}
	return int(val), nil
}
