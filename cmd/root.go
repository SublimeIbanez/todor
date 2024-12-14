/*
Copyright © 2024 Joshua Benn sublimeibanez@protonmail.com and Jon-Micheal Hartway jhartway99@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/SublimeIbanez/todor/file"
	"github.com/spf13/cobra"
)

var input_path string
var output_path string

// root_command represents the base command when called without any subcommands
var root_command = &cobra.Command{
	Use:   "todor",
	Short: "Gathers all TODO's into one easy-to-read markdown file.",
	Long: `Walks through the provided directory (or starts from current working directory if none provided) and searches for all instances of "TODO" within any files. Lists these "TODO"s in a markdown file with links to that file.

Example:
    todor -p src -o output.md`,
	Args: cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// Create the parser
		parser, err := file.NewParser(output_path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = parser.WalkDir(input_path)
		if err != nil {
			fmt.Println("Could not walk directory or file:", err)
			os.Exit(1)
		}

		parser.Shutdown()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := root_command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = config_command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	root_command.Flags().StringVarP(&input_path, "path", "p", ".", "used to specify root path -- default is current working directory")
	root_command.Flags().StringVarP(&output_path, "output", "o", "", "used to specify output path -- default is current working directory")
}
