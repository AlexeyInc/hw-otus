package hw10programoptimization

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	var usr users

	buf := bytes.Buffer{}
	buf.Grow(len(usr))

	if _, err := io.Copy(&buf, r); err != nil {
		return usr, err
	}

	textBytes := buf.Bytes()

	lineIndxs := [len(usr)]struct {
		start, end int
	}{}

	lineIndxs[0].start = 0
	j := 0
	for i, v := range textBytes {
		if v == 10 {
			lineIndxs[j].end = i
			j++
			if i < len(textBytes) {
				lineIndxs[j].start = i + 1
			}
		}
	}

	for i := 0; i < len(usr); i++ {
		startLineIndx := lineIndxs[i].start
		endLineIndx := lineIndxs[i].end
		if endLineIndx != 0 {
			usr[i].UnmarshalJSON(textBytes[startLineIndx:endLineIndx])
			continue
		}
		usr[i].UnmarshalJSON(textBytes[startLineIndx:])
		return usr, nil
	}

	return usr, nil
}

func countDomains(u [100000]User, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		contain := strings.Contains(user.Email, "."+domain)

		if contain {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
