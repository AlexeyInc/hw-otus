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
	keyEmail := "Email\":\""
	keyComma := "\",\""

	for scanner.Scan() {
		text := scanner.Text()

		if indx := strings.Index(text, keyEmail); indx != -1 {
			indxStart := indx + +len(keyEmail)
			offsetEmail := strings.Index(text[indxStart:], keyComma)
			emailText := text[indxStart : indxStart+offsetEmail]

			if strings.Contains(emailText, domain) {
				from := strings.Index(emailText, keyAt) + 1
				key := strings.ToLower(emailText[from:])
				result[key]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
