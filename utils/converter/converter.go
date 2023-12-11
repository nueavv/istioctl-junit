package converter

import (
	"strconv"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Status string

const (
	StatusPassed  Status = "passed"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "Warning"
	StatusError   Status = "Error"
)

type JunitReport interface {
	GetCode() string
	GetMessage() string
	GetOrigin() string
	GetURL() string
	GetLevel() Status
	IsFileAnalze() bool
	GetErrorLine() int
	GetErrorFile() string
}

func GetErrorCount[T JunitReport](reports []T) int {
	count := 0
	for _, r := range reports {
		if r.GetLevel() == StatusError {
			count++
		}
	}
	return count
}

func GetWarningCount[T JunitReport](reports []T) int {
	count := 0
	for _, r := range reports {
		if r.GetLevel() == StatusFailed {
			count++
		}
	}
	return count
}

func GetTotal[T JunitReport](reports []T) int {
	return len(reports)
}

func Yaml2JunitReport(data string) ([]JunitReport, error) {
	var junitReport []JunitReport
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

func Json2JunitReport(data string) ([]JunitReport, error) {
	var junitReport []JunitReport
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

type YamlJunitReport struct {
	Code             string `yaml:"code"`
	Message          string `yaml:"message"`
	Origin           string `yaml:"origin"`
	DocumentationURL string `yaml:"documentationUrl"`
	Level            string `yaml:"level"`
	Reference        string `yaml:"reference,omitempty"`
}

type JsonJunitReport struct {
	Code             string `json:"code"`
	Message          string `json:"message"`
	Origin           string `json:"origin"`
	DocumentationURL string `json:"documentationUrl"`
	Level            string `json:"level"`
	Reference        string `json:"reference,omitempty"`
}

func (y YamlJunitReport) GetCode() string {
	return y.Code
}

func (y YamlJunitReport) GetMessage() string {
	return y.Message
}

func (y YamlJunitReport) GetOrigin() string {
	return y.Origin
}

func (y YamlJunitReport) GetURL() string {
	return y.DocumentationURL
}

func (y YamlJunitReport) GetLevel() Status {
	return Status(y.Level)
}

func (y YamlJunitReport) IsFileAnalze() bool {
	return y.Reference != ""
}

func (y YamlJunitReport) GetErrorLine() int {
	line, _ := strconv.Atoi(strings.Split(y.Reference, ":")[1])
	return line
}

func (y YamlJunitReport) GetErrorFile() string {
	return strings.Split(y.Reference, ":")[0]
}

func (j JsonJunitReport) GetCode() string {
	return j.Code
}

func (j JsonJunitReport) GetMessage() string {
	return j.Message
}

func (j JsonJunitReport) GetOrigin() string {
	return j.Origin
}

func (j JsonJunitReport) GetURL() string {
	return j.DocumentationURL
}

func (j JsonJunitReport) GetLevel() Status {
	return Status(j.Level)
}

func (j JsonJunitReport) IsFileAnalze() bool {
	return j.Reference != ""
}

func (j JsonJunitReport) GetErrorLine() int {
	line, _ := strconv.Atoi(strings.Split(j.Reference, ":")[1])
	return line
}

func (j JsonJunitReport) GetErrorFile() string {
	return strings.Split(j.Reference, ":")[0]
}
