/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"io"
	"fmt"
	"errors"
	"encoding/xml"

	"github.com/spf13/cobra"

	"github.com/nueavv/istioctl-junit/utils/converter"
	"github.com/nueavv/istioctl-junit/utils/junit"
)

var (
	output string
	format string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "istioctl-junit",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var inputReader io.Reader = cmd.InOrStdin()
		var junitReport []converter.JunitReport

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
			junitReport, _ = converter.Yaml2JunitReport(data)
		case "json":
			junitReport, _ = converter.Json2JunitReport(data)
		default:
			return errors.New("")
		}

		return MakeReport(junitReport, output)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&format, "format", "f", "", "istioctl analyze format <json|yaml>")
	rootCmd.Flags().StringVarP(&output, "output", "o", "report.xml", "report filename")
	rootCmd.MarkFlagRequired("format")
}


func MakeReport[T converter.JunitReport](reports []T, output string) error {
	var testsuite junit.TestSuite
	for _, report := range reports {
		var testcase *junit.TestCase
		testcase = &junit.TestCase{
			Name: report.GetOrigin(),
		}

		switch report.GetLevel() {
		case converter.StatusError:
			testcase.Errors = append(testcase.Errors, &junit.Error{
				Message: report.GetMessage(),
				Type: report.GetCode(),
			})
			if report.IsFileAnalze() {
				testcase.Errors[0].File = report.GetErrorFile()
				testcase.Errors[0].Line = report.GetErrorLine()
			}
		case converter.StatusFailed:
			testcase.Failures = append(testcase.Failures, &junit.Failure{
				Message: report.GetMessage(),
				Type: report.GetCode(),
			})
		}

		
		testsuite.TestCases = append(testsuite.TestCases, testcase)
	}

	testsuite.Name = "istioctl analyze"
	testsuite.Errors = converter.GetErrorCount(reports)
	testsuite.Failures = converter.GetWarningCount(reports)

	xmlBytes, err := xml.Marshal(testsuite)
    if err != nil {
        return err
    }

	error := os.WriteFile(output, xmlBytes, 0660)
	if error != nil {
		return error
	}

	return nil
}