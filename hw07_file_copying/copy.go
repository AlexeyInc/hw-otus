package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fSrc, fiErr := os.OpenFile(fromPath, os.O_RDONLY, 0o444)
	if fiErr != nil {
		return fiErr
	}
	defer fSrc.Close()

	fDest, foErr := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if foErr != nil {
		return foErr
	}
	defer fDest.Close()

	fSrcInfo, err := fSrc.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if err := addOffset(fSrc, fSrcInfo, offset); err != nil {
		return err
	}

	copyErr := copyWithProgressBar(fSrc, fDest, fSrcInfo, limit)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func addOffset(file *os.File, fInfo os.FileInfo, offset int64) error {
	if offset > 0 {
		if offset > fInfo.Size() {
			return ErrOffsetExceedsFileSize
		}
		if _, err := file.Seek(offset, 0); err != nil {
			return err
		}
	}
	return nil
}

func copyWithProgressBar(fSrc, fDest *os.File, fSrcInfo fs.FileInfo, limitToCopy int64) error {
	var copyingCompleted bool
	var defBufChunk, barCount int64 = 1024, 0

	if limitToCopy > 0 && defBufChunk > limitToCopy {
		defBufChunk = limitToCopy
		copyingCompleted = true
		barCount = 1
	} else {
		fileSize := fSrcInfo.Size() - offset
		if limitToCopy > 0 && fileSize > limitToCopy {
			setBarCount(&barCount, limitToCopy, defBufChunk)
		} else {
			setBarCount(&barCount, fileSize, defBufChunk)
		}
	}

	buf := make([]byte, defBufChunk)
	bar := pb.New64(barCount).Start()
	bytesRead := 0

	for {
		n, err := fSrc.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		bytesRead += n

		if _, err := fDest.Write(buf[:n]); err != nil {
			return err
		}
		bar.Increment()
		time.Sleep(time.Millisecond * 100)
		if copyingCompleted {
			break
		}

		if limitToCopy > 0 && int64(bytesRead)+defBufChunk > limitToCopy {
			buf = make([]byte, limitToCopy-int64(bytesRead))
			copyingCompleted = true
		}
	}

	bar.Finish()

	return nil
}

func setBarCount(barCount *int64, fSize, buffer int64) {
	*barCount = fSize / buffer
	if fSize%buffer != 0 {
		*barCount++
	}
}
