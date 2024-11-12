package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
				if err = parser.ReadFile(path); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking the directory: %v", err)
		}
	} else {
		if err := parser.ReadFile(full_path); err != nil {
			return err
		}
	}

	return nil
}

func (parser *Parser) ReadFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	todo := ToDo{RelativePath: path}
	line_number := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "TODO:") {
			todo.ToDo = append(todo.ToDo, fmt.Sprintf("Line %d: %s", line_number, strings.TrimSpace(line)))
		}
		line_number += 1
	}
	if len(todo.ToDo) > 0 {
		parser.Input <- todo
	}

	return scanner.Err()
}
