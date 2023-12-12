/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/nueavv/istioctl-junit/utils/istio2junit"
	"github.com/nueavv/istioctl-junit/utils/junit"
)

var (
	check_result bool
	// istio_analyzed_result bool
	output string
	format string
	test_name string
)

const (
	cliName = "istioctl-junit"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     cliName,
	Args:    cobra.MaximumNArgs(1),
	Example: `hello`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputReader io.Reader = cmd.InOrStdin()
		var junitReport []istio2junit.JunitReportConverter

		// the argument received looks like a file, we try to open it
		if len(args) > 0 && args[0] != "-" {
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed open file: %v", err)
			}
			inputReader = file
		}

		raw, err := io.ReadAll(inputReader)
		if err != nil {
			return fmt.Errorf("failed read file: %v", err)
		}
		data := string(raw)

		switch format {
		case "yaml":
			junitReport, _ = istio2junit.Yaml2JunitReport(data)
		case "json":
			junitReport, _ = istio2junit.Json2JunitReport(data)
		default:
			return errors.New("")
		}
		testsuite := istio2junit.MakeReport(junitReport)
		testsuite.Name = test_name
		err = junit.WriteFile(testsuite, output)
		if err != nil {
			return fmt.Errorf("Failed WriteFile %v", err)
		}
		fmt.Printf("Output File : %s\n", output)

		switch check_result {
		case true:
			if testsuite.Errors > 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("Analyze Result Total: %d, Skipped: %d, Failed: %d, Error: %d\n", testsuite.Tests, testsuite.Skipped, testsuite.Failures, testsuite.Errors)
			}
			fmt.Printf("Analyze Result Total: %d, Skipped: %d, Failed: %d, Error: %d\n", testsuite.Tests, testsuite.Skipped, testsuite.Failures, testsuite.Errors)
		}
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&check_result, "check", "c", false, "check report status.")
	rootCmd.Flags().StringVarP(&test_name, "name", "n", "istio-analyze", "istioctl analyze cluster name")
	rootCmd.Flags().StringVarP(&format, "format", "f", "", "istioctl analyze format <json|yaml>")
	rootCmd.Flags().StringVarP(&output, "output", "o", "report.xml", "report filename")
	if err := rootCmd.MarkFlagRequired("format"); err != nil {
		fmt.Printf("error format flag :%v", err)
	}
}
