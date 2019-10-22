package main

import (
	"fmt"
	"os"
)

func main() {
	arg := os.Args[1:]
	fmt.Println("Argument u used is", arg[0])

}
