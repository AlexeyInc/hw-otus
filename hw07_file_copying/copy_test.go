package main

import (
	"bytes"
	"io"
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
		require.True(t, compareFiles(realInputFile, tempOutFile))
	})

	limits := []struct {
		value       int64
		resFilePath string
	}{
		{10, "testdata/out_offset0_limit10.txt"},
		{1000, "testdata/out_offset0_limit1000.txt"},
		{10000, "testdata/out_offset0_limit10000.txt"},
		{1025, "testdata/out_offset0_limit1025.txt"},
	}

	for _, limit := range limits {
		t.Run("Check behaviour with 'limit' value", func(t *testing.T) {
			defer helperRemoveFile(t, tempOutFile)

			Copy(realInputFile, tempOutFile, 0, limit.value)

			fInputSize := getFileSize(realInputFile)
			if limit.value > fInputSize {
				limit.value = fInputSize
			}

			fOutSize := getFileSize(tempOutFile)

			require.Equal(t, limit.value, fOutSize)
			require.True(t, compareFiles(tempOutFile, limit.resFilePath))
		})
	}

	offsetAndLimits := []struct {
		offset      int64
		limit       int64
		resFilePath string
	}{
		{100, 1000, "testdata/out_offset100_limit1000.txt"},
		{6000, 10000, "testdata/out_offset6000_limit1000.txt"},
		{-1, 1000, "testdata/out_offset0_limit1000.txt"},
		{100, -1, "testdata/out_offset100_limit0.txt"},
		{-100, -100, "testdata/out_offset0_limit0.txt"},
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
			require.True(t, compareFiles(tempOutFile, ol.resFilePath))
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

func compareFiles(file1, file2 string) bool {
	chunkSize := 1024

	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
