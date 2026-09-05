package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	builtinCommands := []string{"echo", "type", "exit"}
	for{

		fmt.Print("$ ")
	
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil{
			fmt.Fprint(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}

		command = strings.TrimSpace(command)
		
		if command == "exit"{
			break
		}else if strings.HasPrefix(command, "echo "){
			fmt.Println(command[5: ])
		}else if strings.HasPrefix(command, "type "){
			found := false;
			for _, i := range builtinCommands{
				if command[5: ] == i{
					fmt.Println(command[5: ], "is a shell builtin")
					found = true
					break
				}
			}
			if !found {
				fmt.Println( command[5: ], ": not found")
			}
		}else{
			fmt.Println(command + ": command not found")
		}
	}
}
