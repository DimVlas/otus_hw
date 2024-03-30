package main

import (
	"fmt"

	stringutil "golang.org/x/example/stringutil"
)

const (
	HelloStr string = "Hello, OTUS!"
)

func main() {
	fmt.Println(stringutil.Reverse(HelloStr))
}
