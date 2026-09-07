package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
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
			token := command[5:]

			path, isBuiltin, err := utils.FindCommand(token)
			if isBuiltin {
				fmt.Printf("%s is a shell builtin\n", token)
			} else if err == nil {
				fmt.Printf("%s is %s\n", token, path)
			} else {
				fmt.Printf("%s: not found\n", token)
			}

		}else{
			fmt.Println(command + ": command not found")
		}
	}
}
