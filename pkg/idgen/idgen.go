package idgen

import (
	"crypto/rand"
	"math/big"
)

type Number interface {
	int | int64
}

func IDgenerator() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return int(n.Int64()) + 100000
}
