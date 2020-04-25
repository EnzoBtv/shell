package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	cmdUtils "github.com/EnzoBtv/shell/utils"
)

func scanGitFolders(folders *[]string, folder string) error {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		return errors.New("It wasn't possible to open the specified folder")
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return errors.New("It wasn't possible to read the specified folder")
	}

	defer f.Close()

	path := ""

	for _, file := range files {
		if file.Mode().IsDir() {
			path = fmt.Sprintf("%s/%s", folder, file.Name())
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				*folders = append((*folders), path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" || file.Name() == "dist" {
				continue
			}

			scanGitFolders(folders, path)
		}
	}
	return nil
}

func parseFileLinesToSlice(filePath string) ([]string, error) {
	f, err := cmdUtils.OpenFile(filePath)

	if err != nil {
		return []string{}, err
	}

	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return []string{}, errors.New("There was an error reading dotfile")
	}

	return lines, nil
}

func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	ioutil.WriteFile(filePath, []byte(content), 0755)
}

func addSliceElementsToFile(filePath string, newRepos []string) (bool, error) {
	existingRepos, err := parseFileLinesToSlice(filePath)
	if err != nil {
		return false, err
	}
	cmdUtils.JoinSlices(newRepos, &existingRepos)
	dumpStringsSliceToFile(existingRepos, filePath)
	return true, nil
}

// Gitcontrib gets all your contributions of an specified folder and put into a Graph
func Gitcontrib(args []string) (bool, error) {
	if len(args) > 1 && args[1] == "--add" {
		repos := make([]string, 0)

		err := scanGitFolders(&repos, args[2])
		if err != nil {
			return false, err
		}

		filePath, err := cmdUtils.GetDotFilePath("gitcontrib")
		if err != nil {
			return false, err
		}

		ok, err := addSliceElementsToFile(filePath, repos)
		if !ok {
			return false, err
		}
	}
	return true, nil
}
