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

		err := os.Chdir(args[0])
		if err != nil {
			return "", true, err
		}
		return "", true, nil

	default:
		return "", false, nil
	}
}