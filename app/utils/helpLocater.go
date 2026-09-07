package utils

import (
	"errors"
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