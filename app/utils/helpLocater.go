package utils

import (
	"errors"
	"os"
	"path/filepath"
)	


var BuiltInCommands = []string{
	"echo", 
	"type",
	"exit",
}

var PATHS = []string{
	`C:\Program Files\Git\usr\bin`,
}
var ErrNotFound = errors.New("command not found")

func FindCommand(target string) (string, bool, error) {
	for _, builtIn := range BuiltInCommands{
		if target == builtIn{
			return "", true, nil
		}
	}

	target = target + ".exe"

	for _, dir := range PATHS{
		fullPath := filepath.Join(dir, target)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir(){
			return fullPath, false, nil
		}
	}

	return "", false, ErrNotFound
}