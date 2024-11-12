package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/SublimeIbanez/todor/common"
)

// Recursively walk through the directory and read through all items
func (parser *Parser) WalkDir(input_path string) error {
	defer parser.Context.Done()
	input, err := os.Stat(input_path)
	if err != nil {
		return err
	}

	full_path, err := filepath.Abs(input_path)
	if err != nil {
		return err
	}

	if input.IsDir() {
		err = filepath.WalkDir(full_path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				if err = parser.readFile(path); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking the directory: %v", err)
		}
	} else {
		if err := parser.readFile(full_path); err != nil {
			return err
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
