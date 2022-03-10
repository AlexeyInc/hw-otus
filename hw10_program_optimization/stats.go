package hw10programoptimization

import (
	"fmt"
	"io"
	"io/ioutil"
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

func getUsers(r io.Reader) (result []User, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	result = make([]User, len(lines))

	for i, line := range lines {
		if err = result[i].UnmarshalJSON([]byte(line)); err != nil {
			return
		}
	}

	return
}

func countDomains(u []User, domain string) (DomainStat, error) {
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
