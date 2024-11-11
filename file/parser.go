package file

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const DEFAULT_OUTPUT_FILE_NAME string = "todos.md"

type Parser struct {
	Input      chan ToDo
	Output     chan string
	OutputFile *os.File
	Context    context.Context
	Cancel     context.CancelFunc
}

func NewParser(output_path string) *Parser {
	var output_file *os.File
	var err error

	// Handle the output file
	if output_path.IsDir() {
		f, err := filepath.Abs(output_path.Name())
		if err != nil {
			fmt.Println("Could not generate output absolute path:", err)
			os.Exit(2)
		}

		f = filepath.Join(f, DEFAULT_OUTPUT_FILE_NAME)

		output_file, err = os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}

	fmt.Println(output_path)

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
	parser.init()

	return &parser
}

func (parser *Parser) init() {
	go parser.handleInput()
	go parser.handleOutput()
}

func (parser *Parser) handleInput() {
	for input := range parser.Input {
		select {
		case <-parser.Context.Done():
			return

		default:
			fmt.Println(input)
		}
	}
}

func (parser *Parser) handleOutput() {
	for output := range parser.Output {
		select {
		case <-parser.Context.Done():
			return

		default:
			fmt.Println(output)
		}
	}
}

func (parser *Parser) Shutdown() {
	parser.Cancel()

	close(parser.Input)
	close(parser.Output)
	defer parser.OutputFile.Close()
}
