package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 0
	}
	if err := SetEnv(env); err != nil {
		log.Println("SetEnv:", err)
		return 0x7F
	}
	cmdName := cmd[0]

	cmdArgs := []string{}
	if len(cmd) > 1 {
		cmdArgs = cmd[1:]
	}

	cmdExec := exec.Command(cmdName, cmdArgs...)

	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		log.Println("error Run:", err)
	}

	return cmdExec.ProcessState.ExitCode()
}
