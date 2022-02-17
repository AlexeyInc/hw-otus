package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const newLineSymb byte = byte('\n')

var nullSymb []byte = []byte{0}

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {

	sectionItems, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make(Environment)

	for _, item := range sectionItems {
		if !item.IsDir() {
			needRemove := false

			fileName := validateFileName(item.Name())

			if item.Size() == 0 {
				needRemove = true
				envVars[fileName] = EnvValue{"", needRemove}
				continue
			}

			_, exist := os.LookupEnv(fileName)
			if exist {
				needRemove = true
			}

			filePath := filepath.Join(dir, fileName)
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			fileDataStr := validateFileData(fileData)

			envVars[fileName] = EnvValue{fileDataStr, needRemove}
		}
	}

	return envVars, nil
}

func validateFileName(name string) string {
	name = strings.ReplaceAll(name, "=", "")
	return name
}

func validateFileData(data []byte) string {
	newLineIndx := bytes.IndexByte(data, 10)
	if newLineIndx > 0 {
		data = data[:newLineIndx]
	}
	data = bytes.ReplaceAll(data, nullSymb, []byte{newLineSymb})
	resultStr := strings.TrimRight(string(data), "	 ")
	return resultStr
}
