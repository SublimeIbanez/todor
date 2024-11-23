package cmd

import (
	"fmt"
	"log"
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
	Args:    cobra.MinimumNArgs(1),
}

// Whitelist ***********************************************
var (
	whitelist_add    string
	whitelist_remove string
)

var whitelist_command = &cobra.Command{
	Use:     "whitelist",
	Short:   "Modify the whitelist",
	Aliases: []string{"wl"},
	Args:    cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		if len(whitelist_add) == 0 && len(whitelist_remove) == 0 {
			fmt.Println("Must provide an item to add or remove")
			os.Exit(1)
		}

		add := len(whitelist_add) != 0
		remove := len(whitelist_remove) != 0
		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load configuration: %s", err.Error())
		}

		// TODO: have some way of managing length limits on input
		if add {
			if err := cfg.AddToWhitelist(whitelist_add); err != nil {
				log.Fatalf("Could not add items to whitelist: %v", err)
			}
			fmt.Println("Added items to whitelist")
		}
		if remove {
			if err := cfg.RemoveFromWhitelist(whitelist_remove); err != nil {
				log.Fatalf("Could not remove items from whitelist: %v", err)
			}
			fmt.Println("Removed items from whitelist")
		}
	},
}

// Blacklist ***********************************************
var (
	blacklist_add    string
	blacklist_remove string
)

var blacklist_command = &cobra.Command{
	Use:     "blacklist",
	Short:   "Modify the blacklist",
	Aliases: []string{"bl"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(blacklist_add) == 0 && len(blacklist_remove) == 0 {
			fmt.Println("Must provide an item to add or remove")
			os.Exit(1)
		}

		add := len(blacklist_add) != 0
		remove := len(blacklist_remove) != 0
		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load configuration: %s", err.Error())
		}

		// TODO: have some way of managing length limits on input
		if add {
			cfg.AddToBlacklist(blacklist_add)
		}
		if remove {
			cfg.RemoveFromBlacklist(blacklist_remove)
		}
	},
}

// Git Ignore **********************************************
var git_ignore_command = &cobra.Command{
	Use:     "gitignore [true|false]",
	Short:   "Enable or disable the use of .gitignore when parsing [default: true]",
	Aliases: []string{"gi"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var useGitIgnore bool
		switch args[0] {
		case "true", "t":
			useGitIgnore = true
		case "false", "f":
			useGitIgnore = false
		default:
			fmt.Println("Invalid value. Use 'true' or 'false'.")
			return
		}

		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load configuration file: %s", err.Error())
		}
		if err = cfg.SetGitIgnore(&useGitIgnore); err != nil {
			log.Fatalf("Could not set the gitignore value: %v", err)
		}

		fmt.Printf("Set Git ignore to %v.\n", cfg.Gitignore)
	},
}

// Output Path *********************************************
var (
	set_output_path     string
	set_output_dir      string
	set_output_filename string
)
var output_directory_command = &cobra.Command{
	Use:     "output [path]",
	Short:   "Set the default output directory",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"out"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Must pass a valid path")
		}

		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load config: %v", err)
		}

		err = cfg.SetOutputPath(args[0])
		if err != nil {
			log.Fatalf("Could not set output path: %v", err)
		}
		fmt.Println("Successfully set output directory to")
	},
}

func init() {
	whitelist_command.Flags().StringVarP(&whitelist_add, "add", "a", "", "Add an item to the whitelist")
	whitelist_command.Flags().StringVarP(&whitelist_remove, "remove", "r", "", "Remove an item from the whitelist")
	config_command.AddCommand(whitelist_command)

	blacklist_command.Flags().StringVarP(&blacklist_add, "add", "a", "", "Add an item to the blacklist")
	blacklist_command.Flags().StringVarP(&blacklist_remove, "remove", "r", "", "Remove an item from the blacklist")
	config_command.AddCommand(blacklist_command)

	output_directory_command.Flags().StringVarP(&set_output_dir, "directory", "d", "", "Set the default output directory")
	output_directory_command.Flags().StringVarP(&set_output_filename, "filename", "f", "", "Set the default output filename")
	output_directory_command.Flags().StringVarP(&set_output_path, "path", "p", "", "Set the default output path")
	config_command.AddCommand(output_directory_command)

	config_command.AddCommand(git_ignore_command)

	root_command.AddCommand(config_command)
}
