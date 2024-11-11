package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// Recursively walk through the directory and read through all items
func (parser *Parser) WalkDir(input_path string) error {
	input, _ := os.Stat(input_path)
	// convert the input_path to a []DirEntry because there's no simple function that I know of to do so
	if input.IsDir() {
		full_path, err := filepath.Abs(input_path)
		items, err := os.ReadDir(full_path)
		if err != nil {
			return err
		}

		for _, item := range items {
			entry, err := item.Info()
			if err != nil {
				return err
			}
			err = parser.WalkDir(entry.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (parser *Parser) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return scanner.Err()
}
