package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <path_to_env_dir> <command> [arg1 arg2 ...]\n", filepath.Base(os.Args[0]))
		return
	}

	envDir := os.Args[1]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("ReadDir:", err)
	}

	retCode := RunCmd(os.Args[2:], env)

	os.Exit(retCode)
}
