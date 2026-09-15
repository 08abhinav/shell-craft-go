package utils

import (
	"errors"
	"io"
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
		case '2':
			if i+1 < len(input) && input[i+1] == '>'{
				if current.Len() > 0{
					current.Reset()
				}
				args = append(args, "2>")
				i++
			}else{
				current.WriteByte(char)
				hasToken = true
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
		if token == ">" || token == ">>" || token == "2>" || token == "2>>"{
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

	if cmd.Redirect == ">" || cmd.Redirect == "2>"{
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

type RedirectStdErr struct{
	Name			string
	File1			string
	RedirectToken	string
	RedirectFile	string	
}

func Concatenate(tokens []string) (string, error){
	redirectIdx := -1
	var file *os.File
	var err error

	for i, token := range tokens{
		if token == "2>" || token == "2>>"{
			redirectIdx = i
			break
		}
	}

	if redirectIdx == -1{
		file, err = os.Open(tokens[1])
		if err != nil{
			return "", err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil{
			return "", err
		}
		return string(data), err
	}

	if redirectIdx+1 >= len(tokens) {
		return "", ErrInvalidRedirect
	}

	cmd := &RedirectStdErr{
		Name:			tokens[0],
		File1:			tokens[1],
		RedirectToken: 	tokens[redirectIdx],
		RedirectFile:	tokens[redirectIdx + 1],
	}

	if cmd.RedirectToken == "2>" {
		file, err = os.OpenFile(
			cmd.RedirectFile,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0644,
		)
	} else {
		file, err = os.OpenFile(
			cmd.RedirectFile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
	}

	if err != nil {
		return "", err
	}
	defer file.Close()

	process := exec.Command(cmd.Name, cmd.File1)
	process.Stderr = file

	output, err := process.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}