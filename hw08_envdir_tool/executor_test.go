package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testEnvName := "TESTENV"

	t.Run("'NeedRemove' env var was being removed when file empty", func(t *testing.T) {
		env := make(Environment)
		env[testEnvName] = EnvValue{"", true}

		os.Setenv(testEnvName, env[testEnvName].Value)

		_, envExistsBefore := os.LookupEnv(testEnvName)

		RunCmd(nil, env)

		_, envExistsAfter := os.LookupEnv(testEnvName)

		require.True(t, envExistsBefore)
		require.False(t, envExistsAfter)
	})

	t.Run("'NeedRemove' env var with new value was being updated", func(t *testing.T) {
		env := make(Environment)

		oldEnvValue := "oldValue"
		env[testEnvName] = EnvValue{oldEnvValue, true}

		os.Setenv(testEnvName, env[testEnvName].Value)

		updatedEnvValue := "newValue"
		if envVar, ok := env[testEnvName]; ok {
			envVar.Value = updatedEnvValue

			env[testEnvName] = envVar
		}

		RunCmd(nil, env)

		curEnvValue := os.Getenv(testEnvName)

		require.True(t, curEnvValue == updatedEnvValue)
	})
}
