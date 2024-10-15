package istio2junit

import (
	"strconv"
	"strings"

	"encoding/json"
	"istioctl-junit/internal/junit"
	"gopkg.in/yaml.v3"
)

type Status string

const (
	StatusPassed  Status = "Passed"
	StatusSkipped Status = "Skipped"
	StatusFailed  Status = "Warning"
	StatusError   Status = "Error"
)

type JunitReportConverter interface {
	GetCode() string
	GetMessage() string
	GetOrigin() string
	GetURL() string
	GetLevel() Status
	IsFileAnalyze() bool
	GetErrorLine() int
	GetErrorFile() string
}

func CountReportsByStatus[T JunitReportConverter](reports []T, status Status) int {
	count := 0
	for _, r := range reports {
		if r.GetLevel() == status {
			count++
		}
	}
	return count
}

func GetErrorCount[T JunitReportConverter](reports []T) int {
	return CountReportsByStatus(reports, StatusError)
}

func GetWarningCount[T JunitReportConverter](reports []T) int {
	return CountReportsByStatus(reports, StatusFailed)
}

func GetTotal[T JunitReportConverter](reports []T) int {
	return len(reports)
}

func Yaml2JunitReport(data string) ([]JunitReportConverter, error) {
	var junitReport []JunitReportConverter
	var yamljunitReport []YamlJunitReport
	err := yaml.Unmarshal([]byte(data), &yamljunitReport)
	if err != nil {
		return nil, err
	}

	for _, r := range yamljunitReport {
		junitReport = append(junitReport, r)
	}
	return junitReport, nil
}

func Json2JunitReport(data string) ([]JunitReportConverter, error) {
	var junitReport []JunitReportConverter
	var jsonjunitReport []JsonJunitReport
	err := json.Unmarshal([]byte(data), &jsonjunitReport)
	if err != nil {
		return nil, err
	}

	for _, r := range jsonjunitReport {
		junitReport = append(junitReport, r)
	}
	return junitReport, nil
}

type BaseJunitReport struct {
	Code             string `yaml:"code" json:"code"`
	Message          string `yaml:"message" json:"message"`
	Origin           string `yaml:"origin" json:"origin"`
	DocumentationURL string `yaml:"documentationUrl" json:"documentationUrl"`
	Level            string `yaml:"level" json:"level"`
	Reference        string `yaml:"reference,omitempty" json:"reference,omitempty"`
}

func (b BaseJunitReport) GetCode() string             { return b.Code }
func (b BaseJunitReport) GetMessage() string          { return b.Message }
func (b BaseJunitReport) GetOrigin() string           { return b.Origin }
func (b BaseJunitReport) GetURL() string              { return b.DocumentationURL }
func (b BaseJunitReport) GetLevel() Status            { return Status(b.Level) }
func (b BaseJunitReport) IsFileAnalyze() bool         { return b.Reference != "" }
func (b BaseJunitReport) GetErrorLine() int {
	line, _ := strconv.Atoi(strings.Split(b.Reference, ":")[1])
	return line
}
func (b BaseJunitReport) GetErrorFile() string {
	return strings.Split(b.Reference, ":")[0]
}

type YamlJunitReport BaseJunitReport
type JsonJunitReport BaseJunitReport

func MakeReport[T JunitReportConverter](reports []T) junit.TestSuite {
	var testsuite junit.TestSuite
	for _, report := range reports {
		testcase := &junit.TestCase{
			Name: report.GetOrigin(),
		}

		switch report.GetLevel() {
		case StatusError:
			testcase.Errors = append(testcase.Errors, &junit.Error{
				Message: report.GetMessage(),
				Type:    report.GetCode(),
			})
			if report.IsFileAnalyze() {
				testcase.Errors[0].File = report.GetErrorFile()
				testcase.Errors[0].Line = report.GetErrorLine()
				testcase.Errors[0].Reference = report.GetReference()
			}
		case StatusFailed:
			testcase.Failures = append(testcase.Failures, &junit.Failure{
				Message: report.GetMessage(),
				Type:    report.GetCode(),
			})
		}
		testcase.Status = string(report.GetLevel())
		testsuite.TestCases = append(testsuite.TestCases, testcase)
	}

	testsuite.Name = "istioctl analyze"
	testsuite.Errors = GetErrorCount(reports)
	testsuite.Failures = GetWarningCount(reports)
	return testsuite
}
