package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"istioctl-junit/internal/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	v := version.GetVersion()
	fmt.Fprint(cmd.OutOrStdout(), formatVersion(&v))
}

func formatVersion(v *version.Version) string {
	var output strings.Builder

	fmt.Fprintf(&output, "%s: %s\n", config.CLIName, v)
	fmt.Fprintf(&output, "  BuildDate: %s\n", v.BuildDate)
	fmt.Fprintf(&output, "  GitCommit: %s\n", v.GitCommit)
	fmt.Fprintf(&output, "  GitTreeState: %s\n", v.GitTreeState)
	if v.GitTag != "" {
		fmt.Fprintf(&output, "  GitTag: %s\n", v.GitTag)
	}
	fmt.Fprintf(&output, "  GoVersion: %s\n", v.GoVersion)
	fmt.Fprintf(&output, "  Compiler: %s\n", v.Compiler)
	fmt.Fprintf(&output, "  Platform: %s\n", v.Platform)
	if v.ExtraBuildInfo != "" {
		fmt.Fprintf(&output, "  ExtraBuildInfo: %s\n", v.ExtraBuildInfo)
	}

	return output.String()
}
