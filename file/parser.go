package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SublimeIbanez/todor/common"
)

// TODO: Create a config file and have these constants placed there
const DEFAULT_OUTPUT_FILE_NAME string = "todos.md"

type Parser struct {
	Input      chan ToDo
	Output     chan string
	OutputFile *os.File
	Context    context.Context
	Cancel     context.CancelFunc
	WaitGroup  sync.WaitGroup
	Config     common.ConfigOptions
}

// Generates a new parser to manage i/o of the requisite data. Returns error if file operations fail
func NewParser(output_path string) (*Parser, error) {
	var cfg common.ConfigOptions
	var output_file *os.File
	var err error

	cfg, err = common.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(output_path) == 0 {
		output_dir, err := filepath.Abs(cfg.OutputDirectory)
		if err != nil {
			return nil, err
		}

		output_path = filepath.Join(output_dir, DEFAULT_OUTPUT_FILE_NAME)
	}

	path_info, err := os.Stat(output_path)
	if err != nil {
		if os.IsNotExist(err) {
			output_file, err = os.OpenFile(output_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))
			if err != nil {
				return nil, fmt.Errorf("failed to create output file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to stat output path: %w", err)
		}
	} else if path_info.IsDir() {
		// Output path is a directory; create the default file inside it
		output_file_path := filepath.Join(output_path, DEFAULT_OUTPUT_FILE_NAME)
		output_file, err = os.OpenFile(output_file_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))

		if err != nil {
			return nil, fmt.Errorf("failed to open output file: %w", err)
		}
	} else {
		// Output path is a file
		output_file, err = os.OpenFile(output_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(common.DEFAULT_FILE_PERMISSIONS))
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
		Config:     cfg,
	}
	parser.init()

	return &parser, nil
}

func (parser *Parser) init() {
	// Input the count of goroutines created in init for adding to the wait group
	goroutines := 2
	parser.WaitGroup.Add(goroutines)

	go parser.handleInput()
	go parser.handleOutput()
}

// Handles formatting the requisite data
func (parser *Parser) handleInput() {
	defer parser.WaitGroup.Done()

	for input := range parser.Input {
		select {
		case <-parser.Context.Done():
			close(parser.Output)
			return

		default:
			parser.Output <- fmt.Sprintf("[%s](%s)\n- %s\n\n", input.RelativePath, input.RelativePath, strings.Join(input.ToDo, "\n- "))
		}
	}
	close(parser.Output)
}

// Handles outputing the formatted data to the output file
func (parser *Parser) handleOutput() {
	defer parser.WaitGroup.Done()

	for output := range parser.Output {

		fmt.Println(output)
		_, err := parser.OutputFile.Write([]byte(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to output file: %v\n", err)
			return
		}
	}
}

func (parser *Parser) Shutdown() {
	close(parser.Input)
	parser.WaitGroup.Wait()

	if parser.OutputFile != nil {
		parser.OutputFile.Close()
	}
}
