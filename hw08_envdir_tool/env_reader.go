package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dir = filepath.Clean(dir)

	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("readDir: %w", err)
	}

	env := make(Environment)

	for _, de := range des {
		info, err := de.Info()
		if err != nil {
			log.Printf("info file '%s': %s", de.Name(), err)
			continue
		}
		if !info.Mode().IsRegular() {
			log.Printf("file '%s' is not regular", de.Name())
			continue
		}

		if info.Size() == 0 {
			env[de.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		fName := filepath.Join(dir, de.Name())
		data, err := readData(fName)
		if err != nil {
			log.Printf("ReadFile: %s", err)
			continue
		}

		data = bytes.ReplaceAll(data, []byte("\000"), []byte("\n"))

		env[de.Name()] = EnvValue{Value: strings.TrimRight(string(data), " \t"), NeedRemove: false}
	}

	return env, nil
}

func readData(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open file '%s': %w", fileName, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		return make([]byte, 0), nil
	}

	return scanner.Bytes(), nil
}

func SetEnv(env Environment) error {
	for envName, envVal := range env {
		if envVal.NeedRemove {
			if err := os.Unsetenv(envName); err != nil {
				return err
			}
		} else {
			if err := os.Setenv(envName, envVal.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
