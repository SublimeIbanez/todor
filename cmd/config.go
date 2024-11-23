package cmd

import (
	"fmt"
	"os"

	"github.com/SublimeIbanez/todor/configuration"
	"github.com/spf13/cobra"
)

// Config command ******************************************
var config_command = &cobra.Command{
	Use:     "config",
	Short:   "Modify the config file",
	Long:    `Use to add, remove, or edit configuration settings`,
	Aliases: []string{"cfg"},
}

// Whitelist ***********************************************
var (
	whitelist_add    string = ""
	whitelist_remove string = ""
)

var whitelist_command = &cobra.Command{
	Use:     "whitelist",
	Short:   "Modify the whitelist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"wl"},
}

var whitelist_add_command = &cobra.Command{
	Use:     "add",
	Short:   "Add an item to the whitelist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := configuration.LoadConfig()
		if err != nil {
			fmt.Println("Could not load configuration file:", err)
			os.Exit(1)
		}

		fmt.Println(cfg)

	},
}

var whitelist_remove_command = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an item from the whitelist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Whitelist ***********************************************
var (
	blacklist_add    string = ""
	blacklist_remove string = ""
)

var blacklist_command = &cobra.Command{
	Use:     "blacklist",
	Short:   "Modify the blacklist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"bl"},
}

var blacklist_add_command = &cobra.Command{
	Use:     "add",
	Short:   "Add an item to the blacklist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := configuration.LoadConfig()
		if err != nil {
			fmt.Println("Could not load configuration file:", err)
			os.Exit(1)
		}

		fmt.Println(cfg)

	},
}

var blacklist_remove_command = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an item from the blacklist",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Git Ignore **********************************************
var git_ignore_set bool = true

var git_ignore_command = &cobra.Command{
	Use:     "gitignore [true|false]",
	Short:   "Enable or disable the use of .gitignore when parsing [default: true]",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"gi"},
	Run: func(cmd *cobra.Command, args []string) {
		// if args[0] == "true" {
		// 	useGitIgnore = true
		// } else if args[0] == "false" {
		// 	useGitIgnore = false
		// } else {
		// 	fmt.Println("Invalid value. Use 'true' or 'false'.")
		// 	return
		// }
		// fmt.Printf("Set Git ignore to %v.\n", useGitIgnore)
	},
}

// Output Path *********************************************
var set_output_path string = ""

var default_output_path = &cobra.Command{
	Use:     "output_path [path]",
	Short:   "Set the default output directory",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"out", "outpath", "output", "out_path", "op"},
	Run: func(cmd *cobra.Command, args []string) {
		// defaultOutput = args[0]
		// fmt.Printf("Set default output directory to '%s'.\n", defaultOutput)
	},
}

func init() {
	config_command.AddCommand(whitelist_command)
	config_command.AddCommand(git_ignore_command)
	config_command.AddCommand(default_output_path)

	root_command.AddCommand(config_command)
}
