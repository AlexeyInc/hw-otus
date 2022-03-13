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
	keyDomain := "." + domain

	for scanner.Scan() {
		text := scanner.Text()
		if indx := strings.Index(text, keyDomain); indx != -1 {
			from := strings.Index(text, keyAt) + 1
			to := indx + len(keyDomain)
			key := strings.ToLower(text[from:to])

			result[key]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
