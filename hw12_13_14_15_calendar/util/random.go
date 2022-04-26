package util

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

const (
	titleLen = 6
	descLen  = 20
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func RandomInt(max int64) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Fatal(err)
	}
	res := nBig.Int64()
	return int(res)
}

func RandomIntRange(min, max int64) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		log.Fatal(err)
	}
	return int(nBig.Int64() + min)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[RandomInt(int64(k))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomTitle() string {
	return RandomString(titleLen)
}

func RandomDescription() string {
	return RandomString(descLen)
}

func RandomUserID() int64 {
	users := []int64{1, 2}
	n := len(users)
	return users[RandomInt(int64(n))]
}
