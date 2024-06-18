package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	newFile, err := os.CreateTemp("./testdata/", "test_file_*.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Remove(newFile.Name())
	}()

	size, err := newFile.WriteString("Необходимо реализовать утилиту копирования файлов\n") // size = 95
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(size)

	t.Run("no error full file", func(t *testing.T) {
		toName := "./testdata/out_test_file.txt"
		err := Copy(newFile.Name(), toName, 0, 0)

		require.NoError(t, err)
		require.FileExistsf(t, toName, "file must %v exists", toName)

		toFile, _ := os.Open(toName)
		toFileInfo, _ := toFile.Stat()
		require.Equalf(t, int64(size), toFileInfo.Size(), "file size must be %v", size)

		defer os.Remove(toFile.Name())
	})

	t.Run("no error size less limit", func(t *testing.T) {
		toName := "./testdata/out_test_file.txt"
		err := Copy(newFile.Name(), toName, 90, 30)

		require.NoError(t, err)
		require.FileExistsf(t, toName, "file must %v exists", toName)

		toFile, _ := os.Open(toName)
		toFileInfo, _ := toFile.Stat()
		require.Equalf(t, int64(5), toFileInfo.Size(), "file size must be %v", 5)

		defer os.Remove(toFile.Name())
	})
}

func TestFileCreate(t *testing.T) {
	t.Run("not regular file", func(t *testing.T) {
		fromPath := "/dev/urandom"

		toFile, tstErr := fileCreate(fromPath)
		defer func() {
			if toFile == nil {
				return
			}
			errCl := toFile.Close()
			if errCl != nil {
				panic(errCl)
			}
		}()

		require.EqualError(t, tstErr, ErrUnsupportedFile.Error(), "actual err - %v", tstErr)
		require.Nil(t, toFile, "file must be nil")
	})

	t.Run("no error new file", func(t *testing.T) {
		toPath := "./testdata/out_test_create.txt"

		toFile, tstErr := fileCreate(toPath)
		defer func() {
			if toFile == nil {
				return
			}

			if errCl := toFile.Close(); errCl != nil {
				panic(errCl)
			}
			if errRm := os.Remove(toFile.Name()); errRm != nil {
				panic(errRm)
			}
		}()

		require.NoError(t, tstErr)
		require.NotNil(t, toFile, "file must not be nil")
	})

	t.Run("no error exist file", func(t *testing.T) {
		newFile, err := os.CreateTemp("./testdata/", "out_test_create_*.txt")
		if err != nil {
			panic(err)
		}
		defer func() {
			os.Remove(newFile.Name())
		}()

		toPath := newFile.Name()
		toFile, tstErr := fileCreate(toPath)
		defer func() {
			if toFile == nil {
				return
			}
			errCl := toFile.Close()
			if errCl != nil {
				panic(errCl)
			}
		}()

		require.NoError(t, tstErr)
		require.NotNil(t, toFile, "file must not be nil")
	})
}

func TestFileOpen(t *testing.T) {
	t.Run("not regular file", func(t *testing.T) {
		var offset int64 = 1000

		fromPath := "/dev/urandom"

		fromFile, fromSize, tstErr := fileOpen(fromPath, offset)
		defer func() {
			if fromFile == nil {
				return
			}
			errCl := fromFile.Close()
			if errCl != nil {
				panic(errCl)
			}
		}()

		require.EqualError(t, tstErr, ErrUnsupportedFile.Error(), "actual err - %v", tstErr)
		require.Nil(t, fromFile, "file must be nil")
		require.Equal(t, fromSize, int64(0), "file size must be 0")
	})

	t.Run("big file", func(t *testing.T) {
		var offset int64 = 1000
		fromPath := "./testdata/out_offset6000_limit1000.txt"

		fromFile, fromSize, tstErr := fileOpen(fromPath, offset)
		defer func() {
			if fromFile == nil {
				return
			}
			errCl := fromFile.Close()
			if errCl != nil {
				panic(errCl)
			}
		}()

		require.EqualError(t, tstErr, ErrOffsetExceedsFileSize.Error(), "actual err - %v", tstErr)
		require.Nil(t, fromFile, "file must be nil")
		require.Equal(t, fromSize, int64(0), "file size must be 0")
	})

	t.Run("no error", func(t *testing.T) {
		var offset int64 = 100

		fromPath := "./testdata/out_offset6000_limit1000.txt"

		fromFile, fromSize, tstErr := fileOpen(fromPath, offset)
		log.Println(fromSize)
		defer func() {
			if fromFile == nil {
				return
			}
			errCl := fromFile.Close()
			if errCl != nil {
				panic(errCl)
			}
		}()
		require.NoError(t, tstErr)
		require.NotNil(t, fromFile, "file must not be nil")
		require.Equalf(t, fromSize, int64(617), "file size must be %v", fromSize)
	})
}
