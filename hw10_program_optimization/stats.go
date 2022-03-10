package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const newLineByte = 10

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

	countUsers := countByte(content, newLineByte)
	result = make([]User, countUsers)

	resContent := "[" + strings.ReplaceAll(string(content), "\n", ",") + "]"

	if err = json.Unmarshal([]byte(resContent), &result); err != nil {
		return
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

func countByte(input []byte, search byte) (count int) {
	for _, v := range input {
		if v == search {
			count++
		}
	}
	return
}
