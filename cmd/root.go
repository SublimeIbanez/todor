/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todor",
	Short: "Gathers all TODO's into one easy-to-read markdown file.",
	Long: `Walks through the provided directory (or starts from root directory if none provided) and searches for all instances
of "TODO" within any files. Lists these "TODO"s in a markdown file with links to that file.

Example:
    todor -p src`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var path string
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "used to input a specified path")
}
