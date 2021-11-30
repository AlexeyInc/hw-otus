package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	slashDecCode = 92
	zeroDecCode  = 48
	nineDecCode  = 57
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	var resultStr string
	if isBacktickString(inputStr) {
		errorMsg := validateBacktickString(inputStr)
		if errorMsg != nil {
			return "", errorMsg
		}
		resultStr = unpackBacktickString(inputStr)
	} else {
		errorMsg := validateQuotedString(inputStr)
		if errorMsg != nil {
			return "", errorMsg
		}
		resultStr = unpackQuotedString(inputStr)
	}
	return resultStr, nil
}

func isBacktickString(str string) bool {
	return strings.ContainsAny(str, "\\")
}

func validateBacktickString(str string) error {
	for i := range str {
		if str[i] == '\\' && i+1 < len(str) {
			shieldingChar := str[i+1]

			if (shieldingChar < zeroDecCode || shieldingChar > nineDecCode) && shieldingChar != slashDecCode {
				return ErrInvalidString
			}
		}
	}
	return nil
}

func unpackBacktickString(input string) string {
	var result strings.Builder
	inputStr := []rune(input)
	var tempStr string
	for i := 0; i < len(inputStr); i++ {
		switch {
		case inputStr[i] == slashDecCode:
			tempStr = makeShielding(i, inputStr)
			i++
		case unicode.IsDigit(inputStr[i]):
			numRep, _ := strconv.Atoi(string(inputStr[i]))
			tempStr = strings.Repeat(string(inputStr[i-1]), numRep)
		default:
			tempStr = string(inputStr[i])
		}
		if isNextCharNotDigit(i, inputStr) {
			result.WriteString(tempStr)
		}
	}
	return result.String()
}

func makeShielding(i int, str []rune) string {
	if i+1 > len(str)-1 {
		return ""
	}
	return string(str[i+1])
}

func validateQuotedString(str string) error {
	if len(str) == 0 {
		return nil
	}
	if unicode.IsDigit([]rune(str)[0]) {
		return ErrInvalidString
	}
	if !requireWithoutNumbers(str) {
		return ErrInvalidString
	}
	return nil
}

func unpackQuotedString(input string) string {
	var resultStr strings.Builder
	inputStr := []rune(input)
	for i, v := range inputStr {
		if unicode.IsDigit(v) {
			numRep, _ := strconv.Atoi(string(inputStr[i]))
			strRep := string(inputStr[i-1])
			resultStr.WriteString(strings.Repeat(strRep, numRep))
		} else if isNextCharNotDigit(i, inputStr) {
			resultStr.WriteString(string(inputStr[i]))
		}
	}
	return resultStr.String()
}

func requireWithoutNumbers(str string) bool {
	strArr := []rune(str)
	countDigit := 0
	for i := 0; i < len(strArr); i++ {
		if unicode.IsDigit(strArr[i]) {
			countDigit++
		} else {
			countDigit = 0
		}
		if countDigit > 1 {
			return false
		}
	}
	return true
}

func isNextCharNotDigit(i int, str []rune) bool {
	if i+1 > len(str)-1 {
		return true
	}
	return !unicode.IsDigit(str[i+1])
}
