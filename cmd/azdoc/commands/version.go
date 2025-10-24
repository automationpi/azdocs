package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display the version, git commit, and build date of azdoc",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
