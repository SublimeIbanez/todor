package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TODO: Create a config file and have these constants placed there
const DEFAULT_OUTPUT_FILE_NAME string = "todos.md"
const DEFAULT_FILE_PERMISSIONS int = 0644

type Parser struct {
	Input      chan ToDo
	Output     chan string
	OutputFile *os.File
	Context    context.Context
	Cancel     context.CancelFunc
	WaitGroup  sync.WaitGroup
}

// Generates a new parser to manage i/o of the requisite data. Returns error if file operations fail
func NewParser(output_path string) (*Parser, error) {
	var output_file *os.File
	var err error

	if len(output_path) == 0 {
		output_path = DEFAULT_OUTPUT_FILE_NAME
	}

	path_info, err := os.Stat(output_path)
	if err != nil {
		if os.IsNotExist(err) {
			output_file, err = os.OpenFile(output_path, os.O_CREATE|os.O_WRONLY, fs.FileMode(DEFAULT_FILE_PERMISSIONS))
			if err != nil {
				return nil, fmt.Errorf("failed to create output file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to stat output path: %w", err)
		}
	} else if path_info.IsDir() {
		// Output path is a directory; create the default file inside it
		output_file_path := filepath.Join(output_path, DEFAULT_OUTPUT_FILE_NAME)
		output_file, err = os.OpenFile(output_file_path, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			return nil, fmt.Errorf("failed to open output file: %w", err)
		}
	} else {
		// Output path is a file
		output_file, err = os.OpenFile(output_path, os.O_CREATE|os.O_WRONLY, fs.FileMode(DEFAULT_FILE_PERMISSIONS))
		if err != nil {
			return nil, fmt.Errorf("failed to open output file: %w", err)
		}
	}

	// Create the context
	context, cancel := context.WithCancel(context.Background())

	// Create the parser
	parser := Parser{
		Input:      make(chan ToDo, 1000),
		Output:     make(chan string, 1000),
		OutputFile: output_file,
		Context:    context,
		Cancel:     cancel,
	}
	waiting := parser.init()
	parser.WaitGroup.Add(waiting)

	return &parser, nil
}

func (parser *Parser) init() int {
	// Input the count of goroutines created in init for adding to the wait group
	goroutines := 2

	go parser.handleInput()
	go parser.handleOutput()

	return goroutines
}

// Handles formatting the requisite data
func (parser *Parser) handleInput() {
	defer parser.WaitGroup.Done()
	for input := range parser.Input {
		select {
		case <-parser.Context.Done():
			return

		default:
			parser.Output <- fmt.Sprintf("[%s](%s)\n- %s\n\n", input.RelativePath, input.RelativePath, strings.Join(input.ToDo, "\n- "))
		}
	}
}

// Handles outputing the formatted data to the output file
func (parser *Parser) handleOutput() {
	defer parser.WaitGroup.Done()
	for output := range parser.Output {
		select {
		case <-parser.Context.Done():
			return

		default:
			output_data := []byte(output)
			parser.OutputFile.Write(output_data)
		}
	}
}

func (parser *Parser) Shutdown() {
	parser.Cancel()

	close(parser.Input)
	close(parser.Output)
	if parser.OutputFile != nil {
		parser.OutputFile.Close()
	}
}
