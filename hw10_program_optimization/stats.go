package hw10programoptimization

import (
	"bufio"
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
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	keyAt := "@"
	u := User{}

	for scanner.Scan() {
		text := scanner.Text()

		if indx := strings.Index(text, domain); indx != -1 {
			u.UnmarshalJSON([]byte(text))

			if i := strings.Index(u.Email, domain); i == -1 {
				continue
			}

			from := strings.Index(u.Email, keyAt) + 1
			key := strings.ToLower(u.Email[from:])
			result[key]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
