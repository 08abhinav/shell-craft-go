package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"slices"
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
		
		if command == "" || command == " "{
			continue
		}else if command == "exit"{
			break
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

		}else if command == "pwd" || strings.HasPrefix(command, "cd "){
			output, isBuiltin, err := utils.BuiltinExecute(command)
			if isBuiltin {
				if err != nil {
					args := strings.Fields(command)
					target := ""
					if len(args) > 1 {
						target = args[1]
					}
					fmt.Printf("cd: %s: No such file or directory\n", target)
				} else if output != "" {
					fmt.Println(output)
				}
			}
		}else if strings.HasPrefix(command, "echo"){
			args := utils.QuotingOps(command)

			if slices.Contains(args, ">") || slices.Contains(args, ">>"){
				if err := utils.Redirecting(args); err != nil{
					fmt.Fprintln(os.Stderr, "Error: ", err)
				}
			
			}else if len(args) > 1{
				fmt.Println(strings.Join(args[1: ], " "))
			}else{
				fmt.Println()
			}
		}else if strings.HasPrefix(command, "cat"){
			args := utils.QuotingOps(command)

			data, err := utils.RedirectinStdErr(args)
			if err != nil{
				fmt.Fprintln(os.Stderr, "Error: ", err)
			}
			fmt.Println(data)
		}else{
			output, err := utils.RunProgram(command)

			if err != nil{
				fmt.Fprint(os.Stderr, "Error reading input: ", err)
				os.Exit(1)
			}

			fmt.Print(output)
		}
	}
}
