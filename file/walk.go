package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/SublimeIbanez/todor/common"
	"github.com/boyter/gocodewalker"
)

var buffer = make([]byte, 0, 64*1024)

const buffer_size int = 1024 * 1024

// Recursively walk through the directory and read through all items
func (parser *Parser) WalkDir(input_path string) error {
	_, err := os.Stat(input_path)
	if err != nil {
		return err
	}

	fileListQueue := make(chan *gocodewalker.File, 100)

	fileWalker := gocodewalker.NewFileWalker(input_path, fileListQueue)
	fileWalker.IgnoreGitIgnore = !(*parser.Config.Gitignore)
	fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, parser.Config.Whitelist...)

	// TODO: Fix blacklist implementation
	fileWalker.LocationExcludePattern = append(fileWalker.LocationExcludePattern, parser.Config.Blacklist...)

	errorHandler := func(e error) bool {
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() error {
		if err = fileWalker.Start(); err != nil {
			return err
		}
		return nil
	}()

	for f := range fileListQueue {
		fmt.Println(f.Location)

		if e := parser.readFile(f.Location); e != nil {
			return fmt.Errorf("could not read file at <%s>: %v", f.Location, e)
		}
	}

	return nil
}

// Read the file and find any requisite data. Pass this data to the Input channel in the parser
func (parser *Parser) readFile(path string) error {
	// Open the file for reading
	file, err := os.OpenFile(path, os.O_RDONLY, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))
	if err != nil {
		return err
	}
	defer file.Close()

	// Pull the file into the buffer
	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	// Create a temporary ToDo struct to hold the information
	todo := ToDo{RelativePath: path}

	// TODO: find a more elegant way to handle line number (e.g. range but afaik can't range on scanner.Scan())
	line_number := 1

	// Scan line-by-line searching for requisite callbacks
	// TODO: Create a config that a user can input which triggers they'd like to look for
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "TODO:") {
			todo.ToDo = append(todo.ToDo, fmt.Sprintf("Line %d: %s", line_number, strings.TrimSpace(line)))
		}
		line_number += 1
	}

	// Only if the array length is > 0 should the temporary ToDo struct be added to the Input channel
	if len(todo.ToDo) > 0 {
		parser.Input <- todo
	}

	// Return all errors
	return scanner.Err()
}
