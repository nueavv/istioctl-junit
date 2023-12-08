package report

import (
	//"os"
	//"encoding/xml"
	//"fmt"
)

const (
	error_code = "Error"
	warning_code = "Warning"
)


type JunitReport interface {
	GetCode() string
	GetMessage() string
	GetOrigin() string
	GetURL() string
	GetLevel() string
	IsFileAnalze() bool
	GetErrorLine() (int, error)
	GetErrorCode() (string, error)
}

func GetErrorCount[T JunitReport](reports []T) int {
	count := 0
	for _, r := range reports {
		if r.GetLevel() == error_code {
			count++
		}
	}
	return count
}

func GetWarningCount[T JunitReport](reports []T) int {
	count := 0
	for _, r := range reports {
		if r.GetLevel() == warning_code {
			count++
		}
	}
	return count
}

func GetTotal[T JunitReport](reports []T) int {
	return len(reports)
}



func MakeReport[T JunitReport](reports []T, output string) error {
	// totals := junit.Totals{
	// 	Tests: GetTotal(reports), 
	// 	Passed: 0,
	// 	Skipped: 0,
	// 	Failed: GetWarningCount(reports),
	// 	Error: GetErrorCount(reports),
	// }

	// // 테스트 결과 파싱해서 넣기 
	// var tests []junit.Test
	// for _, report := range reports {
	// 	var status junit.Status
	// 	if report.GetLevel() == error_code {
	// 		status = junit.StatusError
	// 	} else if report.GetLevel() == warning_code {
	// 		status = junit.StatusFailed
	// 	}

	// 	tests = append(tests, junit.Test{
	// 			Name: report.GetCode(),
	// 			Message: report.GetMessage(),
	// 			Status: status,
	// 		})
	// }

	// var suites []junit.Suite
	// suites = append(suites, junit.Suite{
	// 	Name: "istioctl analyze",
	// 	Package: "istioctl", 
	// 	Tests: tests,
	// 	Totals: totals,
	// })

	// for _, suite := range suites {
	// 	fmt.Println(suite.Name)
	// 	for _, test := range suite.Tests {
	// 		fmt.Printf("  %s\n", test.Name)
	// 		if test.Error != nil {
	// 			fmt.Printf("    %s: %v\n", test.Status, test.Error)
	// 		} else {
	// 			fmt.Printf("    %s\n", test.Status)
	// 		}
	// 	}
	// }

	// xmlBytes, err := xml.Marshal(suite)
    // if err != nil {
    //     return err
    // }

	// error := os.WriteFile(output, xmlBytes, 0660)
	// if error != nil {
	// 	return error
	// }

	return nil
}


