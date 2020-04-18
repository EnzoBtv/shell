package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	cmdUtils "github.com/EnzoBtv/shell/utils"
)

// Ls shows all the files and folders of a directory
func Ls() (bool, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return false, errors.New(err.Error())
	}
	files, err := ioutil.ReadDir(currentDir)

	if err != nil {
		return false, errors.New(err.Error())
	}

	newFiles := cmdUtils.MapArray(files, func(item interface{}, i int) interface{} {
		file, ok := item.(os.FileInfo)

		if !ok {
			return errors.New(err.Error())
		}

		return file.Name()
	})

	returnString := ""

	for _, file := range newFiles {
		returnString += fmt.Sprintf("%v\n", file)
	}

	fmt.Println(returnString)
	return true, nil
}
