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
	fromFile, fromSize, err := fileOpen(fromPath, offset)
	if err != nil {
		return fmt.Errorf("fromFileOpen: %w", err)
	}
	defer func() {
		if fromFile == nil {
			return
		}
		if err := fromFile.Close(); err != nil {
			panic(err)
		}
	}()

	toFile, err := fileCreate(toPath)
	if err != nil {
		return fmt.Errorf("toFileCreate: %w", err)
	}
	defer func() {
		if toFile == nil {
			return
		}
		if err := toFile.Close(); err != nil {
			panic(err)
		}
	}()

	copySize := fromSize - offset
	if limit > 0 && limit < copySize {
		copySize = limit
	}

	if offset > 0 {
		_, err = fromFile.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	err = copyRW(fromFile, toFile, copySize)
	if err != nil {
		return err
	}

	return nil
}

func fileOpen(fromPath string, offset int64) (*os.File, int64, error) {
	file, err := os.Open(fromPath)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		// закрыть файл, если ошибка
		if err != nil {
			if errCl := file.Close(); errCl != nil {
				panic(errCl)
			}
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	if !fileInfo.Mode().IsRegular() {
		err = ErrUnsupportedFile
		return nil, 0, err
	}

	size := fileInfo.Size()
	if offset >= size {
		err = ErrOffsetExceedsFileSize
		return nil, 0, err
	}

	return file, size, nil
}

func fileCreate(toPath string) (*os.File, error) {
	file, err := os.Create(toPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		// закрыть файл, если ошибка
		if err != nil {
			errCl := file.Close()
			if errCl != nil {
				panic(errCl)
			}
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}

	return file, nil
}

func copyRW(from io.Reader, to io.Writer, size int64) error {
	var (
		copySize int64 = 512 // кол-во байт копируемых за одну итерацию. из io.ReadAll
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
