package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	defOutName   string = "./testdata/out_test_file.txt"
	defInputName string = "./testdata/input.txt"
)

func prepareTestFile() (string, int64, error) {
	testFile, err := os.CreateTemp("./testdata/", "test_file_*.txt")
	if err != nil {
		log.Fatalf("prepareTestFile: CreateTemp: %s", err)
	}
	defer func() {
		testFile.Close()
	}()

	size, err := testFile.WriteString("Необходимо реализовать утилиту копирования файлов\n") // size = 95
	if err != nil {
		log.Printf("prepareTestFile: %s", err)
		return "", 0, err
	}

	return testFile.Name(), int64(size), nil
}

func TestCopyFromFile(t *testing.T) {
	// исходный файл не корректный
	t.Run("fromFile not regular file", func(t *testing.T) {
		fromName := "/dev/urandom"
		toName := defOutName

		tstErr := Copy(fromName, toName, 0, 0)

		require.EqualError(t, tstErr, "fromFile: "+ErrUnsupportedFile.Error(), "actual err - %v", tstErr)
	})

	// исходный файл дирректория
	t.Run("fromFile is directory", func(t *testing.T) {
		fromName := "./testdata/"
		toName := defOutName

		tstErr := Copy(fromName, toName, 0, 0)

		require.EqualError(t, tstErr, "fromFile: "+ErrUnsupportedFile.Error(), "actual err - %v", tstErr)
	})

	// исходный файл не существует
	t.Run("fromFile no such file", func(t *testing.T) {
		fromName := "./testdata/not_exists_file.txt"
		toName := defOutName

		tstErr := Copy(fromName, toName, 0, 0)

		require.EqualError(t, tstErr,
			fmt.Sprintf("fromFileOpen: open %v: no such file or directory", fromName),
			"actual err - %v", tstErr)
	})
}

func TestCopyToFile(t *testing.T) {
	// целевой файл не корректный
	t.Run("toFile not regular file", func(t *testing.T) {
		tstErr := Copy(defInputName, "/dev/urandom", 0, 0)

		require.EqualError(t, tstErr, "toFile: "+ErrUnsupportedFile.Error(), "actual err - %v", tstErr)
	})

	// целевой файл директория
	t.Run("toFile is a directory", func(t *testing.T) {
		toName := "./testdata/"

		tstErr := Copy(defInputName, toName, 0, 0)

		require.EqualError(t, tstErr, fmt.Sprintf("toFileCreate: open %v: is a directory", toName), "actual err - %v", tstErr)
	})
}

func TestCopy(t *testing.T) {
	testFileName, size, err := prepareTestFile()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		os.Remove(testFileName)
	}()

	// исходный файл меньше offset
	t.Run("fromFile size less offset", func(t *testing.T) {
		tstErr := Copy(testFileName, defOutName, size, 0)

		require.EqualError(t, tstErr, "fromFile: offset exceeds file size", "actual err - %v", tstErr)
	})

	// без ошибок копируем в новый файл
	t.Run("no error full copy new file", func(t *testing.T) {
		tstErr := Copy(testFileName, defOutName, 0, 0)
		require.NoError(t, tstErr)

		toFile, tstErr := os.Open(defOutName)
		if tstErr != nil {
			t.Error(tstErr)
			return
		}

		defer func() {
			toFile.Close()
			os.Remove(toFile.Name())
		}()

		toFileInfo, tstErr := toFile.Stat()
		if tstErr != nil {
			t.Error(tstErr)
			return
		}

		require.Equalf(t, size, toFileInfo.Size(), "file size must be %v", size)

		require.NotNil(t, toFile, "file must not be nil")
	})

	// без ошибок копируем файл целиком
	t.Run("no error full file", func(t *testing.T) {
		err := Copy(testFileName, defOutName, 0, 0)

		require.NoError(t, err)
		require.FileExistsf(t, defOutName, "file must %v exists", defOutName)

		toFile, tstErr := os.Open(defOutName)
		if tstErr != nil {
			t.Error(tstErr)
			return
		}

		defer func() {
			toFile.Close()
			os.Remove(toFile.Name())
		}()

		toFileInfo, tstErr := toFile.Stat()
		if tstErr != nil {
			t.Error(tstErr)
			return
		}

		require.Equalf(t, size, toFileInfo.Size(), "file size must be %v", size)
	})

	// без ошибок копируем, но limit превышает размер
	t.Run("no error size less limit", func(t *testing.T) {
		err := Copy(testFileName, defOutName, 90, 30)

		require.NoError(t, err)
		require.FileExistsf(t, defOutName, "file must %v exists", defOutName)

		toFile, tstErr := os.Open(defOutName)
		if tstErr != nil {
			t.Error(tstErr)
			return
		}

		defer func() {
			toFile.Close()
			os.Remove(toFile.Name())
		}()

		toFileInfo, tstErr := toFile.Stat()
		if tstErr != nil {
			t.Error(tstErr)
			return
		}
		require.Equalf(t, int64(5), toFileInfo.Size(), "file size must be %v", 5)
	})

	// без ошибок копируем в существующий файл
	t.Run("no error exists toFile", func(t *testing.T) {
		var limit int64 = 10
		fromName := "./testdata/input.txt"
		toName := testFileName
		err := Copy(fromName, toName, 0, limit)

		require.NoError(t, err)
		require.FileExistsf(t, toName, "file must %v exists", toName)

		toFile, tstErr := os.Open(toName)
		if tstErr != nil {
			t.Error(err)
			return
		}

		defer func() {
			toFile.Close()
			os.Remove(toFile.Name())
		}()

		toFileInfo, tstErr := toFile.Stat()
		if tstErr != nil {
			t.Error(err)
			return
		}

		require.Equalf(t, limit, toFileInfo.Size(), "file size must be %v", size)
	})
}
