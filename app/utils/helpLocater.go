package utils

import (
	"fmt"
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
var ErrNoRedirect = errors.New("could not redirect")
var ErrInvalidRedirect = errors.New("invalid command")

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
			continue
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
			continue
		}

		switch char {
		case '\'':
			inSingleQuotes = true
			hasToken = true

		case '"':
			inDoubleQuotes = true
			hasToken = true

		case ' ', '\t':
			if hasToken {
				args = append(args, current.String())
				current.Reset()
				hasToken = false
			}

		case '>':
			if hasToken {
				args = append(args, current.String())
				current.Reset()
				hasToken = false
			}

			if i+1 < len(input) && input[i+1] == '>' {
				args = append(args, ">>")
				i++
			} else {
				args = append(args, ">")
			}

		default:
			current.WriteByte(char)
			hasToken = true
		}
	}

	if hasToken{
		args = append(args, current.String())
	}
	return args
}

type Command struct{
	Name 		string
	Args 	  []string
	Redirect 	string
	FileName 	string
}

func Redirecting(tokens []string) error{
	redirectIdx := -1

	for i, token := range tokens{
		if token == ">" || token == ">>"{
			redirectIdx = i
			break
		}
	}

	if redirectIdx == -1{
		return ErrNoRedirect
	}

	if redirectIdx+1 >= len(tokens) {
		return ErrInvalidRedirect
	}

	cmd := &Command{
		Name: tokens[0],
		Args: tokens[1 : redirectIdx],
		Redirect: tokens[redirectIdx],
		FileName: tokens[redirectIdx + 1],
	}

	var file *os.File
	var err error

	if cmd.Redirect == ">"{
		file, err = os.OpenFile(
			cmd.FileName, 
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 
			0644,
		)
	}else if cmd.Redirect == ">>"{
		file, err = os.OpenFile(
			cmd.FileName, 
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 
			0644,
		)
	}
	if err != nil{
		return err
	}
	defer file.Close()

	process := exec.Command(cmd.Name, cmd.Args...)
	process.Stdout = file
	process.Stderr = os.Stderr

	return process.Run()
}