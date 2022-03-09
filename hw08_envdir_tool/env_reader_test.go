package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	pathToEnv := "./testdata/env/"
	envFilesCount := 5
	notValidFileName := "envTEMP="
	fileNameAsExistingEnvVar := "USER"
	emptyDirName := "tempDir"
	emptyFileName := "UNSET"

	t.Run("default behaviour", func(t *testing.T) {
		envVariables, err := ReadDir(pathToEnv)

		require.Nil(t, err)
		require.Equal(t, len(envVariables), envFilesCount)
	})

	t.Run("default behaviour if there are no files in directory", func(t *testing.T) {
		defer os.Remove(emptyDirName)

		if err := os.Mkdir(emptyDirName, os.ModePerm); err != nil {
			t.Error(err)
		}

		envVariables, err := ReadDir(emptyDirName)

		require.Nil(t, err)
		require.Equal(t, len(envVariables), 0)
	})

	t.Run("wrong path to env files", func(t *testing.T) {
		envVariables, err := ReadDir(pathToEnv + "/notExistingFolder")

		require.Nil(t, envVariables)
		require.True(t, strings.Contains(err.Error(), "no such file or directory"))
	})

	t.Run("check file name validation", func(t *testing.T) {
		defer os.Remove(pathToEnv + notValidFileName)

		if _, err := os.Create(pathToEnv + notValidFileName); err != nil {
			t.Error(err)
		}

		envVariables, _ := ReadDir(pathToEnv)

		_, invalidEnvNameExist := envVariables[notValidFileName]

		validFileName := strings.ReplaceAll(notValidFileName, "=", "")
		_, validEnvNameExist := envVariables[validFileName]

		require.False(t, invalidEnvNameExist)
		require.True(t, validEnvNameExist)
	})

	t.Run("empty file should be deleted", func(t *testing.T) {
		envVariables, _ := ReadDir(pathToEnv)

		envValue := envVariables[emptyFileName]

		require.True(t, envValue.NeedRemove)
	})

	t.Run("existing env variable should be deleted", func(t *testing.T) {
		defer os.Remove(pathToEnv + fileNameAsExistingEnvVar)

		if _, err := os.Create(pathToEnv + fileNameAsExistingEnvVar); err != nil {
			t.Error(err)
		}

		envVariables, _ := ReadDir(pathToEnv)

		envValue := envVariables[fileNameAsExistingEnvVar]

		require.True(t, envValue.NeedRemove)
	})
}
