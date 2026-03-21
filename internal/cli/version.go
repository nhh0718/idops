package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build-time variables injected via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("idops %s (%s) built %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// SetVersionInfo sets build info from main package ldflags.
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}
