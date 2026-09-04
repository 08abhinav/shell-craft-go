package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	var input string
	fmt.Print("$ ")
	fmt.Scanln(&input)
	fmt.Printf("%v: command not found", input)
}
