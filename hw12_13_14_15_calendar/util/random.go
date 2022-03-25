package util

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	titleLen = 6
	descLen  = 20
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func RandomInt(max int) int {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(randNum.Int64())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[RandomInt(k)]
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

func RandomUserID() int64 { // TODO: Remove after add User logic
	users := []int64{1, 2}
	n := len(users)
	return users[RandomInt(n)]
}
