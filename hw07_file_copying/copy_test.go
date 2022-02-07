package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyErrors(t *testing.T) {
	realInputFile := "testdata/input.txt"
	notExistInputFile := "testdata/dummyInput.txt"
	tempOutFile := "testdata/dummyOut.txt"

	t.Run("Should return err if 'input' file not exists", func(t *testing.T) {
		err := Copy(notExistInputFile, "", 0, 0)

		require.Error(t, err)
	})

	t.Run("Should NOT return err if 'out' file not exists", func(t *testing.T) {
		defer helperRemoveFile(t, tempOutFile)

		err := Copy(realInputFile, tempOutFile, 0, 0)

		require.Equal(t, nil, err)
	})

	t.Run("Should return err if 'offset' bigger than 'input' file size", func(t *testing.T) {
		defer helperRemoveFile(t, tempOutFile)

		err := Copy(realInputFile, tempOutFile, 10000, 0)

		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("Check default behaviour", func(t *testing.T) {
		defer helperRemoveFile(t, tempOutFile)

		Copy(realInputFile, tempOutFile, 0, 0)

		fInputSize := getFileSize(realInputFile)
		fOutSize := getFileSize(tempOutFile)

		require.Equal(t, fInputSize, fOutSize)
	})

	limits := []int64{10, 1000, 1025, 10000}

	for _, limit := range limits {
		t.Run("Check behaviour with 'limit' value", func(t *testing.T) {
			defer helperRemoveFile(t, tempOutFile)

			Copy(realInputFile, tempOutFile, 0, limit)

			fInputSize := getFileSize(realInputFile)
			if limit > fInputSize {
				limit = fInputSize
			}

			fOutSize := getFileSize(tempOutFile)

			require.Equal(t, limit, fOutSize)
		})
	}

	offsetAndLimits := []struct {
		offset int64
		limit  int64
	}{
		{100, 1000},
		{6000, 10000},
		{-1, 1000},
		{100, -1},
		{-100, -100},
	}

	for _, ol := range offsetAndLimits {
		t.Run("Check behaviour with 'offset' and 'limit' values", func(t *testing.T) {
			defer helperRemoveFile(t, tempOutFile)

			Copy(realInputFile, tempOutFile, ol.offset, ol.limit)

			fInputSize := getFileSize(realInputFile)

			expectedSize := fInputSize
			if ol.offset > 0 {
				expectedSize = fInputSize - ol.offset
			}
			if expectedSize > ol.limit {
				if ol.limit > 0 {
					expectedSize = ol.limit
				}
			}

			fOutSize := getFileSize(tempOutFile)

			require.Equal(t, expectedSize, fOutSize)
		})
	}
}

func helperRemoveFile(tb testing.TB, removeFilePath string) {
	tb.Helper()
	os.Remove(removeFilePath)
}

func getFileSize(fPath string) int64 {
	info, err := os.Stat(fPath)
	if err != nil {
		return 0
	}
	return info.Size()
}
