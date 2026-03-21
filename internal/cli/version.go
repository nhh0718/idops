package cli

import (
	"fmt"
	"runtime"

	"github.com/nhh0718/idops/internal/ui"
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
		fmt.Println(ui.TitleStyle.Render("idops") + " " + version)
		fmt.Printf("  Commit:   %s\n", commit)
		fmt.Printf("  Built:    %s\n", date)
		fmt.Printf("  Go:       %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:  %s/%s\n", runtime.GOOS, runtime.GOARCH)
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

// GetVersion returns the current version string.
func GetVersion() string {
	return version
}
