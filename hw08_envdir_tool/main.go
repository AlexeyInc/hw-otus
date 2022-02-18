package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Args == nil || len(os.Args) < 1 {
		fmt.Println("Haven't found any args")
		os.Exit(0)
	}
	envVariables, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	res := RunCmd(os.Args[2:], envVariables)

	os.Exit(res)
}
