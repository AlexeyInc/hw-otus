package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const validateTag string = "validate:"
const lenTag string = "len:"
const inTag string = "in:"
const regexpTag string = "regexp:"
const minTag string = "min:"
const maxTag string = "max:"
const separatorSymb string = "|"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (validationErrors ValidationErrors) Error() string {
	var errMsg strings.Builder

	for _, curErr := range validationErrors {
		errMsg.WriteString(
			fmt.Sprintf("Field Name: %v\nErrors: %v\n", curErr.Field, curErr.Err),
		)
	}
	return strings.TrimRight(errMsg.String(), "\n")
}

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)

	if reflectValue.Kind() != reflect.Struct {
		return nil
	}

	validationErrors := make(ValidationErrors, 0)
	reflectType := reflectValue.Type()
	numField := reflectValue.NumField()

	for i := 0; i < numField; i++ {
		fieldReflectType := reflectType.Field(i)
		fieldReflectValue := reflectValue.Field(i)
		fieldTag := fieldReflectType.Tag
		fieldTagStr := strings.Replace(string(fieldTag), "\\\\", "\\", -1)

		if requires := getRequirements(fieldTagStr); requires != "" {
			errors := getValidationErrors(fieldReflectType, fieldReflectValue, requires)
			if errors != nil {
				validationErrors = append(validationErrors, errors...)
			}
		}
	}

	return validationErrors
}

func getRequirements(require string) string {
	if !strings.Contains(require, validateTag) {
		return ""
	} else {
		fromIndx := strings.Index(require, validateTag)
		r := require[fromIndx+len(validateTag):]
		return r
	}
}

func getValidationErrors(fieldT reflect.StructField, fieldV reflect.Value, require string) []ValidationError {
	validationErros := make([]ValidationError, 0)
	require = strings.Trim(require, "\"")
	requirements := strings.Split(require, separatorSymb)

	isValidField := true
	var err error

	switch fieldV.Kind() {
	case reflect.String:
		isValidField, err = validString(fieldV.String(), requirements)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		isValidField, err = validInt(int(fieldV.Int()), requirements)
	case reflect.Slice:
		isValidField, err = validSlice(fieldV, requirements)
	default:
		fmt.Println("Can't define type of field") //TODO: remove after tests
	}
	if !isValidField {
		validationErros = append(validationErros, ValidationError{
			Field: fieldT.Name,
			Err:   err,
		})
	}
	return validationErros
}

func validString(str string, requirements []string) (bool, error) {
	isValidStr := true
	var resultErrors strings.Builder

	for _, req := range requirements {
		//check len:
		if indx := strings.Index(req, lenTag); indx != -1 {
			lenReqStr := req[indx+len(lenTag):]
			lenReq, err := strconv.Atoi(lenReqStr)
			if err != nil {
				isValidStr = false
				addErrMsg(&resultErrors, err.Error())
			}
			if len(str) != lenReq {
				isValidStr = false
				addErrMsg(&resultErrors, "string len should be: "+strconv.Itoa(lenReq))
			}
		}
		//check in:
		if indx := strings.Index(req, inTag); indx != -1 {
			inReqStr := req[indx+len(inTag):]
			inReqSlice := strings.Split(inReqStr, ",")
			exist := false
			for _, s := range inReqSlice {
				if str == s {
					exist = true
					break
				}
			}
			if !exist {
				isValidStr = false
				addErrMsg(&resultErrors, "value should be in: "+strings.Join(inReqSlice, ","))
			}
		}
		//check regexp:
		if indx := strings.Index(req, regexpTag); indx != -1 {
			pattern := req[indx+len(regexpTag):]
			match, _ := regexp.MatchString(pattern, str)
			if !match {
				isValidStr = false
				addErrMsg(&resultErrors, "string should match pattern: "+pattern)
			}
		}
	}
	if !isValidStr {
		return false, errors.New(resultErrors.String())
	}
	return true, nil
}

func validInt(num int, requirements []string) (bool, error) {
	isValidNum := true
	var resultErrors strings.Builder

	for _, req := range requirements {
		// check min:
		if indx := strings.Index(req, minTag); indx != -1 {
			minReqStr := req[indx+len(minTag):]
			minReq, err := strconv.Atoi(minReqStr)
			if err != nil {
				isValidNum = false
				addErrMsg(&resultErrors, err.Error())
			}
			if num < minReq {
				isValidNum = false
				addErrMsg(&resultErrors, "min should be: "+strconv.Itoa(minReq))
			}
			continue
		}
		//check max:
		if indx := strings.Index(req, maxTag); indx != -1 {
			maxReqStr := req[indx+len(maxTag):]
			maxReq, err := strconv.Atoi(maxReqStr)
			if err != nil {
				isValidNum = false
				addErrMsg(&resultErrors, err.Error())
			}
			if num > maxReq {
				isValidNum = false
				addErrMsg(&resultErrors, "max should be: "+strconv.Itoa(maxReq))
			}
			continue
		}
		//check in:
		matchIndxes := regexp.MustCompile("^" + inTag).FindStringIndex(req)
		if matchIndxes != nil {
			stringsReq := strings.Split(req[matchIndxes[1]:], ",")
			intsReq := make([]int, len(stringsReq))
			isInSequence := false

			for i, s := range stringsReq {
				intsReq[i], _ = strconv.Atoi(s)
			}
			for _, s := range intsReq {
				if num == s {
					isInSequence = true
					break
				}
			}
			if !isInSequence {
				isValidNum = false
				addErrMsg(&resultErrors, "value should be in: "+strings.Join(stringsReq, ","))
			}
			continue
		}
	}
	if !isValidNum {
		return false, errors.New(resultErrors.String())
	}
	return true, nil
}

func validSlice(items reflect.Value, requirements []string) (bool, error) {
	isValidSlice := true
	var err error
	var resultErrors strings.Builder

	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		switch item.Kind() {
		case reflect.String:
			_, err = validString(item.String(), requirements)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_, err = validInt(int(item.Int()), requirements)
		case reflect.Uint8:
			_, err = validInt(int(item.Uint()), requirements)
		default:
			fmt.Println("Type not defined <-") //TODO: remove after tests
		}
		if err != nil {
			isValidSlice = false
			addErrMsg(&resultErrors, err.Error())
		}
	}

	if !isValidSlice {
		return false, errors.New(resultErrors.String())
	} else {
		return true, nil
	}
}

func addErrMsg(erros *strings.Builder, errText string) {
	if erros.Len() > 0 {
		erros.WriteString("; ")
	}
	erros.WriteString(errText)
}
