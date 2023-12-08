package report

import (
	"strings"
	"errors"
	"strconv"
)

type JsonJunitReport struct {
	Code             string `json:"code"`
	Message          string `json:"message"`
	Origin           string `json:"origin"`
	DocumentationURL string `json:"documentationUrl"`
	Level            string `json:"level"`
	Reference		 string `json:"reference,omitempty"`
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

func (j JsonJunitReport) GetLevel() string {
	return j.Level
}

func (j JsonJunitReport) IsFileAnalze() bool {
	if j.Reference == "" {
		return false
	}
	return true
}

func (j JsonJunitReport) GetErrorLine() (int, error) {
	if j.IsFileAnalze() {
		line := strings.Split(j.Reference, ":")[1]
		i, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, errors.New("No Reference")
}

func (j JsonJunitReport) GetErrorCode() (string, error) {
	if j.IsFileAnalze() {
		return strings.Split(j.Reference, ":")[0], nil
	}
	return "", errors.New("No Reference")
}

