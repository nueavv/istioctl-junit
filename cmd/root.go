/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"io"
	"fmt"
	"errors"
	"encoding/json"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	report "github.com/nueavv/istioctl-junit/report"
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
		var junitReport []report.JunitReport

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
			junitReport, _ = Yaml2JunitReport(data)
		case "json":
			junitReport, _ = Json2JunitReport(data)
		default:
			return errors.New("")
		}

		return report.MakeReport(junitReport, output)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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


// func LoadYamlJunitReport(filepath string) (report.JunitReport, error) {
// 	var report report.YamlJunitReport
//     yamlFile, _ := os.ReadFile(filepath)
//     err := yaml.Unmarshal(yamlFile, &report)
//     if err != nil {
// 		return report, err
//     }
	
// 	return report.JunitReport(report), nil
// }

func Yaml2JunitReport(data string) ([]report.JunitReport, error) {
	var junitReport []report.JunitReport
	var yamljunitReport []report.YamlJunitReport
	err := yaml.Unmarshal([]byte(data), &yamljunitReport)
    if err != nil {
		return nil, err
    }

	for _, r := range yamljunitReport {
		junitReport = append(junitReport, r)
	}
	return junitReport, nil
}


// func LoadJsonJunitReport(filepath string) (report.JunitReport, error) {
// 	var report report.JsonJunitReport
//     yamlFile, _ := os.ReadFile(filepath)
//     err := yaml.Unmarshal(yamlFile, &report)
//     if err != nil {
// 		return report, err
//     }
// 	return report.JunitReport(report), nil
// }

func Json2JunitReport(data string) ([]report.JunitReport, error) {
	var junitReport []report.JunitReport
	var jsonjunitReport []report.JsonJunitReport
	err := json.Unmarshal([]byte(data), &jsonjunitReport)
    if err != nil {
		return nil, err
    }

	for _, r := range jsonjunitReport {
		junitReport = append(junitReport, r)
	}
	return junitReport, nil
}