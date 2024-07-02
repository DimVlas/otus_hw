package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// fromFile open and check
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("fromFileOpen: %w", err)
	}
	defer fromFile.Close()

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("fromFile: %w", err)
	}

	if !fromInfo.Mode().IsRegular() {
		return fmt.Errorf("fromFile: %w", ErrUnsupportedFile)
	}

	fromSize := fromInfo.Size()
	if offset >= fromSize {
		return fmt.Errorf("fromFile: %w", ErrOffsetExceedsFileSize)
	}

	// toFile open and check
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("toFileCreate: %w", err)
	}
	defer toFile.Close()

	toInfo, err := toFile.Stat()
	if err != nil {
		return fmt.Errorf("toFile: %w", err)
	}

	if !toInfo.Mode().IsRegular() {
		return fmt.Errorf("toFile: %w", ErrUnsupportedFile)
	}

	// prepare copy
	copySize := fromSize - offset
	if limit > 0 && limit < copySize {
		copySize = limit
	}

	if offset > 0 {
		_, err = fromFile.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("fromFile: %w", err)
		}
	}

	err = copyRW(fromFile, toFile, copySize)
	if err != nil {
		return fmt.Errorf("copyRW: %w", err)
	}

	return toFile.Sync()
}

// значение по умолчанию кол-ва байт копируемых за одну итерацию. из io.ReadAll.
const defCopySize int64 = 512

func copyRW(from io.Reader, to io.Writer, size int64) error {
	var (
		copySize = defCopySize // кол-во байт копируемых за одну итерацию.
		copied   int64
	)
	if size < copySize {
		copySize = size
	}

	// bar := pb.New64(size)

	for i := size; i > 0; i -= copySize {
		if i < copySize {
			copySize = i
		}

		c, err := io.CopyN(to, from, copySize)
		if err != nil {
			return err
		}
		copied += c
		fmt.Printf("\r%v%%", (100*copied)/size)
		// bar.Add64(copySize)
	}
	fmt.Println()
	// bar.Finish()

	return nil
}
