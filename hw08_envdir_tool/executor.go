package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {

	for k, ev := range env {
		if env[k].NeedRemove {
			os.Unsetenv(k)
		}
		if ev.Value != "" {
			os.Setenv(k, ev.Value)
		}
	}

	resCmd := exec.Command(cmd[0], cmd[1])
	resCmd.Env = os.Environ()
	resCmd.Env = append(resCmd.Env, "ADDED=from original env")
	resCmd.Args = append(resCmd.Args, cmd[2], cmd[3])
	resCmd.Stdout = os.Stdout

	if err := resCmd.Run(); err != nil {
		return 0
	}

	return 1
}
