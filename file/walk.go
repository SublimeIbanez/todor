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

var buffer = make([]byte, 0, 64*1024)

const buffer_size int = 1024 * 1024

// Recursively walk through the directory and read through all items
func (parser *Parser) WalkDir(input_path string) error {
	input, err := os.Stat(input_path)
	if err != nil {
		return err
	}

	full_path, err := filepath.Abs(input_path)
	if err != nil {
		return err
	}

	ignore_list := parser.Config.Blacklist
	// Find .gitgnore if UseGitIgnore from the config file is correct -- extract into ignore files
	if *parser.Config.Gitignore {
		filepath.WalkDir(full_path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.Name() == string(common.GIT_IGNORE) {
				file, e := os.OpenFile(path, os.O_RDONLY, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))
				if e != nil {
					return e
				}

				scanner := bufio.NewScanner(file)
				scanner.Buffer(buffer, buffer_size)

				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					// Ignore comments
					if len(line) == 0 || strings.HasPrefix(line, string(common.GIT_IGNORE_COMMENT_PREFIX)) {
						continue
					}

					// TODO: Handle `*` properly
					ignore_list = append(ignore_list, strings.Replace(line, "*", "", -1))
				}

				file.Close()
				return scanner.Err()
			}

			return nil
		})
	}
	fmt.Println(ignore_list)

	if input.IsDir() {
		err = filepath.WalkDir(full_path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			marked_ignore := true
			if len(ignore_list) > 0 {
				for _, ignore := range ignore_list {
					if strings.Contains(path, ignore) {
						marked_ignore = true
					}
				}
			}

			if !d.IsDir() && !marked_ignore && len(parser.Config.Whitelist) > 0 {
				for _, allowed := range parser.Config.Whitelist {
					if strings.Contains(path, allowed) {
						if e := parser.readFile(path); e != nil {
							return fmt.Errorf("could not read file at <%s>: %v", path, e)
						}
					}
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
