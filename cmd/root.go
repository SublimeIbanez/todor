/*
Copyright © 2024 Joshua Benn sublimeibanez@protonmail.com
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/SublimeIbanez/todor/file"
	"github.com/spf13/cobra"
)

var input_path string
var output_path string

// root_command represents the base command when called without any subcommands
var root_command = &cobra.Command{
	Use:   "",
	Short: "Gathers all TODO's into one easy-to-read markdown file.",
	Long: `Walks through the provided directory (or starts from root directory if none provided) and searches for all instances
of "TODO" within any files. Lists these "TODO"s in a markdown file with links to that file.

Example:
    todor -p src`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {

		// Check the path
		input, err := os.Stat(input_path)
		if err != nil {
			fmt.Println("Could not verify input path:", err)
			os.Exit(1)
		}
		output, err := os.Stat(input_path)
		if err != nil {
			fmt.Println("Could not verify output path:", err)
			os.Exit(1)
		}

		fmt.Println("Input", input, " -- Output", output)

		// Create the parser
		parser := file.NewParser(fs.FileInfoToDirEntry(output))
		defer parser.Shutdown()

		err = parser.WalkDir(input_path)
		if err != nil {
			fmt.Println("Could not walk directory or file", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := root_command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	root_command.Flags().StringVarP(&input_path, "path", "p", ".", "used to specify root path -- default is current working directory")
	root_command.Flags().StringVarP(&output_path, "output", "o", "todos.md", "used to specify output path -- default is current working directory")
}
