package cmd

import (
	"github.com/spf13/cobra"
)

var config_command = &cobra.Command{
	Use:   "config",
	Short: "Modify the config file",
	Long:  `Use to add, remove, or edit configuration settings`,
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the config

		// Set the config

		// Save the config
	},
}

func init() {
	config_command.Flags().StringVarP(&input_path, "path", "p", ".", "used to specify root path -- default is current working directory")
	config_command.Flags().StringVarP(&output_path, "output", "o", "", "used to specify output path -- default is current working directory")
}
