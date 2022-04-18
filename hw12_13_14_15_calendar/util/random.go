package util

import (
	mathRand "math/rand"
	"strings"
	"time"
)

const (
	titleLen = 6
	descLen  = 20
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func RandomInt(max int) int {
	mathRand.Seed(time.Now().UnixNano())
	return mathRand.Intn(max)
}

func RandomIntRange(min, max int) int {
	mathRand.Seed(time.Now().UnixNano())
	return (mathRand.Intn(max-min+1) + min)
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

func RandomUserID() int64 { // Remove after add User logic
	users := []int64{1, 2}
	n := len(users)
	return users[RandomInt(n)]
}
