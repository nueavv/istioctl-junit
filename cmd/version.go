/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nueavv/istioctl-junit/common"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cv := common.GetVersion()
		fmt.Fprint(cmd.OutOrStdout(), printVersion(&cv))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(version *common.Version) string {
	output := fmt.Sprintf("%s: %s\n", cliName, version)
	output += fmt.Sprintf("  BuildDate: %s\n", version.BuildDate)
	output += fmt.Sprintf("  GitCommit: %s\n", version.GitCommit)
	output += fmt.Sprintf("  GitTreeState: %s\n", version.GitTreeState)
	if version.GitTag != "" {
		output += fmt.Sprintf("  GitTag: %s\n", version.GitTag)
	}
	output += fmt.Sprintf("  GoVersion: %s\n", version.GoVersion)
	output += fmt.Sprintf("  Compiler: %s\n", version.Compiler)
	output += fmt.Sprintf("  Platform: %s\n", version.Platform)
	if version.ExtraBuildInfo != "" {
		output += fmt.Sprintf("  ExtraBuildInfo: %s\n", version.ExtraBuildInfo)
	}
	return output
}
