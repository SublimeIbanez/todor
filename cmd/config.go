package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/SublimeIbanez/todor/configuration"
	"github.com/spf13/cobra"
)

// TODO: For all: better error management and handling rather than log.Fatal

// Config command ******************************************
var config_command = &cobra.Command{
	Use:     "config",
	Short:   "Modify the config file",
	Long:    `Use to add, remove, or edit configuration settings`,
	Aliases: []string{"cfg"},
	Args:    cobra.NoArgs,
}

// Whitelist ***********************************************
var (
	whitelist_add    []string
	whitelist_remove []string
)

var whitelist_command = &cobra.Command{
	Use:     "whitelist",
	Short:   "Modify the whitelist",
	Aliases: []string{"wl"},
	Args:    cobra.NoArgs,
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

		// TODO: Rework error handling so log.Fatal isn't used all the time
		if add {
			for _, item := range whitelist_add {
				if err := cfg.AddToWhitelist(item); err != nil {
					log.Fatalf("Could not add item <%s> to whitelist: %v", item, err)
				}
			}
			fmt.Println("Added items to whitelist")
		}
		if remove {
			for _, item := range whitelist_remove {
				if err := cfg.RemoveFromWhitelist(item); err != nil {
					log.Fatalf("Could not remove item <%s> from whitelist: %v", item, err)
				}
			}
			fmt.Println("Removed items from whitelist")
		}
	},
}

// Blacklist ***********************************************
var (
	blacklist_add    []string
	blacklist_remove []string
)

var blacklist_command = &cobra.Command{
	Use:     "blacklist",
	Short:   "Modify the blacklist",
	Aliases: []string{"bl"},
	Args:    cobra.NoArgs,
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

		// TODO: Rework error handling so log.Fatal isn't used all the time
		if add {
			for _, item := range blacklist_add {
				if err := cfg.AddToBlacklist(item); err != nil {
					log.Fatalf("Could not add item <%s> to blacklist: %v", item, err)
				}
			}
		}
		if remove {
			for _, item := range blacklist_remove {
				if err := cfg.RemoveFromBlacklist(item); err != nil {
					log.Fatalf("Could not remove item <%s> from blacklist: %v", item, err)
				}
			}
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
			log.Fatal("Invalid value. Use 'true', 't', 'false' or 'f'.")
			return
		}

		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load configuration file: %v", err)
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
var set_output_command = &cobra.Command{
	Use:     "output",
	Short:   "Set the default output path, directory, and/or filename",
	Args:    cobra.NoArgs,
	Aliases: []string{"out"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(set_output_path) == 0 && len(set_output_dir) == 0 && len(set_output_filename) == 0 {
			log.Fatal("Path cannot be empty")
		}

		handle_output_path := len(set_output_path) != 0
		handle_output_dir := len(set_output_path) != 0
		handle_output_filename := len(set_output_path) != 0

		cfg, err := configuration.LoadConfig()
		if err != nil {
			log.Fatalf("Could not load config: %v", err)
		}

		if handle_output_path {
			log.Fatalf("Setting the output path is currently under development")
		} else {
			if handle_output_dir {
				if err = cfg.SetOutputDirectory(set_output_dir); err != nil {
					log.Fatalf("Could not set output directory: %v", err)
				}
				fmt.Printf("Successfully set output directory: %s\n", set_output_dir)
			}

			if handle_output_filename {
				if err = cfg.SetOutputFilename(set_output_filename); err != nil {
					log.Fatalf("Could not set output filename: %v", err)
				}
				fmt.Printf("Successfully set output filename: %s\n", set_output_filename)
			}
		}

		fmt.Println("Successfully set output directory to")
	},
}

func init() {
	whitelist_command.Flags().StringSliceVarP(&whitelist_add, "add", "a", nil, "Add an item to the whitelist")
	whitelist_command.Flags().StringSliceVarP(&whitelist_remove, "remove", "r", nil, "Remove an item from the whitelist")
	config_command.AddCommand(whitelist_command)

	blacklist_command.Flags().StringSliceVarP(&blacklist_add, "add", "a", nil, "Add an item to the blacklist")
	blacklist_command.Flags().StringSliceVarP(&blacklist_remove, "remove", "r", nil, "Remove an item from the blacklist")
	config_command.AddCommand(blacklist_command)

	set_output_command.Flags().StringVarP(&set_output_dir, "directory", "d", "", "Set the default output directory")
	set_output_command.Flags().StringVarP(&set_output_filename, "filename", "f", "", "Set the default output filename")
	set_output_command.Flags().StringVarP(&set_output_path, "path", "p", "", "Set the default output path")
	config_command.AddCommand(set_output_command)

	config_command.AddCommand(git_ignore_command)

	root_command.AddCommand(config_command)
}
