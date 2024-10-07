package main

import (
	"fmt"

	reverse "golang.org/x/example/hello/reverse"
)

const HelloStr string = "Hello, OTUS!"

func main() {
	fmt.Println(reverse.String(HelloStr))
}
