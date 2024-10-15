package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"istioctl-junit/internal/istio2junit"
	"istioctl-junit/internal/junit"
)

var (
	checkResult bool
	outputFile  string
	inputFormat string
	testName    string
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
		inputReader, err := getInputReader(cmd, args)
		if err != nil {
			return fmt.Errorf("입력 파일 열기 실패: %v", err)
		}

		data, err := io.ReadAll(inputReader)
		if err != nil {
			return fmt.Errorf("파일 읽기 실패: %v", err)
		}

		junitReport, err := convertToJunitReport(string(data), inputFormat)
		if err != nil {
			return err
		}

		testsuite := istio2junit.MakeReport(junitReport)
		testsuite.Name = testName

		if err := junit.WriteFile(testsuite, outputFile); err != nil {
			return fmt.Errorf("파일 쓰기 실패: %v", err)
		}

		fmt.Printf("출력 파일: %s\n", outputFile)

		if checkResult {
			return checkAnalyzeResult(testsuite, cmd)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&checkResult, "check", "c", false, "리포트 상태 확인")
	rootCmd.Flags().StringVarP(&testName, "name", "n", "istio-analyze", "istioctl 분석 클러스터 이름")
	rootCmd.Flags().StringVarP(&inputFormat, "format", "f", "", "istioctl 분석 형식 <json|yaml>")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "report.xml", "리포트 파일 이름")
	if err := rootCmd.MarkFlagRequired("format"); err != nil {
		fmt.Printf("형식 플래그 오류: %v", err)
	}
}
