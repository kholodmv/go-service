package main

import (
	"fmt"
	"os"
)

func main() {
	// we formulate expectations: the analyzer must find an error,
	// described in the comment want
	i := 0
	fmt.Println(i)
	os.Exit(0) // want "not allowed using of os.Exit()"
}
