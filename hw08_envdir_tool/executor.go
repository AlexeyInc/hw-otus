package main

import (
	"os"
	"os/exec"
)

var envVarsToAdd = []string{"ADDED=from original env"}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, ev := range env {
		if env[k].NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return 0
			}
		}
		if ev.Value != "" {
			if err := os.Setenv(k, ev.Value); err != nil {
				return 0
			}
		}
	}

	if len(cmd) >= 2 {
		resCmd := exec.Command(cmd[0], cmd[1]) //nolint
		resCmd.Env = os.Environ()
		resCmd.Env = append(resCmd.Env, envVarsToAdd...)
		if len(cmd) >= 4 {
			resCmd.Args = append(resCmd.Args, cmd[2], cmd[3])
		}
		resCmd.Stdout = os.Stdout

		if err := resCmd.Run(); err != nil {
			return 0
		}
	}
	return 1
}
