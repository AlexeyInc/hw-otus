package hw10programoptimization

import (
	json "encoding/json"
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

func getUsers(r io.Reader) (result users, err error) {
	dec := json.NewDecoder(r)

	var u User
	i := 0
	for dec.More() {
		err = dec.Decode(&u)
		if err != nil {
			return
		}
		result[i] = u
		i++
	}

	return
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
