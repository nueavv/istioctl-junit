package report

import (
	"strings"
	"errors"
	"strconv"
)

type YamlJunitReport struct {
	Code             string `yaml:"code"`
	Message          string `yaml:"message"`
	Origin           string `yaml:"origin"`
	DocumentationURL string `yaml:"documentationUrl"`
	Level            string `yaml:"level"`
	Reference		 string `yaml:"reference,omitempty"`
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

func (y YamlJunitReport) GetLevel() string {
	return y.Level
}

func (y YamlJunitReport) IsFileAnalze() bool {
	if y.Reference == "" {
		return false
	}
	return true
}

func (y YamlJunitReport) GetErrorLine() (int, error) {
	if y.IsFileAnalze() {
		line := strings.Split(y.Reference, ":")[1]
		i, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, errors.New("No Reference")
}

func (y YamlJunitReport) GetErrorCode() (string, error) {
	if y.IsFileAnalze() {
		return strings.Split(y.Reference, ":")[0], nil
	}
	return "", errors.New("No Reference")
}


// func (y YamlJunitReport) GetTotal() int {
// 	return len(y)
// }

// func (y YamlJunitReport) GetErrorCount() int {
// 	count := 0
// 	for _, raw_data := range y {
// 		if raw_data.Level == error_code {
// 			count++
// 		}
// 	}
// 	return count
// }

// func (y YamlJunitReport) GetWarningCount() int {
// 	count := 0
// 	for _, raw_data := range y {
// 		if raw_data.Level == warning_code {
// 			count++
// 		}
// 	}
// 	return count
// }

// func (y YamlJunitReport) MakeReport() error {
// 	return nil
// }

