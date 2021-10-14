package augo

import "fmt"

const (
	Windows = "Windows"
	MacOS   = "MacOS"
	Linux   = "Linux"
)

//默認使用MacOS環境的路徑格式
var system_version = MacOS
var pathChar = `/`

func SetSystemVersion(version string) {
	switch version {
	case Windows:
		pathChar = `/`
	case MacOS:
		pathChar = `/`
	case Linux:
		pathChar = `/`
	default:
		panic(fmt.Sprintf("version:= %s,input version is not format", version))
	}
	system_version = version
}

func GetSystemVersion() string {
	return system_version
}
