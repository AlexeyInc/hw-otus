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

func unpackBacktickString(inputStr string) string {
	var result strings.Builder
	i := 0
	for ; i < len(inputStr); i++ {
		if inputStr[i] == slashDecCode {
			addNewStr := makeShielding(i, inputStr)
			i++
			result.WriteString(addNewStr)
		} else if !unicode.IsDigit([]rune(inputStr)[i]) {
			result.WriteString(string(inputStr[i]))
		}
	}
	return result.String()
}

func makeShielding(i int, str string) string {
	if i+1 > len(str)-1 {
		return ""
	}
	if i+2 > len(str)-1 || !unicode.IsDigit([]rune(str)[i+2]) {
		return string(str[i+1])
	}
	var resultStr string
	if numRep, err := strconv.Atoi(string(str[i+2])); err == nil {
		resultStr = strings.Repeat(string(str[i+1]), numRep)
	} else {
		return ""
	}
	return resultStr
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

func unpackQuotedString(inputStr string) string {
	var resultStr strings.Builder
	for i, v := range inputStr {
		if unicode.IsDigit(v) {
			if numRep, err := strconv.Atoi(string(inputStr[i])); err == nil {
				strRep := string(inputStr[i-1])
				resultStr.WriteString(strings.Repeat(strRep, numRep))
			}
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

func isNextCharNotDigit(i int, str string) bool {
	if i+1 > len(str)-1 {
		return true
	}
	return !unicode.IsDigit([]rune(str)[i+1])
}
