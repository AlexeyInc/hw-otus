package main

import (
	"fmt"
	"os"
)

func main() {
	//
	// "./testdata/env"
	envVariables, err := ReadDir(os.Args[1])

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	//
	// []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
	res := RunCmd(os.Args[2:], envVariables)

	os.Exit(res)
}
