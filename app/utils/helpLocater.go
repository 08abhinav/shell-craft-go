package utils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)	


var BuiltInCommands = []string{
	"echo", 
	"type",
	"exit",
}

var ErrNotFound = errors.New("command not found")

func FindCommand(target string) (string, bool, error) {
	for _, builtIn := range BuiltInCommands{
		if target == builtIn{
			return "", true, nil
		}
	}

	path, err := exec.LookPath(target)
	if err == nil{
		return path, false, nil
	}
	
	return "", false, ErrNotFound
}

func RunProgram(input string) (string, error){
	value := strings.Fields(input)

	command := value[0]
	args := value[1: ]
	
	cmd := exec.Command(command, args...)

	output, err := cmd.Output()
	if err == nil{
		return string(output), nil
	}

	return "", ErrNotFound
}

func BuiltinExecute(cmd string) (string, bool, error) {
	value := strings.Fields(cmd)
	if len(value) == 0 {
		return "", false, nil
	}

	command := value[0]
	args := value[1:]

	switch command {
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return "", true, err
		}
		return dir, true, nil

	case "cd":
		if len(args) != 1{
			return "", true, ErrNotFound
		}

		arg := args[0]

		if arg == "~"{
			err := os.Chdir(os.Getenv("HOME"))
			if err != nil {
				return "", true, err
			}

			return "", true, nil
		}

		err := os.Chdir(args[0])
		if err != nil {
			return "", true, err
		}
		return "", true, nil

	default:
		return "", false, nil
	}
}

func QuotingOps(input string) []string{
	var args []string
	var current strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	hasToken := false

	input = strings.ReplaceAll(input, "\u00a0", " ")

	for i := 0; i < len(input); i++{
		char := input[i]
		
		if inSingleQuotes{
			if char == '\''{
				inSingleQuotes = false
			}else{
				current.WriteByte(char)
			}
			hasToken = true
		}else if inDoubleQuotes{
			if char == '"'{
				inDoubleQuotes = false
			}else if char == '\\' && i+1 < len(input) && (input[i+1] == '"' || input[i+1] == '\\' || input[i+1] == '$' || input[i+1] == '\n') {
				i++
				current.WriteByte(input[i])
			} else {
				current.WriteByte(char)
			}
			hasToken = true
		}else {
			if char == '\'' {
				inSingleQuotes = true
				hasToken = true
			} else if char == '"' {
				inDoubleQuotes = true
				hasToken = true
			} else if char == ' ' || char == '\t' {

				if hasToken {
					args = append(args, current.String())
					current.Reset()
					hasToken = false
				}
			} else {
				current.WriteByte(char)
				hasToken = true
			}
		}
	}

	if hasToken{
		args = append(args, current.String())
	}
	return args
}