package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	cmdUtils "github.com/EnzoBtv/shell/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type column []int

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
		fmt.Println(err)
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

func fillCommits(email string, path string, commits *map[int]int) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return errors.New("It wasn't possible to plain open the Repository on path" + path)
	}

	ref, err := repo.Head()
	if err != nil {
		return errors.New("It wasn't possible to get the Head of the Repository on path" + path)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return errors.New("It wasn't possible to get the Logs of the Repository on path" + path)
	}

	offset := cmdUtils.CalcOffset()

	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := cmdUtils.CountDaysSinceDate(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}

		if daysAgo != cmdUtils.OutOfRange {
			(*commits)[daysAgo]++
		}
		return nil
	})

	if err != nil {
		return errors.New("Couldn't iterate over all commits")
	}
	return nil
}

func processRepositories(email string) (map[int]int, error) {
	filePath, err := cmdUtils.GetDotFilePath("gitcontrib")
	if err != nil {
		return map[int]int{}, err
	}

	repos, err := parseFileLinesToSlice(filePath)
	if err != nil {
		return map[int]int{}, err
	}

	commits := make(map[int]int, cmdUtils.CommitDays)
	for i := cmdUtils.CommitDays; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		fillCommits(email, path, &commits)
	}

	return commits, nil
}

func sortMapIntoSlice(newMap map[int]int) []int {
	keys := make([]int, 0)
	for key := range newMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	return keys
}

func buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7)
		dayinweek := k % 7

		if dayinweek == 0 {
			col = column{}
		}

		col = append(col, commits[k])

		if dayinweek == 6 {
			cols[week] = col
		}
	}

	return cols
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Printf(out)
}

func printMonths() {
	week := cmdUtils.GetBeginningOfDay(time.Now()).Add(-(cmdUtils.CommitDays * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := cmdUtils.CommitWeeks + 1; i >= 0; i-- {
			if i == cmdUtils.CommitWeeks+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				if i == 0 && j == cmdUtils.CalcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func printCommitsStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}

func gitAddFolder(args []string) (bool, error) {
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
	return true, nil
}

// Gitcontrib gets all your contributions of an specified folder and put into a Graph
func Gitcontrib(args []string) (bool, error) {
	if len(args) >= 3 {
		if args[1] == "--add" {
			return gitAddFolder(args)
		}
		if args[1] == "--email" {
			commits, err := processRepositories(args[2])
			if err != nil {
				return false, err
			}
			printCommitsStats(commits)
		}
	}
	return true, nil
}
