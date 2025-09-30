package seed

import (
	"crypto/rand"
	"math/big"
)

func generateRandomNumber(upperBound int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(upperBound)))
	return int(num.Int64())
}
