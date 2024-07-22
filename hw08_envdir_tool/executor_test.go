package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("ls no error no args", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls"}, Environment{})

		require.Equal(t, exitCode, 0, "exitCode incorrect")
	})

	t.Run("unknow command", func(t *testing.T) {
		exitCode := RunCmd([]string{"lssssss"}, Environment{})

		require.Equal(t, exitCode, -1, "exitCode incorrect")
	})

	t.Run("ls correct argument", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls", "-l"}, Environment{})

		require.Equal(t, exitCode, 0, "exitCode incorrect")
	})

	t.Run("ls unknow argument", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls", "-y"}, Environment{})

		require.Equal(t, exitCode, 2, "exitCode incorrect")
	})
}
