package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testdataPath = "." + string(os.PathSeparator) + "testdata"

func TestReadDir(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		envDir := testdataPath + string(os.PathSeparator) + "env"
		testEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		fmt.Println("envDir =", envDir)
		env, err := ReadDir(envDir)

		require.NoError(t, err, "no error")
		require.Equal(t, env, testEnv, "Environment incorrect")
	})

	t.Run("no such directory", func(t *testing.T) {
		envDir := testdataPath + string(os.PathSeparator) + "env_nosuch"

		_, tstErr := ReadDir(envDir)

		require.EqualError(
			t,
			tstErr,
			"readDir: open testdata/env_nosuch: no such file or directory",
			"actual err - %v", tstErr)
	})

	t.Run("empty directory", func(t *testing.T) {
		envDir, err := os.MkdirTemp(testdataPath, "env_*")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(envDir)

		env, tstErr := ReadDir(envDir)

		require.Len(t, env, 0, "Env must be empty")
		require.NoError(t, tstErr, "no error")
	})
}

func TestSetEnv(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		errTest := SetEnv(Environment{})

		require.NoError(t, errTest, "no error")
	})

	t.Run("var remove", func(t *testing.T) {
		os.Setenv("var1", "Value1")
		os.Setenv("var2", "Value2")

		errTest := SetEnv(Environment{
			"var1": {"", true},
			"var2": {"value2", false},
			"var3": {"", false},
		})

		require.NoError(t, errTest, "no error")

		_, var1Exists := os.LookupEnv("var1")
		require.Equal(t, false, var1Exists, "var1 must be remove")

		var2Val, var2Exists := os.LookupEnv("var2")
		require.Equal(t, true, var2Exists, "var2 should be")
		require.Equal(t, "value2", var2Val, "var2 must be 'value2'")

		var3Val, var3Exists := os.LookupEnv("var3")
		require.Equal(t, true, var3Exists, "var3 should be")
		require.Equal(t, "", var3Val, "var3 must be empty")
	})

	t.Run("error variable name", func(t *testing.T) {
		errTest := SetEnv(Environment{
			"var=": {"varVal", false},
		})

		require.EqualError(t, errTest, "setenv: invalid argument", "expected error: 'setenv: invalid argument'")
	})

	t.Run("error variable value", func(t *testing.T) {
		varVal := "var" + string([]byte("\000")) + "Val"

		errTest := SetEnv(Environment{
			"var": {varVal, false},
		})

		require.EqualError(t, errTest, "setenv: invalid argument", "expected error: 'setenv: invalid argument'")
	})
}
