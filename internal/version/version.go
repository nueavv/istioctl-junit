package version

import (
	"fmt"
	"runtime"
)

const (
	defaultVersion   = "99.99.99"
	defaultBuildDate = "1970-01-01T00:00:00Z"
)

// 빌드 정보 (링크 플래그로 설정됨)
var (
	version        = defaultVersion
	buildDate      = defaultBuildDate
	gitCommit      = ""
	gitTag         = ""
	gitTreeState   = ""
	extraBuildInfo = ""
)

// Version은 Argo 버전 정보를 포함합니다.
type Version struct {
	Version        string
	BuildDate      string
	GitCommit      string
	GitTag         string
	GitTreeState   string
	GoVersion      string
	Compiler       string
	Platform       string
	ExtraBuildInfo string
}

func (v Version) String() string {
	return v.Version
}

// GetVersion은 버전 정보를 반환합니다.
func GetVersion() Version {
	versionStr := formVersionString()

	return Version{
		Version:        versionStr,
		BuildDate:      buildDate,
		GitCommit:      gitCommit,
		GitTag:         gitTag,
		GitTreeState:   gitTreeState,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		ExtraBuildInfo: extraBuildInfo,
	}
}

func formVersionString() string {
	if isOfficialRelease() {
		return gitTag
	}

	versionStr := "v" + version
	if len(gitCommit) >= 7 {
		versionStr += "+" + gitCommit[:7]
		if gitTreeState != "clean" {
			versionStr += ".dirty"
		}
	} else {
		versionStr += "+unknown"
	}
	return versionStr
}

func isOfficialRelease() bool {
	return gitCommit != "" && gitTag != "" && gitTreeState == "clean"
}
