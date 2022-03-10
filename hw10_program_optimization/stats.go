package hw10programoptimization

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pquerna/ffjson/ffjson"
)

const newLineByte = 10

type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
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
	result = make([]User, 0, countUsers)

	resContent := "[" + strings.ReplaceAll(string(content), "\n", ",") + "]"

	if err = ffjson.Unmarshal([]byte(resContent), &result); err != nil {
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
