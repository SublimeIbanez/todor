package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/SublimeIbanez/todor/common"
	"github.com/boyter/gocodewalker"
)

// Recursively walk through the directory and read through all items
func (parser *Parser) WalkDir(input_path string) error {
	full_path, err := filepath.Abs(input_path)
	if err != nil {
		return err
	}

	_, err = os.Stat(input_path)
	if err != nil {
		return err
	}

	fileListQueue := make(chan *gocodewalker.File, 100)

	fileWalker := gocodewalker.NewFileWalker(full_path, fileListQueue)
	fileWalker.IgnoreGitIgnore = !(*parser.Config.Gitignore)

	var extension_list []string
	for _, e := range parser.Config.Whitelist {
		extension_list = append(extension_list, strings.Replace(e, ".", "", 1))
	}
	fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, extension_list...)
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
		if e := parser.readFile(f.Location); e != nil {
			return fmt.Errorf("could not read file at <%s>: %v", f.Location, e)
		}
	}

	return nil
}

// Read the file and find any requisite data. Pass this data to the Input channel in the parser
func (parser *Parser) readFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	todo := ToDo{RelativePath: path}

	// TODO: Create a config that a user can input which triggers they'd like to look for
	line_number := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "todo") {
			todo.ToDo = append(todo.ToDo, fmt.Sprintf("Line %d: %s", line_number, strings.TrimSpace(line)))
		}
		line_number += 1
	}

	if len(todo.ToDo) > 0 {
		parser.Input <- todo
	}

	return scanner.Err()
}
