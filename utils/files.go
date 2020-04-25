package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// GetDotFilePath gets the homedir of the user and
// concatenates with a .file specified by the parameter
func GetDotFilePath(dotFileName string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("It wasn't possible to get the current user")
	}

	dotFileName = strings.TrimPrefix(dotFileName, ".")

	separator := "/"

	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	dotFile := fmt.Sprintf("%s%s.%s", usr.HomeDir, separator, dotFileName)
	return dotFile, nil
}

// OpenFile opens a file if it exists, else creates one
func OpenFile(filePath string) (*os.File, error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("It was not possible to access the dotfile")
		}
		_, err := os.Create(filePath)
		if err != nil {
			return nil, errors.New("It was not possible to create the file")
		}

	}

	return f, nil
}
